package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	//contractBytecode := "0x6080604052600436106100775763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663025313a281146101085780630add8140146101395780633659cfe61461014e5780635c60da1b146101715780639965b3d614610186578063f1739cae1461019b575b60006100816101bc565b9050600160a060020a03811615156100e3576040805160e560020a62461bcd02815260206004820152601f60248201527f696d706c656d656e746174696f6e20636f6e7472616374206e6f742073657400604482015290519081900360640190fd5b60405136600082376000803683855af43d806000843e818015610104578184f35b8184fd5b34801561011457600080fd5b5061011d6101f2565b60408051600160a060020a039092168252519081900360200190f35b34801561014557600080fd5b5061011d610228565b34801561015a57600080fd5b5061016f600160a060020a036004351661025e565b005b34801561017d57600080fd5b5061011d6101bc565b34801561019257600080fd5b5061016f6102dc565b3480156101a757600080fd5b5061016f600160a060020a03600435166103b8565b604080517f747275655553442e70726f78792e696d706c656d656e746174696f6e000000008152905190819003601c0190205490565b604080517f747275655553442e70726f78792e6f776e657200000000000000000000000000815290519081900360130190205490565b604080517f747275655553442e70656e64696e672e70726f78792e6f776e657200000000008152905190819003601b0190205490565b6102666101f2565b600160a060020a031633600160a060020a03161415156102d0576040805160e560020a62461bcd02815260206004820152601060248201527f6f6e6c792050726f7879204f776e657200000000000000000000000000000000604482015290519081900360640190fd5b6102d981610496565b50565b6102e4610228565b600160a060020a031633600160a060020a031614151561034e576040805160e560020a62461bcd02815260206004820152601860248201527f6f6e6c792070656e64696e672050726f7879204f776e65720000000000000000604482015290519081900360640190fd5b610356610228565b600160a060020a03166103676101f2565b600160a060020a03167f5a3e66efaa1e445ebd894728a69d6959842ea1e97bd79b892797106e270efcd960405160405180910390a36103ac6103a7610228565b6104fc565b6103b66000610531565b565b6103c06101f2565b600160a060020a031633600160a060020a031614151561042a576040805160e560020a62461bcd02815260206004820152601060248201527f6f6e6c792050726f7879204f776e657200000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038116151561043f57600080fd5b61044881610531565b7fb3d55174552271a4f1aaf36b72f50381e892171636b3fb5447fe00e995e7a37b6104716101f2565b60408051600160a060020a03928316815291841660208301528051918290030190a150565b60006104a06101bc565b9050600160a060020a0380821690831614156104bb57600080fd5b6104c482610566565b604051600160a060020a038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a25050565b604080517f747275655553442e70726f78792e6f776e6572000000000000000000000000008152905190819003601301902055565b604080517f747275655553442e70656e64696e672e70726f78792e6f776e657200000000008152905190819003601b01902055565b604080517f747275655553442e70726f78792e696d706c656d656e746174696f6e000000008152905190819003601c019020555600a165627a7a7230582066d4b5932adad66dd7b19bdf5b36c25e9e99e563f06c779eeef9427bb3c9c7670029"
	//contractABI := `[{"constant":true,"inputs":[],"name":"proxyOwner","outputs":[{"name":"owner","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"pendingProxyOwner","outputs":[{"name":"pendingOwner","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"implementation","type":"address"}],"name":"upgradeTo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"implementation","outputs":[{"name":"impl","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"claimProxyOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferProxyOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"ProxyOwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"currentOwner","type":"address"},{"indexed":false,"name":"pendingOwner","type":"address"}],"name":"NewPendingOwner","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"implementation","type":"address"}],"name":"Upgraded","type":"event"}]`
	contractABI := `[{"constant":false,"inputs":[{"name":"newImplementation","type":"address"}],"name":"upgradeTo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newImplementation","type":"address"},{"name":"data","type":"bytes"}],"name":"upgradeToAndCall","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[],"name":"implementation","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAdmin","type":"address"}],"name":"changeAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"admin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_implementation","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":false,"name":"previousAdmin","type":"address"},{"indexed":false,"name":"newAdmin","type":"address"}],"name":"AdminChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"implementation","type":"address"}],"name":"Upgraded","type":"event"}]`
	contractBytecode := `0x60806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633659cfe6146100775780634f1ef286146100ba5780635c60da1b146101085780638f2839701461015f578063f851a440146101a2575b6100756101f9565b005b34801561008357600080fd5b506100b8600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610213565b005b610106600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001908201803590602001919091929391929390505050610268565b005b34801561011457600080fd5b5061011d610308565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561016b57600080fd5b506101a0600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610360565b005b3480156101ae57600080fd5b506101b761051e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610201610576565b61021161020c610651565b610682565b565b61021b6106a8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561025c57610257816106d9565b610265565b6102646101f9565b5b50565b6102706106a8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156102fa576102ac836106d9565b3073ffffffffffffffffffffffffffffffffffffffff163483836040518083838082843782019150509250505060006040518083038185875af19250505015156102f557600080fd5b610303565b6103026101f9565b5b505050565b60006103126106a8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156103545761034d610651565b905061035d565b61035c6101f9565b5b90565b6103686106a8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561051257600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610466576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260368152602001807f43616e6e6f74206368616e6765207468652061646d696e206f6620612070726f81526020017f787920746f20746865207a65726f20616464726573730000000000000000000081525060400191505060405180910390fd5b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f61048f6106a8565b82604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a161050d81610748565b61051b565b61051a6101f9565b5b50565b60006105286106a8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561056a576105636106a8565b9050610573565b6105726101f9565b5b90565b61057e6106a8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151515610647576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260328152602001807f43616e6e6f742063616c6c2066616c6c6261636b2066756e6374696f6e20667281526020017f6f6d207468652070726f78792061646d696e000000000000000000000000000081525060400191505060405180910390fd5b61064f610777565b565b6000807f7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c36001029050805491505090565b3660008037600080366000845af43d6000803e80600081146106a3573d6000f35b3d6000fd5b6000807f10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b6001029050805491505090565b6106e281610779565b7fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b81604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150565b60007f10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b60010290508181555050565b565b60006107848261084b565b151561081e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603b8152602001807f43616e6e6f742073657420612070726f787920696d706c656d656e746174696f81526020017f6e20746f2061206e6f6e2d636f6e74726163742061646472657373000000000081525060400191505060405180910390fd5b7f7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c360010290508181555050565b600080823b9050600081119150509190505600a165627a7a72305820a4a547cfc7202c5acaaae74d428e988bc62ad5024eb0165532d3a8f91db4ed240029`
	isEIP897 := isEIP897Compliant(contractBytecode, contractABI)
	fmt.Printf("合约是否符合 EIP-897 标准：%v\n", isEIP897)
}

func isEIP897Compliant(bytecode, abiJSON string) bool {
	// 解析 ABI
	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Fatalf("解析 ABI 时出错：%v", err)
	}

	// 解析字节码以查找必需的操作码
	bytecodeBytes := common.FromHex(bytecode)
	requiredOpcodes := []string{
		"60", // PUSH1
		"73", // PUSH4
		"80", // DUP1
		"90", // SWAP1
		"a1", // SSTORE
		"f3", // RETURN
		"ff", // SELFDESTRUCT
	}

	for _, op := range requiredOpcodes {
		if !containsOpcode(bytecodeBytes, op) {
			return false
		}
	}

	// 检查是否存在 'implementation' 函数
	_, found := contractABI.Methods["implementation"]
	if !found {
		return false
	}

	return true
}

func containsOpcode(bytecode []byte, opcode string) bool {
	opcodeBytes := common.FromHex(opcode)
	for i := 0; i < len(bytecode); i++ {
		if bytecode[i] == opcodeBytes[0] {
			if i+len(opcodeBytes) <= len(bytecode) && bytesEqual(bytecode[i:i+len(opcodeBytes)], opcodeBytes) {
				return true
			}
		}
	}
	return false
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
