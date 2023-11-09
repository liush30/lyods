package evm

import (
	"database/sql"
)

// AddABIRecords 增加表中还未存储abi的合约地址
func (e *EthClient) AddABIRecords(db *sql.DB) error {
	//var err error
	//// 检查TOKEN_NAME是否不为空，并获取TWAR_KEY、PROXY_ADDR和TW_ADDR信息
	//rows, err := db.Query("SELECT TWAR_KEY,PROXY_ADDR, TW_ADDR FROM T_WHITELIST_ADDR WHERE TOKEN_NAME IS NOT NULL && ABI IS NULL")
	//if err != nil {
	//	return fmt.Errorf("AddABIRecords:Fail query->%v", err)
	//}
	//defer rows.Close()
	//
	////为TOKEN_NAME不为空的记录添加ABI信息
	//for rows.Next() {
	//	var key, twAddr, abiStr, addr string
	//	var proxyAddr sql.NullString
	//	if err := rows.Scan(&key, &proxyAddr, &twAddr); err != nil {
	//		return fmt.Errorf("AddABIRecords Fail scan row ->%v", err.Error())
	//	}
	//	//若proxyAddr不为空,查询proxyAddr的abi信息，否则直接查询合约地址信息
	//	if proxyAddr.Valid {
	//		addr = proxyAddr.String
	//	} else {
	//		addr = twAddr
	//	}
	//	abiStr, err = utils.GetContractAbi(addr)
	//	if err != nil {
	//		return fmt.Errorf("AddABIRecords:Fail get contract abi->%v", err)
	//	}
	//	// 将ABI记录插入表中
	//	_, err := db.Exec("UPDATE T_WHITELIST_ADDR SET ABI = ? WHERE TWAR_KEY = ? ", []byte(abiStr), key)
	//	if err != nil {
	//		return fmt.Errorf("AddABIRecords:Fail insert %s abi info->%v\n", addr, err)
	//	}
	//}
	//if err := rows.Err(); err != nil {
	//	return fmt.Errorf("AddABIRecords:Fail traverse rows->%v", err)
	//}
	return nil
}
