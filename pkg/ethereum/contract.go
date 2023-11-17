// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"aftr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"before\",\"type\":\"uint256\"}],\"name\":\"IncorrectRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"PubAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"PubDoesNotExist\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"pub\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"CIDAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"pub\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"PubCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PUB_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"pub\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"addCID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"pub\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"cidsAtTimestamp\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"pub\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"aftr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"before\",\"type\":\"uint256\"}],\"name\":\"cidsInRange\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"pub\",\"type\":\"string\"}],\"name\":\"createPub\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"pubsOfOwner\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461009d573360009081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604081205460ff161561004f575b5061127b90816100a38239f35b808052806020526040812033825260205260408120600160ff19825416179055339033907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d8180a438610042565b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806301ffc9a7146100e7578063248a9ca3146100e257806326294a77146100dd5780632f2ff15d146100d857806336568abe146100d357806352b62b3e146100ce578063822ba40b146100c957806391d14854146100c4578063a217fddf146100bf578063d41bc3ae146100ba578063d547741f146100b5578063de665dbc146100b05763fd936858146100ab57600080fd5b610924565b6107f8565b6107b9565b610695565b610679565b610627565b6105ec565b6104da565b610418565b610354565b610247565b610142565b3461013d57602036600319011261013d5760043563ffffffff60e01b811680910361013d57602090637965db0b60e01b811490811561012c575b506040519015158152f35b6301ffc9a760e01b14905038610121565b600080fd5b3461013d57602036600319011261013d5760043560005260006020526020600160406000200154604051908152f35b600435906001600160a01b038216820361013d57565b602435906001600160a01b038216820361013d57565b60005b8381106101b05750506000910152565b81810151838201526020016101a0565b906020916101d98151809281855285808601910161019d565b601f01601f1916010190565b602080820190808352835180925260408301928160408460051b8301019501936000915b8483106102195750505050505090565b9091929394958480610237600193603f198682030187528a516101c0565b9801930193019194939290610209565b3461013d5760208060031936011261013d576001600160a01b03610269610171565b16600090815260028252604080822091825490610285826110f2565b9361029284519586610cae565b82855281528481209481908086015b8483106102b9578551806102b589826101e5565b0390f35b85518285928a54926102ca84610f41565b8082526001948086169081156103385750600114610300575b506102f2816001960382610cae565b8152019801920191966102a1565b8c8952838920955088905b80821061032157508101830194506102f26102e3565b86548383018601529585019587949091019061030b565b60ff19168584015250151560051b8101830194506102f26102e3565b3461013d57604036600319011261013d57600435610370610187565b6000918083528260205261038a6001604085200154610b9d565b808352602083815260408085206001600160a01b0385166000908152925290205460ff16156103b7578280f35b808352602083815260408085206001600160a01b038516600090815292529020805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d8480a438808280f35b3461013d57604036600319011261013d57610431610187565b336001600160a01b0382160361044f5761044d90600435610ce4565b005b60405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b6064820152608490fd5b9181601f8401121561013d5782359167ffffffffffffffff831161013d576020838186019501011161013d57565b3461013d57604036600319011261013d576104f3610171565b60243567ffffffffffffffff811161013d576105139036906004016104ac565b61051b610a25565b60405191818184376001838301908152839003602001909220546001600160a01b03929083166105cc57908161057785610558846105a296610ece565b80546001600160a01b0319166001600160a01b03909216919091179055565b61059d82826105988860018060a01b03166000526002602052604060002090565b610fcf565b6110ce565b9116907ff8debc2f1745eba86909890f2dc061624705c74329348829e04aba43c015b9a2600080a3005b6105e8604051928392635c78f6ed60e11b845260048401610f19565b0390fd5b3461013d57600036600319011261013d5760206040517fafda658ee731b8f86292e3b52a311534cd93642b12a698012439316e0c3a09958152f35b3461013d57604036600319011261013d57602060ff61066d610647610187565b6004356000526000845260406000209060018060a01b0316600052602052604060002090565b54166040519015158152f35b3461013d57600036600319011261013d57602060405160008152f35b3461013d5760408060031936011261013d5760043567ffffffffffffffff811161013d576106ca6106d09136906004016104ac565b90610ee7565b9060009060243582526020928352808220918254906106ee826110f2565b936106fb84519586610cae565b82855281528481209481908086015b84831061071e578551806102b589826101e5565b85518285928a549261072f84610f41565b80825260019480861690811561079d5750600114610765575b50610757816001960382610cae565b81520198019201919661070a565b8c8952838920955088905b8082106107865750810183019450610757610748565b865483830186015295850195879490910190610770565b60ff19168584015250151560051b810183019450610757610748565b3461013d57604036600319011261013d5761044d6004356107d8610187565b908060005260006020526107f3600160406000200154610b9d565b610ce4565b3461013d57606036600319011261013d5760043567ffffffffffffffff811161013d576108299036906004016104ac565b60243591604435908184101561090257929061084e6108488486610f00565b546111e7565b9260009161085c8394610d84565b925b81841061087657848652604051806102b588826101e5565b6108a061089b8561088c868b9a9997989a610ee7565b90600052602052604060002090565b61110a565b9381935b85518510156108e6576108da6108e0916108be8789611231565b516108c9828b611231565b526108d4818a611231565b506110e3565b946110e3565b936108a4565b979295969093506108f89194506110e3565b939095939261085e565b5060405163bc0c888560e01b8152600481018490526024810191909152604490fd5b3461013d57606036600319011261013d5767ffffffffffffffff60043581811161013d576109569036906004016104ac565b9160243590811161013d5761096f9036906004016104ac565b61097a939193610a25565b6040518284823760018184019081528190036020019020546001600160a01b0316938415610a095781816109c06109e195946109db9461059860443561088c898c610ee7565b6109ca8487610f00565b6109d481546110e3565b90556110ce565b926110ce565b907fe3f9a45ba3cdf7457d983d516788bd8a5d69a802c3bcb430d08e865b114986f0600080a4005b6040516315e6e0eb60e21b8152806105e8858760048401610f19565b3360009081527f1b025c5f7493127e9e4262519d0b051a3767d7b241de3e0684fd56f9e8235b6060205260409020547fafda658ee731b8f86292e3b52a311534cd93642b12a698012439316e0c3a09959060ff1615610a815750565b610a8a33610e4e565b610a92610d97565b916030610a9e84610dc8565b536078610aaa84610dd5565b5360415b60018111610b56576105e86048610b3e85610b3088610acd8815610e03565b6040519485937f416363657373436f6e74726f6c3a206163636f756e74200000000000000000006020860152610b0d81518092602060378901910161019d565b84017001034b99036b4b9b9b4b733903937b6329607d1b60378201520190610c44565b03601f198101835282610cae565b60405162461bcd60e51b815291829160048301610cd0565b90600f8116906010821015610b9857610b93916f181899199a1a9b1b9c1cb0b131b232b360811b901a610b898487610de5565b5360041c91610df6565b610aae565b610db2565b60008181526020818152604080832033845290915290205460ff1615610bc05750565b610bc933610e4e565b610bd1610d97565b916030610bdd84610dc8565b536078610be984610dd5565b5360415b60018111610c0c576105e86048610b3e85610b3088610acd8815610e03565b90600f8116906010821015610b9857610c3f916f181899199a1a9b1b9c1cb0b131b232b360811b901a610b898487610de5565b610bed565b90610c576020928281519485920161019d565b0190565b634e487b7160e01b600052604160045260246000fd5b6080810190811067ffffffffffffffff821117610c8d57604052565b610c5b565b6060810190811067ffffffffffffffff821117610c8d57604052565b90601f8019910116810190811067ffffffffffffffff821117610c8d57604052565b906020610ce19281815201906101c0565b90565b6000818152602081815260408083206001600160a01b038616845290915281205490919060ff16610d1457505050565b808252602082815260408084206001600160a01b038616600090815292529020805460ff1916905533926001600160a01b0316917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4565b634e487b7160e01b600052601160045260246000fd5b9060018201809211610d9257565b610d6e565b60405190610da482610c71565b604282526060366020840137565b634e487b7160e01b600052603260045260246000fd5b805115610b985760200190565b805160011015610b985760210190565b908151811015610b98570160200190565b8015610d92576000190190565b15610e0a57565b606460405162461bcd60e51b815260206004820152602060248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152fd5b60405190610e5b82610c92565b602a825260403660208401376030610e7283610dc8565b536078610e7e83610dd5565b536029905b60018211610e9657610ce1915015610e03565b600f8116906010821015610b9857610ec8916f181899199a1a9b1b9c1cb0b131b232b360811b901a610b898486610de5565b90610e83565b6020908260405193849283378101600181520301902090565b6020908260405193849283378101600481520301902090565b6020908260405193849283378101600381520301902090565b90918060409360208452816020850152848401376000828201840152601f01601f1916010190565b90600182811c92168015610f71575b6020831014610f5b57565b634e487b7160e01b600052602260045260246000fd5b91607f1691610f50565b90601f8111610f8957505050565b600091825260208220906020601f850160051c83019410610fc5575b601f0160051c01915b828110610fba57505050565b818155600101610fae565b9092508290610fa5565b9091815468010000000000000000811015610c8d57600192838201808255821015610b985760005260209081600020019367ffffffffffffffff8311610c8d576110238361101d8754610f41565b87610f7b565b600091601f841160011461106557506110569350600091908361105a575b50508160011b916000199060031b1c19161790565b9055565b013590503880611041565b9183601f19811661107b88600052602060002090565b9483905b888383106110b4575050501061109a575b505050811b019055565b0135600019600384901b60f8161c19169055388080611090565b86860135885590960195938401938793509081019061107f565b81604051928392833781016000815203902090565b6000198114610d925760010190565b67ffffffffffffffff8111610c8d5760051b60200190565b90815491611117836110f2565b9260409161112783519586610cae565b81855260009081526020808220938291908188015b85841061114c5750505050505050565b815183869289549261115d84610f41565b8082526001948086169081156111cb5750600114611193575b50611185816001960382610cae565b81520197019301929561113c565b8b8a52838a20955089905b8082106111b45750810183019450611185611176565b86548383018601529585019588949091019061119e565b60ff19168584015250151560051b810183019450611185611176565b906111f1826110f2565b6111fe6040519182610cae565b828152809261120f601f19916110f2565b019060005b82811061122057505050565b806060602080938501015201611214565b8051821015610b985760209160051b01019056fea2646970667358221220a67108681c070b8892513194cde0d9c1950175e4854f2faf7838c4b627be29c364736f6c63430008150033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Contract *ContractCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Contract *ContractSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Contract.Contract.DEFAULTADMINROLE(&_Contract.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Contract *ContractCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Contract.Contract.DEFAULTADMINROLE(&_Contract.CallOpts)
}

