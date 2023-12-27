package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

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
type SmsRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func main() {
	e := echo.New()
	e.POST("/v1/sms", func(c echo.Context) error {
		req := new(SmsWithTemplateRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		fmt.Printf("req: %+v\n", req)
		resp := new(SmsResponse)
		resp.MSGID = 1
		resp.Message = "halosis-success"
		resp.Status = "ok"
		return c.JSON(200, resp)
	})
	if err := e.Start(":8082"); err != nil {
		e.Logger.Fatal(err)
	}
}
