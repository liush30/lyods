package db

import (
	"database/sql"
	"log"
	"lyods-adsTool/domain/db"
)

// AddToken 添加数据到表
func AddToken(db *sql.DB, token db.Token) {
	insertSQL := `
        INSERT INTO t_token (TOKEN_KEY, CONTRACT_ADDRESS, SYMBOL, DECIMALS, BLOCKCHAIN, ABI,PROXY_ADDR, CREATE_DATE)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(insertSQL, token.TokenKey, token.ContractAddress, token.Symbol, token.Decimals, token.Blockchain, token.Abi, token.ProxyAddr, token.CreateDate)
	if err != nil {
		log.Fatal(err)
	}
}

// 判断合约地址是否为erc20合约地址
func existsToken(db *sql.DB, contractAddress, blockchain string) (bool, error) {
	querySQL := `
        SELECT COUNT(*)
        FROM t_token
        WHERE CONTRACT_ADDRESS = ? AND BLOCKCHAIN = ?
    `
	var count int
	err := db.QueryRow(querySQL, contractAddress, blockchain).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTokenByContractAddress 根据CONTRACT_ADDRESS查询信息
func GetTokenByContractAddress(db *sql.DB, contractAddress string) ([]byte, error) {
	return nil, nil
}
