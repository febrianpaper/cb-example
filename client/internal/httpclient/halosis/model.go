package halosis

type SmsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	MSGID   int    `json:"msgid"`
}

type SmsComponent struct {
	Text string `json:"text"`
}

type SmsTemplateSetting struct {
	Name       string         `json:"name"`
	Components []SmsComponent `json:"components"`
}

type SmsWithTemplateRequest struct {
	To       string             `json:"to"`
	Template SmsTemplateSetting `json:"template"`
}
