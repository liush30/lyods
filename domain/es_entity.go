package domain

//es存储结构-实体信息 risk-domain

type DateOfBirth struct {
	DateOfBirth string `json:"dateOfBirth"`
	MainEntry   bool   `json:"mainEntry"`
}

type PlaceOfBirth struct {
	PlaceOfBirth string `json:"placeOfBirth"`
	MainEntry    bool   `json:"mainEntry"`
}

type Nationality struct {
	Country   string `json:"country"`
	MainEntry bool   `json:"mainEntry"`
}

type ID struct {
	IDType         string `json:"idType"`   //ID列表
	IDNumber       string `json:"idNumber"` //
	IDCountry      string `json:"idCountry"`
	ExpirationDate string `json:"expirationDate"`
}

type OtherInfo struct {
	Type string `json:"type"` //类型
	Info string `json:"info"` //信息
}

type AddressList struct {
	Country         string   `json:"country"`
	StateOrProvince string   `json:"stateOrProvince"`
	City            string   `json:"city"`
	Other           []string `json:"other"`
}

type Entity struct {
	IsIndividual     bool           `json:"isIndividual"`     //是否为个体
	Name             string         `json:"name"`             //名字
	AkaList          []string       `json:"akaList"`          //别名列表
	AddressList      []AddressList  `json:"addressList"`      //地址列表
	DateOfBirthList  []DateOfBirth  `json:"dateOfBirthList"`  //出生日期列表
	PlaceOfBirth     []PlaceOfBirth `json:"placeOfBirth"`     //出生地址列表
	Gender           string         `json:"gender"`           //性别
	Email            []string       `json:"emailList"`        //邮箱列表
	Website          []string       `json:"websiteList"`      //网站列表
	PhoneNumber      []string       `json:"phoneNumberList"`  //电话号码
	IDList           []ID           `json:"idList"`           //ID列表信息
	NationalityList  []Nationality  `json:"nationalityList"`  //国籍列表
	OrganizationType string         `json:"organizationType"` //机构类型
	CitizenshipList  []Nationality  `json:"citizenshipList"`  //公民身份列表
	OrgEstDate       string         `json:"orgEstDate"`       //机构成立日期
	OtherInfo        []OtherInfo    `json:"otherInfo"`        //其他信息
}
