package handler

import "github.com/thnthien/great-plateau/api/fiber/response"

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
