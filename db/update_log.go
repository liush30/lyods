package db

import (
	"database/sql"
	"fmt"
	"lyods-adsTool/domain"
)

// AddUpdateRecord 增加记录
func AddUpdateRecord(db *sql.DB, log domain.UpdateLog) error {
	query := "INSERT INTO t_update_log (LOG_KEY, UPDATE_DATE, UPDATE_NAME, UPDATE_RECORD) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, log.LogKey, log.UpdateDate, log.UpdateName, log.UpdateRecord)
	if err != nil {
		return fmt.Errorf("failed to add update record:%v", err)
	}
	return nil
}

// DeleteUpdateRecord 删除记录
func DeleteUpdateRecord(db *sql.DB, logKey string) error {
	query := "DELETE FROM t_update_log WHERE LOG_KEY = ?"
	result, err := db.Exec(query, logKey)
	if err != nil {
		return err
	}
	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return fmt.Errorf("no records deleted for LOG_KEY: %s", logKey)
	}
	return nil
}

// QueryAllUpdateRecords 查询所有记录
func QueryAllUpdateRecords(db *sql.DB) ([]domain.UpdateLog, error) {
	var updateLogs []domain.UpdateLog
	query := "SELECT LOG_KEY, UPDATE_DATE, UPDATE_NAME, UPDATE_RECORD FROM t_update_log"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query update records:%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var updateLog domain.UpdateLog
		if err := rows.Scan(&updateLog.LogKey, &updateLog.UpdateDate, &updateLog.UpdateName, &updateLog.UpdateRecord); err != nil {
			return nil, fmt.Errorf("failed to scan update record:%v", err)
		}
		updateLogs = append(updateLogs, updateLog)
	}
	return updateLogs, nil
}

// GetLatestUpdateRecord 根据 update_name 获取最新的更新记录
func GetLatestUpdateRecord(db *sql.DB, updateName string) (domain.UpdateLog, error) {
	query := "SELECT LOG_KEY, UPDATE_DATE, UPDATE_NAME, UPDATE_RECORD FROM t_update_log WHERE UPDATE_NAME = ? ORDER BY UPDATE_DATE DESC LIMIT 1"
	row := db.QueryRow(query, updateName)

	var latestUpdate domain.UpdateLog
	if err := row.Scan(&latestUpdate.LogKey, &latestUpdate.UpdateDate, &latestUpdate.UpdateName, &latestUpdate.UpdateRecord); err != nil {
		return domain.UpdateLog{}, fmt.Errorf("failed to scan update record:%v", err)
	}
	return latestUpdate, nil
}
