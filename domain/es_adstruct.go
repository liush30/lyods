package domain

import "time"

// es存储结构-风险名单及风险来源 risk-address

type WalletAddr struct {
	WaAddr      string          `json:"waAddr"`        //风险钱包地址
	EntityId    string          `json:"entityId"`      //entityID
	WaRiskLevel uint            `json:"waRiskLevel"`   //最高风险层级
	WaChain     string          `json:"waChain"`       //所在链
	DsAddr      []AdsDataSource `json:"adsDataSource"` //来源地址
	LevelNumber []Level         `json:"levelNumber"`   //被标记层级信息
	Rules       []string        `json:"rules"`         //被中规则id列表
	IsTrace     bool            `json:"isTrace"`       //是否追查子交易
	IsNeedTrace bool            `json:"isNeedTrace"`   //是否需要追查子交易
}
type AdsDataSource struct {
	DsAddr      string    `json:"dsAddr"`      //来源地址
	DsTransHash []string  `json:"dsTransHash"` //涉及风险交易哈希列表
	DsType      string    `json:"dsType"`      //涉及风险交易哈希列表
	Illustrate  string    `json:"illustrate"`  //风险说明
	Time        time.Time `json:"time"`        //被标记时间
	DsRules     []string  `json:"dsRules"`     //规则id
}
type Level struct {
	Level  int16 `json:"level"`  //所在层级
	Number int16 `json:"number"` //被标记次数
}
