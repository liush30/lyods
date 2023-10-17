package pre

import (
	"encoding/csv"
	"fmt"
	"io"
	"lyods-adsTool/db"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/utils"
	"os"
	"strconv"
	"time"
)

func GetContractAddressToDb(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	_, err = csvReader.Read() // 跳过文件头
	if err != nil && err != io.EOF {
		return err
	}

	dbClient := db.GetDb()
	defer dbClient.Close()

	tx, err := dbClient.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin a database transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO t_token (TOKEN_KEY, CONTRACT_ADDRESS, SYMBOL, DECIMALS, BLOCKCHAIN, CREATE_DATE, LAST_MODIFY_DATE) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	const batchSize = 100
	valueArgs := make([]interface{}, 0, batchSize*7)
	query := "INSERT INTO t_token (TOKEN_KEY, CONTRACT_ADDRESS, SYMBOL, DECIMALS, BLOCKCHAIN, CREATE_DATE, LAST_MODIFY_DATE) VALUES"

	for {
		// 读取一行
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV: %v", err)
		}

		decimals, err := strconv.Atoi(record[3])
		if err != nil {
			return fmt.Errorf("failed to convert decimals: %v", err)
		}

		token := domain.Token{
			TokenKey:        utils.GenerateUuid(),
			Blockchain:      record[0],
			ContractAddress: record[1],
			Symbol:          record[2],
			Decimals:        decimals,
			CreateDate:      time.Now().Format(time.RFC3339),
			LastModifyDate:  time.Now().Format(time.RFC3339),
		}

		// 生成批量插入的查询语句
		query += "(?, ?, ?, ?, ?, ?, ?),"
		valueArgs = append(valueArgs, token.TokenKey, token.ContractAddress, token.Symbol, token.Decimals, token.Blockchain, token.CreateDate, token.LastModifyDate)

		if len(valueArgs) >= batchSize*7 {
			// 移除最后的逗号
			query = query[:len(query)-1]

			_, err := tx.Exec(query, valueArgs...)
			if err != nil {
				return fmt.Errorf("failed to insert tokens into database: %v", err)
			}

			// 清空切片信息
			valueArgs = valueArgs[:0]
			query = "INSERT INTO t_token (TOKEN_KEY, CONTRACT_ADDRESS, SYMBOL, DECIMALS, BLOCKCHAIN, CREATE_DATE, LAST_MODIFY_DATE) VALUES"
		}
	}

	if len(valueArgs) > 0 {
		// 移除最后的逗号
		query = query[:len(query)-1]

		_, err := tx.Exec(query, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to insert remaining tokens into database: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the database transaction: %v", err)
	}

	return nil
}
