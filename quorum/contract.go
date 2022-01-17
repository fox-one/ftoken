package quorum

//go:generate solc --optimize --abi contract/factory.sol -o contract/abi --overwrite
//go:generate abigen --abi=contract/abi/FTokenFactory.abi --pkg=quorum --out=factory.abi.go
