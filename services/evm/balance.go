package evm

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"lyods-adsTool/pkg/constants"
	"math/big"
	"strings"
)

// GetBalanceChange 获取指定地址，在交易(指定区块)前后原生代币金额的变化
func (e *EthClient) GetBalanceChange(blockNumber *big.Int, address string) (*big.Int, error) {
	// 获取前一个区块号
	beforeNumber := GetLastBlockNumber(blockNumber)

	// 查询地址在交易前的余额
	//balanceBefore, err := e.BalanceAt(context.Background(), common.HexToAddress(address), beforeNumber)
	//if err != nil {
	//	log.Println("Fail get before balance:", err.Error())
	//	return nil, err
	//}
	balanceBefore, err := e.GetBalance(address, beforeNumber)
	if err != nil {
		log.Println("Fail get before balance:", err.Error())
		return nil, err
	}
	balanceAfter, err := e.GetBalance(address, blockNumber)
	// 查询地址在交易后的余额
	//balanceAfter, err := e.BalanceAt(context.Background(), common.HexToAddress(address), blockNumber)
	if err != nil {
		log.Println("Fail get after balance:", err.Error())
		return nil, err
	}

	// 计算余额变化
	balanceChange := new(big.Int).Sub(balanceAfter, balanceBefore)
	return balanceChange, nil
}

// GetERC20TokenBalance 查询指定账户在指定 ERC-20 代币合约中的余额
func (e *EthClient) GetERC20TokenBalance(tokenContractAddress string, accountAddress string, blockNumber *big.Int) (*big.Int, error) {
	// 创建 ERC-20 代币合约的 ABI
	tokenABI, err := abi.JSON(strings.NewReader(constants.ABI_ERC20))
	if err != nil {
		return nil, err
	}
	// 构造 balanceOf 函数调用
	balanceCallData, err := tokenABI.Pack("balanceOf", common.HexToAddress(accountAddress))
	if err != nil {
		return nil, err
	}
	contractAddress := common.HexToAddress(tokenContractAddress)
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: balanceCallData,
	}
	result, err := e.CallContract(context.Background(), callMsg, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	return new(big.Int).SetBytes(result), nil
}
func (e *EthClient) GetBalance(address string, blockNumber *big.Int) (*big.Int, error) {
	// 查询地址在交易前的余额
	balanceBefore, err := e.BalanceAt(context.Background(), common.HexToAddress(address), blockNumber)
	if err != nil {
		log.Println("Fail get before balance:", err.Error())
		return nil, err
	}
	return balanceBefore, nil
}

// GetERC20TokenBalanceChange 查询指定账户在指定 ERC-20 代币合约中的余额变化
func (e *EthClient) GetERC20TokenBalanceChange(tokenContractAddress string, accountAddress string, blockNumber *big.Int) (*big.Int, error) {
	// 获取前一个区块号
	beforeNumber := GetLastBlockNumber(blockNumber)
	//查询在blockNumber前账户的余额
	balanceBefore, err := e.GetERC20TokenBalance(tokenContractAddress, accountAddress, beforeNumber)
	if err != nil {
		log.Println("Fail get before balance:", err.Error())
		return nil, err
	}
	//查询在blockNumber后
	//
	//账户的余额
	balanceAfter, err := e.GetERC20TokenBalance(tokenContractAddress, accountAddress, blockNumber)
	if err != nil {
		log.Println("Fail get after balance:", err.Error())
		return nil, err
	}
	// 计算余额变化
	balanceChange := new(big.Int).Sub(balanceAfter, balanceBefore)
	return balanceChange, nil
}

// GetLatestBlockNumber 获得最新区块号
func (e *EthClient) GetLatestBlockNumber() (int64, error) {
	// 获取最新的区块头
	header, err := e.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Int64(), nil
}
