// Package entity 风险名单及风险来源
package entity

//	type AdsDataSource struct {
//		DsName  string `json:"ds_name"`  //来源名称
//		DsAddr  string `json:"ds_addr"`  //来源地址
//		DsLevel uint   `json:"ds_level"` //监控层
//	}
type AdsDataSource struct {
	DsAddr string `json:"dsAddr"` //来源地址
	DsType string `json:"dsType"` //来源类型
	Number uint   `json:"number"` //标记次数
}
type WalletAddr struct {
	WaAddr      string          `json:"waAddr"`        //风险钱包地址
	WaRiskLevel uint            `json:"waRiskLevel"`   //风险层级
	WaChain     string          `json:"waChain"`       //所在链
	WaTicker    string          `json:"waTicker"`      //货币代码
	DsAddr      []AdsDataSource `json:"adsDataSource"` //来源地址
}
