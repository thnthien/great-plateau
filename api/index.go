package api

import (
	"github.com/thnthien/great-plateau/api/fiber/response"

	cerrors "github.com/thnthien/great-plateau/errors"
)

type HTTPErrorResponse struct {
	Status     string         `json:"status"`
	Code       cerrors.Code   `json:"code"`
	Message    string         `json:"message"`
	DevMessage any            `json:"dev_message"`
	Errors     map[string]any `json:"errors"`
	RID        string         `json:"rid"`
}

// Handler ...
type Handler struct {
}

// Resp ...
func (Handler) Resp() response.IResponse {
	return response.NewResponse()
}

type IHealthController interface {
	SetReady(b bool)
}
