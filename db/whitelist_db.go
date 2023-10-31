package db

// BulkAddWhitelistAddr 批量增加白名单信息
//func BulkAddWhitelistAddr(db *sql.DB, addrs []domain.WhitelistAddr) error {
//	// 开始事务
//	tx, err := db.Begin()
//	if err != nil {
//		log.Println("BulkAddWhitelistAddr: Fail begin tx->", err.Error())
//		return err
//	}
//	// 准备插入语句
//	stmt, err := db.Prepare("INSERT INTO T_WHITELIST_ADDR (TWAR_KEY, CID, TW_ADDR, TW_CHAIN, TW_TYPE, ADD_TYPE, ADDR_ILL, ADDR_SOURCE, TAG_KEY,TOKEN_NAME, ABI,PROXY_ADDR,WEBSITE,CREATOR_ID, CREATE_DATE, MODIFIER_ID, LAST_MODIFY_DATE, VERSION,TOKEN_DECIMAL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?,?)")
//	if err != nil {
//		log.Println("BulkAddWhitelistAddr: Fail prepare sql->", err.Error())
//		return err
//	}
//	defer stmt.Close()
//	// 执行批量插入
//	for _, addr := range addrs {
//		//_, err = stmt.Exec(addr.TWARKey, addr.CID, addr.TWAddr, addr.TWChain, addr.TWType, addr.AddType, addr.AddrIll, addr.AddrSource, addr.TagKey, addr.TokenName, addr.Abi, addr.ProxyAddr, addr.Website, addr.CreatorID, addr.CreateDate, addr.ModifierID, addr.LastModifyDate, addr.Version, addr.TokenDecimal)
//		//if err != nil {
//		//	tx.Rollback()
//		//	log.Println("BulkAddWhitelistAddr: Fail exec insert sql->", err.Error())
//		//	return err
//		//}
//	}
//	// 提交事务
//	err = tx.Commit()
//	if err != nil {
//		log.Println("BulkAddWhitelistAddr: Fail commit tx->", err.Error())
//		return err
//	}
//	return nil
//}
//
//// GetAbiAndTokenByAddr 根据指定的合约地址查询abi信息以及token,token decimal信息
//func GetAbi(db *sql.DB, addr string) ([]byte, error) {
//	query := "SELECT ABI FROM T_WHITELIST_ADDR WHERE TW_ADDR = ?"
//	row := db.QueryRow(query, addr)
//	// 解析查询结果
//	var abi []byte
//	err := row.Scan(&abi)
//	if err != nil {
//		//若指定地址的信息不存在，返回空
//		if err == sql.ErrNoRows {
//			log.Printf("GetAbiInfoByAddr: TW_ADDR %s does not exist or abi not exist", addr)
//			return nil, nil
//		} else {
//			return nil, fmt.Errorf("GetAbiInfoByAddr:Failed to retrieve ABI-> %v", err.Error())
//		}
//	}
//	return abi, nil
//}
//
//// AddWhitelistAddr 添加白名单信息
//func AddWhitelistAddr(db *sql.DB, addr domain.WhitelistAddr) error {
//	stmt, err := db.Prepare("INSERT INTO T_WHITELIST_ADDR (TWAR_KEY, CID, TW_ADDR, TW_CHAIN, TW_TYPE, ADD_TYPE, ADDR_ILL, ADDR_SOURCE, TAG_KEY, CREATOR_ID, CREATE_DATE, MODIFIER_ID, LAST_MODIFY_DATE, VERSION) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	_, err = stmt.Exec(addr.TWARKey, addr.CID, addr.TWAddr, addr.TWChain, addr.TWType, addr.AddType, addr.AddrIll, addr.AddrSource, addr.TagKey, addr.CreatorID, addr.CreateDate, addr.ModifierID, addr.LastModifyDate, addr.Version)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// UpdateWhitelistAddr 修改白名单地址信息
//func UpdateWhitelistAddr(db *sql.DB, addr domain.WhitelistAddr) error {
//	stmt, err := db.Prepare("UPDATE T_WHITELIST_ADDR SET TW_ADDR = ?, TW_CHAIN = ?, TW_TYPE = ?, ADD_TYPE = ?, ADDR_ILL = ?, ADDR_SOURCE = ?, MODIFIER_ID = ?, LAST_MODIFY_DATE = ?, VERSION = ? WHERE TWAR_KEY = ?")
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	_, err = stmt.Exec(addr.TWAddr, addr.TWChain, addr.TWType, addr.AddType, addr.AddrIll, addr.AddrSource, addr.ModifierID, addr.LastModifyDate, addr.Version, addr.TWARKey)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// DeleteWhitelistAddr 删除白名单地址信息
//func deleteWhitelistAddr(db *sql.DB, twarKey string) error {
//	stmt, err := db.Prepare("DELETE FROM T_WHITELIST_ADDR WHERE TWAR_KEY = ?")
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	_, err = stmt.Exec(twarKey)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
