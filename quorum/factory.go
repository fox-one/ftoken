package quorum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fox-one/ftoken/core"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type (
	Factory struct {
		privkey        *ecdsa.PrivateKey
		address        common.Address
		factoryAddress common.Address
		client         *ethclient.Client
		quorum         *Quorum
		transactor     *QuorumTransactor
	}
)

func New(ethurl, priv, factoryContract string) *Factory {
	client, err := ethclient.Dial(ethurl)
	if err != nil {
		panic(err)
	}

	privkey, err := crypto.HexToECDSA(priv)
	if err != nil {
		panic("invalid private key: " + priv)
	}
	pubKey, ok := privkey.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("cannot cast to ecdsa public key")
	}
	address := crypto.PubkeyToAddress(*pubKey)
	factoryAddress := common.HexToAddress(factoryContract)

	quorum, err := NewQuorum(factoryAddress, client)
	if err != nil {
		panic(err)
	}

	transactor, err := NewQuorumTransactor(factoryAddress, client)
	if err != nil {
		panic(err)
	}

	return &Factory{
		privkey:        privkey,
		address:        address,
		factoryAddress: factoryAddress,
		client:         client,
		quorum:         quorum,
		transactor:     transactor,
	}
}

func (*Factory) Platform() string {
	return "Ethereum"
}

func (*Factory) GasAsset() string {
	return EthAsset
}

func (f *Factory) CreateTransaction(ctx context.Context, tokens core.TokenItems, trace string) (*core.Transaction, error) {
	data, err := core.EncodeTokens(tokens)
	if err != nil {
		return nil, err
	}

	chainID, err := f.client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := f.client.PendingNonceAt(ctx, f.address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := f.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(f.privkey, chainID)
	if err != nil {
		return nil, err
	}

	opts.Nonce = big.NewInt(int64(nonce))
	opts.GasFeeCap = gasPrice
	opts.GasTipCap = big.NewInt(1000000000)       // 1 Gwei
	opts.Value = big.NewInt(0)                    // in wei
	opts.GasLimit = uint64(1000000 * len(tokens)) // in units
	opts.NoSend = true

	var traceID = &big.Int{}
	if trace != "" {
		if trace, err := uuid.FromString(trace); err == nil {
			traceID.SetBytes(trace.Bytes())
		}
	}

	txRaw, err := f.transactor.CreateContractRaw(opts, data, traceID)
	if err != nil {
		return nil, err
	}

	rawData, err := txRaw.MarshalBinary()
	if err != nil {
		return nil, err
	}

	gas, err := f.client.EstimateGas(ctx, ethereum.CallMsg{
		From:      f.address,
		To:        &f.factoryAddress,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
		Value:     opts.Value,
		Data:      txRaw.Data(),
	})
	if err != nil {
		return nil, err
	}

	fee := decimal.NewFromInt(int64(gas * gasPrice.Uint64())).Shift(-10).Ceil().Shift(-8)
	return &core.Transaction{
		Hash:  txRaw.Hash().String(),
		Raw:   hex.EncodeToString(rawData),
		Gas:   fee,
		State: core.TransactionStateNew,
	}, nil
}

func (f *Factory) SendTransaction(ctx context.Context, tx *core.Transaction) error {
	data, err := hex.DecodeString(tx.Raw)
	if err != nil {
		return err
	}

	var txEth types.Transaction
	if err := txEth.UnmarshalBinary(data); err != nil {
		return err
	}
	return f.client.SendTransaction(ctx, &txEth)
}

func (f *Factory) ReadTransaction(ctx context.Context, hash string) (*core.Transaction, error) {
	txRaw, pending, err := f.client.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}

	rawData, err := txRaw.MarshalBinary()
	if err != nil {
		return nil, err
	}

	tx := core.Transaction{
		Hash:  txRaw.Hash().String(),
		Raw:   hex.EncodeToString(rawData),
		Gas:   decimal.NewFromBigInt(txRaw.Cost(), 0).Shift(-10).Ceil().Shift(-8),
		State: core.TransactionStatePending,
	}

	if !pending {
		receipt, err := f.client.TransactionReceipt(ctx, txRaw.Hash())
		if err != nil {
			return nil, err
		}

		tx.Gas = decimal.NewFromInt(int64(receipt.GasUsed * txRaw.GasPrice().Uint64())).Shift(-10).Ceil().Shift(-8)
		if receipt.Status == 0 {
			tx.State = core.TransactionStateFailed
			return &tx, nil
		}

		tx.State = core.TransactionStateSuccess
		var opts = &bind.CallOpts{
			Context: ctx,
			Pending: false,
		}
		for i := 0; i < len(receipt.Logs)-1; i = i + 2 {
			address := receipt.Logs[i].Address
			contract, err := f.quorum.ReadToken(opts, address)
			if err != nil {
				return nil, err
			}
			if contract.Cap == nil || contract.Cap.Sign() <= 0 {
				return nil, errors.New("invalid contract cap")
			}
			tx.Tokens = append(tx.Tokens, &core.TokenItem{
				Name:        contract.Name,
				Symbol:      contract.Symbol,
				TotalSupply: contract.Cap.Uint64(),
				AssetKey:    address.String(),
				AssetID:     MixinAssetID(address.String()),
			})
		}
	}

	return &tx, nil
}
