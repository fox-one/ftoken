// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package quorum

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FTokenFactoryToken is an auto generated low-level Go binding around an user-defined struct.
type FTokenFactoryToken struct {
	Name   string
	Symbol string
	Cap    *big.Int
	Trace  *big.Int
	Minter common.Address
}

// QuorumABI is the input ABI used to generate the binding from.
const QuorumABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"cap\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"trace\",\"type\":\"uint128\"}],\"name\":\"createContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"raw\",\"type\":\"bytes\"},{\"internalType\":\"uint128\",\"name\":\"trace\",\"type\":\"uint128\"}],\"name\":\"createContractRaw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"readToken\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"cap\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"trace\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"internalType\":\"structFTokenFactory.Token\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"receiver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"setReceiverAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Quorum is an auto generated Go binding around an Ethereum contract.
type Quorum struct {
	QuorumCaller     // Read-only binding to the contract
	QuorumTransactor // Write-only binding to the contract
	QuorumFilterer   // Log filterer for contract events
}

// QuorumCaller is an auto generated read-only Go binding around an Ethereum contract.
type QuorumCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuorumTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QuorumTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuorumFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuorumFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuorumSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuorumSession struct {
	Contract     *Quorum           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuorumCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuorumCallerSession struct {
	Contract *QuorumCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// QuorumTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuorumTransactorSession struct {
	Contract     *QuorumTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuorumRaw is an auto generated low-level Go binding around an Ethereum contract.
type QuorumRaw struct {
	Contract *Quorum // Generic contract binding to access the raw methods on
}

// QuorumCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuorumCallerRaw struct {
	Contract *QuorumCaller // Generic read-only contract binding to access the raw methods on
}

// QuorumTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuorumTransactorRaw struct {
	Contract *QuorumTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQuorum creates a new instance of Quorum, bound to a specific deployed contract.
func NewQuorum(address common.Address, backend bind.ContractBackend) (*Quorum, error) {
	contract, err := bindQuorum(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Quorum{QuorumCaller: QuorumCaller{contract: contract}, QuorumTransactor: QuorumTransactor{contract: contract}, QuorumFilterer: QuorumFilterer{contract: contract}}, nil
}

// NewQuorumCaller creates a new read-only instance of Quorum, bound to a specific deployed contract.
func NewQuorumCaller(address common.Address, caller bind.ContractCaller) (*QuorumCaller, error) {
	contract, err := bindQuorum(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuorumCaller{contract: contract}, nil
}

// NewQuorumTransactor creates a new write-only instance of Quorum, bound to a specific deployed contract.
func NewQuorumTransactor(address common.Address, transactor bind.ContractTransactor) (*QuorumTransactor, error) {
	contract, err := bindQuorum(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuorumTransactor{contract: contract}, nil
}

// NewQuorumFilterer creates a new log filterer instance of Quorum, bound to a specific deployed contract.
func NewQuorumFilterer(address common.Address, filterer bind.ContractFilterer) (*QuorumFilterer, error) {
	contract, err := bindQuorum(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuorumFilterer{contract: contract}, nil
}

// bindQuorum binds a generic wrapper to an already deployed contract.
func bindQuorum(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(QuorumABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Quorum *QuorumRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Quorum.Contract.QuorumCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Quorum *QuorumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quorum.Contract.QuorumTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Quorum *QuorumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Quorum.Contract.QuorumTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Quorum *QuorumCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Quorum.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Quorum *QuorumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quorum.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Quorum *QuorumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Quorum.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Quorum *QuorumCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Quorum.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Quorum *QuorumSession) Owner() (common.Address, error) {
	return _Quorum.Contract.Owner(&_Quorum.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Quorum *QuorumCallerSession) Owner() (common.Address, error) {
	return _Quorum.Contract.Owner(&_Quorum.CallOpts)
}

// ReadToken is a free data retrieval call binding the contract method 0x8fddced8.
//
// Solidity: function readToken(address _address) view returns((string,string,uint256,uint128,address))
func (_Quorum *QuorumCaller) ReadToken(opts *bind.CallOpts, _address common.Address) (FTokenFactoryToken, error) {
	var out []interface{}
	err := _Quorum.contract.Call(opts, &out, "readToken", _address)

	if err != nil {
		return *new(FTokenFactoryToken), err
	}

	out0 := *abi.ConvertType(out[0], new(FTokenFactoryToken)).(*FTokenFactoryToken)

	return out0, err

}

// ReadToken is a free data retrieval call binding the contract method 0x8fddced8.
//
// Solidity: function readToken(address _address) view returns((string,string,uint256,uint128,address))
func (_Quorum *QuorumSession) ReadToken(_address common.Address) (FTokenFactoryToken, error) {
	return _Quorum.Contract.ReadToken(&_Quorum.CallOpts, _address)
}

// ReadToken is a free data retrieval call binding the contract method 0x8fddced8.
//
// Solidity: function readToken(address _address) view returns((string,string,uint256,uint128,address))
func (_Quorum *QuorumCallerSession) ReadToken(_address common.Address) (FTokenFactoryToken, error) {
	return _Quorum.Contract.ReadToken(&_Quorum.CallOpts, _address)
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() view returns(address)
func (_Quorum *QuorumCaller) Receiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Quorum.contract.Call(opts, &out, "receiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() view returns(address)
func (_Quorum *QuorumSession) Receiver() (common.Address, error) {
	return _Quorum.Contract.Receiver(&_Quorum.CallOpts)
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() view returns(address)
func (_Quorum *QuorumCallerSession) Receiver() (common.Address, error) {
	return _Quorum.Contract.Receiver(&_Quorum.CallOpts)
}

// CreateContract is a paid mutator transaction binding the contract method 0x6adb3a13.
//
// Solidity: function createContract(string name, string symbol, uint256 cap, uint128 trace) returns()
func (_Quorum *QuorumTransactor) CreateContract(opts *bind.TransactOpts, name string, symbol string, cap *big.Int, trace *big.Int) (*types.Transaction, error) {
	return _Quorum.contract.Transact(opts, "createContract", name, symbol, cap, trace)
}

// CreateContract is a paid mutator transaction binding the contract method 0x6adb3a13.
//
// Solidity: function createContract(string name, string symbol, uint256 cap, uint128 trace) returns()
func (_Quorum *QuorumSession) CreateContract(name string, symbol string, cap *big.Int, trace *big.Int) (*types.Transaction, error) {
	return _Quorum.Contract.CreateContract(&_Quorum.TransactOpts, name, symbol, cap, trace)
}

// CreateContract is a paid mutator transaction binding the contract method 0x6adb3a13.
//
// Solidity: function createContract(string name, string symbol, uint256 cap, uint128 trace) returns()
func (_Quorum *QuorumTransactorSession) CreateContract(name string, symbol string, cap *big.Int, trace *big.Int) (*types.Transaction, error) {
	return _Quorum.Contract.CreateContract(&_Quorum.TransactOpts, name, symbol, cap, trace)
}

// CreateContractRaw is a paid mutator transaction binding the contract method 0x990ac85d.
//
// Solidity: function createContractRaw(bytes raw, uint128 trace) returns()
func (_Quorum *QuorumTransactor) CreateContractRaw(opts *bind.TransactOpts, raw []byte, trace *big.Int) (*types.Transaction, error) {
	return _Quorum.contract.Transact(opts, "createContractRaw", raw, trace)
}

// CreateContractRaw is a paid mutator transaction binding the contract method 0x990ac85d.
//
// Solidity: function createContractRaw(bytes raw, uint128 trace) returns()
func (_Quorum *QuorumSession) CreateContractRaw(raw []byte, trace *big.Int) (*types.Transaction, error) {
	return _Quorum.Contract.CreateContractRaw(&_Quorum.TransactOpts, raw, trace)
}

// CreateContractRaw is a paid mutator transaction binding the contract method 0x990ac85d.
//
// Solidity: function createContractRaw(bytes raw, uint128 trace) returns()
func (_Quorum *QuorumTransactorSession) CreateContractRaw(raw []byte, trace *big.Int) (*types.Transaction, error) {
	return _Quorum.Contract.CreateContractRaw(&_Quorum.TransactOpts, raw, trace)
}

// SetReceiverAddress is a paid mutator transaction binding the contract method 0x8279c7db.
//
// Solidity: function setReceiverAddress(address _receiver) returns()
func (_Quorum *QuorumTransactor) SetReceiverAddress(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _Quorum.contract.Transact(opts, "setReceiverAddress", _receiver)
}

// SetReceiverAddress is a paid mutator transaction binding the contract method 0x8279c7db.
//
// Solidity: function setReceiverAddress(address _receiver) returns()
func (_Quorum *QuorumSession) SetReceiverAddress(_receiver common.Address) (*types.Transaction, error) {
	return _Quorum.Contract.SetReceiverAddress(&_Quorum.TransactOpts, _receiver)
}

// SetReceiverAddress is a paid mutator transaction binding the contract method 0x8279c7db.
//
// Solidity: function setReceiverAddress(address _receiver) returns()
func (_Quorum *QuorumTransactorSession) SetReceiverAddress(_receiver common.Address) (*types.Transaction, error) {
	return _Quorum.Contract.SetReceiverAddress(&_Quorum.TransactOpts, _receiver)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _owner) returns()
func (_Quorum *QuorumTransactor) TransferOwnership(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Quorum.contract.Transact(opts, "transferOwnership", _owner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _owner) returns()
func (_Quorum *QuorumSession) TransferOwnership(_owner common.Address) (*types.Transaction, error) {
	return _Quorum.Contract.TransferOwnership(&_Quorum.TransactOpts, _owner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _owner) returns()
func (_Quorum *QuorumTransactorSession) TransferOwnership(_owner common.Address) (*types.Transaction, error) {
	return _Quorum.Contract.TransferOwnership(&_Quorum.TransactOpts, _owner)
}
