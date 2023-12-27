package domain

type VendorSetting struct {
	BaseUrl           string `json:"base_url"`
	IsSupportTemplate bool   `json:"is_support_template"`
	AllowSms          bool   `json:"allow_sms"`
	SmsEndpoint       string `json:"sms_endpoint"`
	AllowWa           bool   `json:"allow_wa"`
	WAEndpoint        string `json:"wa_endpoint"`
}

func (st *VendorSetting) GetSmsUrl() string {
	return st.BaseUrl + st.SmsEndpoint
}

func (st *VendorSetting) GetWAUrl() string {
	return st.BaseUrl + st.WAEndpoint
}
