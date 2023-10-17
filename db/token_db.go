package db

import (
	"database/sql"
	"fmt"
	"lyods-adsTool/domain"
)

// AddToken 添加数据到表
func AddToken(db *sql.DB, token domain.Token) (int64, error) {
	insertSQL := `
        INSERT INTO t_token (TOKEN_KEY, CONTRACT_ADDRESS, SYMBOL, DECIMALS, BLOCKCHAIN, CREATE_DATE,LAST_MODIFY_DATE)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	result, err := db.Exec(insertSQL, token.TokenKey, token.ContractAddress, token.Symbol, token.Decimals, token.Blockchain, token.CreateDate, token.LastModifyDate)
	if err != nil {
		return 0, fmt.Errorf("insert to database error :%v\n", err.Error())
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected error :%v\n", err.Error())
	}
	return count, nil
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
