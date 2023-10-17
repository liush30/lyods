package domain

// Token 结构体
type Token struct {
	TokenKey        string
	ContractAddress string
	Symbol          string
	Decimals        int
	Blockchain      string
	CreateDate      string
	LastModifyDate  string
	Abi             []byte
	ProxyAddr       string
}
