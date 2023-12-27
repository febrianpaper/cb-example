package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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

func main() {
	e := echo.New()
	e.POST("/v1/sms", func(c echo.Context) error {
		req := new(SmsRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		fmt.Printf("req: %+v\n", req)
		resp := new(SmsResponse)
		resp.BulkId = uuid.New().String()
		resp.Messages = make([]SmsMessageResponse, 1)
		resp.Messages[0].MessageId = uuid.New().String()
		resp.Messages[0].Status.Description = "success"
		resp.Messages[0].Status.GroupId = 1
		resp.Messages[0].Status.GroupName = "infobib-success"
		resp.Messages[0].Status.ID = 1
		resp.Messages[0].Status.Name = "success"
		resp.Messages[0].To = req.Messages[0].Destinations[0].To

		// send sms
		return c.JSON(http.StatusOK, resp)
	})
	if err := e.Start(":8081"); err != nil {
		e.Logger.Fatal(err)
	}
}