// PUBADMINROLE is a free data retrieval call binding the contract method 0x822ba40b.
//
// Solidity: function PUB_ADMIN_ROLE() view returns(bytes32)
func (_Contract *ContractCaller) PUBADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "PUB_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PUBADMINROLE is a free data retrieval call binding the contract method 0x822ba40b.
//
// Solidity: function PUB_ADMIN_ROLE() view returns(bytes32)
func (_Contract *ContractSession) PUBADMINROLE() ([32]byte, error) {
	return _Contract.Contract.PUBADMINROLE(&_Contract.CallOpts)
}

// PUBADMINROLE is a free data retrieval call binding the contract method 0x822ba40b.
//
// Solidity: function PUB_ADMIN_ROLE() view returns(bytes32)
func (_Contract *ContractCallerSession) PUBADMINROLE() ([32]byte, error) {
	return _Contract.Contract.PUBADMINROLE(&_Contract.CallOpts)
}

// CidsAtTimestamp is a free data retrieval call binding the contract method 0xd41bc3ae.
//
// Solidity: function cidsAtTimestamp(string pub, uint256 epoch) view returns(string[])
func (_Contract *ContractCaller) CidsAtTimestamp(opts *bind.CallOpts, pub string, epoch *big.Int) ([]string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "cidsAtTimestamp", pub, epoch)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// CidsAtTimestamp is a free data retrieval call binding the contract method 0xd41bc3ae.
//
// Solidity: function cidsAtTimestamp(string pub, uint256 epoch) view returns(string[])
func (_Contract *ContractSession) CidsAtTimestamp(pub string, epoch *big.Int) ([]string, error) {
	return _Contract.Contract.CidsAtTimestamp(&_Contract.CallOpts, pub, epoch)
}

