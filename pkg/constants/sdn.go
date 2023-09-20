package constants

// SDN列表字段信息

const (
	GENDER       = "Gender"
	EMAIL        = "Email Address"
	WEBSITE      = "Website"
	PHONE_NUMBER = "Phone Number"
	ORGAN_TYPE   = "Organization Type:"
	ORGAN_DATE   = "Organization Established Date"
	OTHER_INFO1  = "Transactions Prohibited For Persons Owned or Controlled By U.S. Financial Institutions:"
	OTHER_INFO2  = "Additional Sanctions Information -"
	OTHER_INFO3  = "Secondary sanctions risk:"
	ID_COUNTRY   = "idCountry"
	EXPRI_DATE   = "expirationDate"
)
const (
	CONDITION = `^Digital Currency Address - ([\D]{3,16}$)` //SDN数字地址过滤条件
)
