package entity

type AdsDataSource struct {
	DsName  string `json:"ds_name"`  //来源名称
	DsAddr  string `json:"ds_addr"`  //来源地址
	DsLevel int    `json:"ds_level"` //监控层次
}
type WalletAddr struct {
	WaAddr      string `json:"wa_addr"`       //风险钱包地址
	WaRiskLevel int    `json:"wa_risk_level"` //风险层级
	WaChain     string `json:"wa_chain"`      //所在链
}
