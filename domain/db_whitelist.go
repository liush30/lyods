package domain

import "database/sql"

type AddrTag struct {
	TagKey         string
	CID            uint
	TagName        string
	TagStatus      string
	TagIll         string
	CreatorID      string
	CreateDate     string
	ModifierID     string
	LastModifyDate string
	Version        int
}

type WhitelistAddr struct {
	TWARKey        string
	CID            uint
	TWAddr         string
	TWChain        string
	TWType         string
	AddType        string
	AddrIll        string
	AddrSource     string
	TagKey         string
	TokenName      string
	Abi            []byte
	ProxyAddr      sql.NullString
	Website        string
	CreatorID      string
	CreateDate     string
	ModifierID     string
	LastModifyDate string
	Version        uint
	TokenDecimal   int
}
type EventInfo struct {
	TEKey        string
	EventSign    string
	SignString   string
	EventName    string
	ContractAddr string
	IsAnonymous  string
	Nature       string
}
