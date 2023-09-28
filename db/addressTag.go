package db

import (
	"database/sql"
	"lyods-adsTool/domain/db"
)

// AddAddrTag 添加地址标签
func AddAddrTag(db *sql.DB, tag db.AddrTag) error {
	stmt, err := db.Prepare("INSERT INTO T_ADDR_TAG (TAG_KEY, CID, TAG_NAME, TAG_STATUS, TAG_ILL, CREATOR_ID, CREATE_DATE, MODIFIER_ID, LAST_MODIFY_DATE, VERSION) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.TagKey, tag.CID, tag.TagName, tag.TagStatus, tag.TagIll, tag.CreatorID, tag.CreateDate, tag.ModifierID, tag.LastModifyDate, tag.Version)
	if err != nil {
		return err
	}

	return nil
}

// UpdateAddrTag 修改标签信息
func UpdateAddrTag(db *sql.DB, tag db.AddrTag) error {
	stmt, err := db.Prepare("UPDATE T_ADDR_TAG SET TAG_NAME = ?, TAG_STATUS = ?, TAG_ILL = ?, MODIFIER_ID = ?, LAST_MODIFY_DATE = ?, VERSION = ? WHERE TAG_KEY = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.TagName, tag.TagStatus, tag.TagIll, tag.ModifierID, tag.LastModifyDate, tag.Version, tag.TagKey)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAddrTag 删除标签信息
func DeleteAddrTag(db *sql.DB, tagKey string) error {
	stmt, err := db.Prepare("DELETE FROM T_ADDR_TAG WHERE TAG_KEY = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tagKey)
	if err != nil {
		return err
	}

	return nil
}

// GetAddrTags 获取地址标签信息
func GetAddrTags(db *sql.DB) ([]db.AddrTag, error) {
	rows, err := db.Query("SELECT * FROM T_ADDR_TAG")
	if err != nil {
		return nil, err
	}
	var tags []db.AddrTag
	for rows.Next() {
		var tag db.AddrTag
		err := rows.Scan(&tag.TagKey, &tag.CID, &tag.TagName, &tag.TagStatus, &tag.TagIll, &tag.CreatorID, &tag.CreateDate, &tag.ModifierID, &tag.LastModifyDate, &tag.Version)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