// CidsAtTimestamp is a free data retrieval call binding the contract method 0xd41bc3ae.
//
// Solidity: function cidsAtTimestamp(string pub, uint256 epoch) view returns(string[])
func (_Contract *ContractCallerSession) CidsAtTimestamp(pub string, epoch *big.Int) ([]string, error) {
	return _Contract.Contract.CidsAtTimestamp(&_Contract.CallOpts, pub, epoch)
}

// CidsInRange is a free data retrieval call binding the contract method 0xde665dbc.
//
// Solidity: function cidsInRange(string pub, uint256 aftr, uint256 before) view returns(string[])
func (_Contract *ContractCaller) CidsInRange(opts *bind.CallOpts, pub string, aftr *big.Int, before *big.Int) ([]string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "cidsInRange", pub, aftr, before)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// CidsInRange is a free data retrieval call binding the contract method 0xde665dbc.
//
// Solidity: function cidsInRange(string pub, uint256 aftr, uint256 before) view returns(string[])
func (_Contract *ContractSession) CidsInRange(pub string, aftr *big.Int, before *big.Int) ([]string, error) {
	return _Contract.Contract.CidsInRange(&_Contract.CallOpts, pub, aftr, before)
}

// CidsInRange is a free data retrieval call binding the contract method 0xde665dbc.
//
// Solidity: function cidsInRange(string pub, uint256 aftr, uint256 before) view returns(string[])
func (_Contract *ContractCallerSession) CidsInRange(pub string, aftr *big.Int, before *big.Int) ([]string, error) {
	return _Contract.Contract.CidsInRange(&_Contract.CallOpts, pub, aftr, before)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Contract *ContractCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Contract *ContractSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Contract.Contract.GetRoleAdmin(&_Contract.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Contract *ContractCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Contract.Contract.GetRoleAdmin(&_Contract.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Contract *ContractCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Contract *ContractSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Contract.Contract.HasRole(&_Contract.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Contract *ContractCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Contract.Contract.HasRole(&_Contract.CallOpts, role, account)
}

// PubsOfOwner is a free data retrieval call binding the contract method 0x26294a77.
//
// Solidity: function pubsOfOwner(address owner) view returns(string[])
func (_Contract *ContractCaller) PubsOfOwner(opts *bind.CallOpts, owner common.Address) ([]string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "pubsOfOwner", owner)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// PubsOfOwner is a free data retrieval call binding the contract method 0x26294a77.
//
// Solidity: function pubsOfOwner(address owner) view returns(string[])
func (_Contract *ContractSession) PubsOfOwner(owner common.Address) ([]string, error) {
	return _Contract.Contract.PubsOfOwner(&_Contract.CallOpts, owner)
}

// PubsOfOwner is a free data retrieval call binding the contract method 0x26294a77.
//
// Solidity: function pubsOfOwner(address owner) view returns(string[])
func (_Contract *ContractCallerSession) PubsOfOwner(owner common.Address) ([]string, error) {
	return _Contract.Contract.PubsOfOwner(&_Contract.CallOpts, owner)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Contract *ContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Contract *ContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Contract.Contract.SupportsInterface(&_Contract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Contract *ContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Contract.Contract.SupportsInterface(&_Contract.CallOpts, interfaceId)
}

// AddCID is a paid mutator transaction binding the contract method 0xfd936858.
//
// Solidity: function addCID(string pub, string cid, uint256 timestamp) returns()
func (_Contract *ContractTransactor) AddCID(opts *bind.TransactOpts, pub string, cid string, timestamp *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addCID", pub, cid, timestamp)
}

// AddCID is a paid mutator transaction binding the contract method 0xfd936858.
//
// Solidity: function addCID(string pub, string cid, uint256 timestamp) returns()
func (_Contract *ContractSession) AddCID(pub string, cid string, timestamp *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.AddCID(&_Contract.TransactOpts, pub, cid, timestamp)
}

// AddCID is a paid mutator transaction binding the contract method 0xfd936858.
//
// Solidity: function addCID(string pub, string cid, uint256 timestamp) returns()
func (_Contract *ContractTransactorSession) AddCID(pub string, cid string, timestamp *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.AddCID(&_Contract.TransactOpts, pub, cid, timestamp)
}

// CreatePub is a paid mutator transaction binding the contract method 0x52b62b3e.
//
// Solidity: function createPub(address owner, string pub) returns()
func (_Contract *ContractTransactor) CreatePub(opts *bind.TransactOpts, owner common.Address, pub string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createPub", owner, pub)
}

// CreatePub is a paid mutator transaction binding the contract method 0x52b62b3e.
//
// Solidity: function createPub(address owner, string pub) returns()
func (_Contract *ContractSession) CreatePub(owner common.Address, pub string) (*types.Transaction, error) {
	return _Contract.Contract.CreatePub(&_Contract.TransactOpts, owner, pub)
}

// CreatePub is a paid mutator transaction binding the contract method 0x52b62b3e.
//
// Solidity: function createPub(address owner, string pub) returns()
func (_Contract *ContractTransactorSession) CreatePub(owner common.Address, pub string) (*types.Transaction, error) {
	return _Contract.Contract.CreatePub(&_Contract.TransactOpts, owner, pub)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Contract *ContractTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Contract *ContractSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.Contract.GrantRole(&_Contract.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Contract *ContractTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.Contract.GrantRole(&_Contract.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Contract *ContractTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Contract *ContractSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RenounceRole(&_Contract.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Contract *ContractTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RenounceRole(&_Contract.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Contract *ContractTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Contract *ContractSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RevokeRole(&_Contract.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Contract *ContractTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RevokeRole(&_Contract.TransactOpts, role, account)
}

// ContractCIDAddedIterator is returned from FilterCIDAdded and is used to iterate over the raw logs and unpacked data for CIDAdded events raised by the Contract contract.
type ContractCIDAddedIterator struct {
	Event *ContractCIDAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractCIDAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractCIDAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractCIDAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractCIDAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractCIDAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractCIDAdded represents a CIDAdded event raised by the Contract contract.
type ContractCIDAdded struct {
	Cid   common.Hash
	Pub   common.Hash
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCIDAdded is a free log retrieval operation binding the contract event 0xe3f9a45ba3cdf7457d983d516788bd8a5d69a802c3bcb430d08e865b114986f0.
//
// Solidity: event CIDAdded(string indexed cid, string indexed pub, address indexed owner)
func (_Contract *ContractFilterer) FilterCIDAdded(opts *bind.FilterOpts, cid []string, pub []string, owner []common.Address) (*ContractCIDAddedIterator, error) {

	var cidRule []interface{}
	for _, cidItem := range cid {
		cidRule = append(cidRule, cidItem)
	}
	var pubRule []interface{}
	for _, pubItem := range pub {
		pubRule = append(pubRule, pubItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "CIDAdded", cidRule, pubRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ContractCIDAddedIterator{contract: _Contract.contract, event: "CIDAdded", logs: logs, sub: sub}, nil
}

// WatchCIDAdded is a free log subscription operation binding the contract event 0xe3f9a45ba3cdf7457d983d516788bd8a5d69a802c3bcb430d08e865b114986f0.
//
// Solidity: event CIDAdded(string indexed cid, string indexed pub, address indexed owner)
func (_Contract *ContractFilterer) WatchCIDAdded(opts *bind.WatchOpts, sink chan<- *ContractCIDAdded, cid []string, pub []string, owner []common.Address) (event.Subscription, error) {

	var cidRule []interface{}
	for _, cidItem := range cid {
		cidRule = append(cidRule, cidItem)
	}
	var pubRule []interface{}
	for _, pubItem := range pub {
		pubRule = append(pubRule, pubItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "CIDAdded", cidRule, pubRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractCIDAdded)
				if err := _Contract.contract.UnpackLog(event, "CIDAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCIDAdded is a log parse operation binding the contract event 0xe3f9a45ba3cdf7457d983d516788bd8a5d69a802c3bcb430d08e865b114986f0.
//
// Solidity: event CIDAdded(string indexed cid, string indexed pub, address indexed owner)
func (_Contract *ContractFilterer) ParseCIDAdded(log types.Log) (*ContractCIDAdded, error) {
	event := new(ContractCIDAdded)
	if err := _Contract.contract.UnpackLog(event, "CIDAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPubCreatedIterator is returned from FilterPubCreated and is used to iterate over the raw logs and unpacked data for PubCreated events raised by the Contract contract.
type ContractPubCreatedIterator struct {
	Event *ContractPubCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPubCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPubCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPubCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPubCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPubCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPubCreated represents a PubCreated event raised by the Contract contract.
type ContractPubCreated struct {
	Pub   common.Hash
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterPubCreated is a free log retrieval operation binding the contract event 0xf8debc2f1745eba86909890f2dc061624705c74329348829e04aba43c015b9a2.
//
// Solidity: event PubCreated(string indexed pub, address indexed owner)
func (_Contract *ContractFilterer) FilterPubCreated(opts *bind.FilterOpts, pub []string, owner []common.Address) (*ContractPubCreatedIterator, error) {

	var pubRule []interface{}
	for _, pubItem := range pub {
		pubRule = append(pubRule, pubItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "PubCreated", pubRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ContractPubCreatedIterator{contract: _Contract.contract, event: "PubCreated", logs: logs, sub: sub}, nil
}

// WatchPubCreated is a free log subscription operation binding the contract event 0xf8debc2f1745eba86909890f2dc061624705c74329348829e04aba43c015b9a2.
//
// Solidity: event PubCreated(string indexed pub, address indexed owner)
func (_Contract *ContractFilterer) WatchPubCreated(opts *bind.WatchOpts, sink chan<- *ContractPubCreated, pub []string, owner []common.Address) (event.Subscription, error) {

	var pubRule []interface{}
	for _, pubItem := range pub {
		pubRule = append(pubRule, pubItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "PubCreated", pubRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPubCreated)
				if err := _Contract.contract.UnpackLog(event, "PubCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePubCreated is a log parse operation binding the contract event 0xf8debc2f1745eba86909890f2dc061624705c74329348829e04aba43c015b9a2.
//
// Solidity: event PubCreated(string indexed pub, address indexed owner)
func (_Contract *ContractFilterer) ParsePubCreated(log types.Log) (*ContractPubCreated, error) {
	event := new(ContractPubCreated)
	if err := _Contract.contract.UnpackLog(event, "PubCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Contract contract.
type ContractRoleAdminChangedIterator struct {
	Event *ContractRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRoleAdminChanged represents a RoleAdminChanged event raised by the Contract contract.
type ContractRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Contract *ContractFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ContractRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ContractRoleAdminChangedIterator{contract: _Contract.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Contract *ContractFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ContractRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRoleAdminChanged)
				if err := _Contract.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Contract *ContractFilterer) ParseRoleAdminChanged(log types.Log) (*ContractRoleAdminChanged, error) {
	event := new(ContractRoleAdminChanged)
	if err := _Contract.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Contract contract.
type ContractRoleGrantedIterator struct {
	Event *ContractRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRoleGranted represents a RoleGranted event raised by the Contract contract.
type ContractRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Contract *ContractFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractRoleGrantedIterator{contract: _Contract.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Contract *ContractFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ContractRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRoleGranted)
				if err := _Contract.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Contract *ContractFilterer) ParseRoleGranted(log types.Log) (*ContractRoleGranted, error) {
	event := new(ContractRoleGranted)
	if err := _Contract.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Contract contract.
type ContractRoleRevokedIterator struct {
	Event *ContractRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRoleRevoked represents a RoleRevoked event raised by the Contract contract.
type ContractRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Contract *ContractFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractRoleRevokedIterator{contract: _Contract.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Contract *ContractFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ContractRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRoleRevoked)
				if err := _Contract.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Contract *ContractFilterer) ParseRoleRevoked(log types.Log) (*ContractRoleRevoked, error) {
	event := new(ContractRoleRevoked)
	if err := _Contract.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
