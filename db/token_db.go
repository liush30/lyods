package db

import (
	"database/sql"
	"fmt"
	"lyods-adsTool/domain"
	"time"
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

// ExistsToken 判断合约地址是否是存储于数据库中的erc20合约地址
func ExistsToken(db *sql.DB, contractAddress, blockchain string) (bool, error) {
	querySQL := `
        SELECT COUNT(*)
        FROM t_token
        WHERE CONTRACT_ADDRESS = ? AND CHAIN = ?
    `
	var count int
	err := db.QueryRow(querySQL, contractAddress, blockchain).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// QueryAbiAndCheckByAddress 根据合约地址查询abi信息,以及记录count数
func QueryAbiAndCheckByAddress(db *sql.DB, contractAddress, chain string) (string, int, error) {
	querySQL := `SELECT abi,DECIMALS  FROM t_token
WHERE contract_address = ? AND chain = ?;`
	var abi string
	var decimals int
	err := db.QueryRow(querySQL, contractAddress, chain).Scan(&abi, &decimals)
	if err != nil {
		return "", 0, fmt.Errorf("query from t_token abi info error: %v", err)
	}
	return abi, decimals, nil
}

func GetContractAddressAll(db *sql.DB, chain string) ([]string, error) {
	var contractAddresses []string

	// 使用占位符来构建 SQL 查询，以避免 SQL 注入风险
	var sql string
	if chain != "" {
		sql = "SELECT CONTRACT_ADDRESS FROM t_token WHERE CHAIN = ?  and ABI is  null"
	} else {
		sql = "SELECT CONTRACT_ADDRESS FROM t_token"
	}

	// 查询数据库
	rows, err := db.Query(sql, chain)
	if err != nil {
		return nil, fmt.Errorf("query from t_token error: %v", err)
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var contractAddress string

		// 扫描查询结果中的数据
		if err := rows.Scan(&contractAddress); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		contractAddresses = append(contractAddresses, contractAddress)
	}

	// 检查是否有查询错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query from t_token error: %v", err)
	}
	return contractAddresses, nil
}

func SaveContractABI(db *sql.DB, chain, address, abiInfo, proxy string) error {
	updateSQL := "update t_token set ABI=?,PROXY_ADDR=?,LAST_MODIFY_DATE=? where CONTRACT_ADDRESS=? and CHAIN=?"
	timeStr := time.DateTime
	result, err := db.Exec(updateSQL, abiInfo, proxy, timeStr, address, chain)
	if err != nil {
		return fmt.Errorf("update to t_token error :%v\n", err.Error())
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error :%v\n", err.Error())
	} else if count < 1 {
		return fmt.Errorf("update to t_token error :%v\n", "affected 0 rows")

	}
	return nil
}
