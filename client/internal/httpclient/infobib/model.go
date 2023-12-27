package infobib

type SmsDestination struct {
	To string `json:"to"`
}

type SmsRequestMessage struct {
	Destinations []SmsDestination `json:"destination"`
	From         string           `json:"from"`
	Text         string           `json:"text"`
}

type SmsRequest struct {
	Messages []SmsRequestMessage `json:"message"`
}

func newSmsRequest(to string, message string) SmsRequest {
	return SmsRequest{
		Messages: []SmsRequestMessage{
			{
				Destinations: []SmsDestination{
					{
						To: to,
					},
				},
				From: "paper.id",
				Text: message,
			},
		},
	}
}

type SmsResponseStatus struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupId     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
}

type SmsMessageResponse struct {
	MessageId string            `json:"messageId"`
	Status    SmsResponseStatus `json:"status"`
	To        string            `json:"to"`
}

type SmsResponse struct {
	BulkId   string               `json:"bulkId"`
	Messages []SmsMessageResponse `json:"messages"`
}
