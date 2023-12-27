package echo

import (
	"context"
	"fbriansyah/client/internal/usecase/commands"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Server) SmsSendHandler(c echo.Context) error {
	var req struct {
		PhoneNumber  string   `json:"phone_number"`
		Message      string   `json:"message"`
		TemplateName string   `json:"template_name"`
		TemplateData []string `json:"template_data"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.uc.SendSms(ctx, commands.SendSmsParams{
		To:           req.PhoneNumber,
		Message:      req.Message,
		TemplateName: req.TemplateName,
		TemplateData: req.TemplateData,
	})

	var response struct {
		IsError bool   `json:"is_error"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	if err != nil {
		response.IsError = true
		response.Message = err.Error()
		response.Status = http.StatusInternalServerError
	} else {
		response.IsError = false
		response.Message = "success"
		response.Status = http.StatusOK
	}

	return c.JSON(response.Status, response)
}
