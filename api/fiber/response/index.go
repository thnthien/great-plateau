package response

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"

	"github.com/thnthien/great-plateau/api"
	cerrors "github.com/thnthien/great-plateau/errors"
)

// IResponse ...
type IResponse interface {
	WithData(data any) IResponse
	WithPaging(data any) IResponse
	WithMessage(data string) IResponse
	WithStatus(status int) IResponse
	WithCode(code cerrors.Code) IResponse
	Json(c *fiber.Ctx) error
	NoContent(c *fiber.Ctx) error
}

type response struct {
	status  int
	code    cerrors.Code
	data    interface{}
	link    *api.Links
	message string
}

// NewResponse ...
func NewResponse() IResponse {
	return response{}
}

// WithData ...
func (r response) WithData(data interface{}) IResponse {
	r.data = data
	return r
}

// WithMessage ...
func (r response) WithMessage(data string) IResponse {
	r.message = data
	return r
}

// WithPaging ...
func (r response) WithPaging(data interface{}) IResponse {
	if r.link == nil {
		r.link = &api.Links{}
	}
	err := copier.Copy(r.link, data)
	if err != nil {
		fmt.Printf("[WithPaging] failed to copy paging data [%+v] [%v]", data, err)
	}
	return r
}

// WithStatus ...
func (r response) WithStatus(status int) IResponse {
	r.status = status
	return r
}

func (r response) WithCode(code cerrors.Code) IResponse {
	r.code = code
	return r
}

// Json ...
func (r response) Json(c *fiber.Ctx) error {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	if r.code == 0 {
		r.code = cerrors.Code(r.status)
	}
	return c.Status(r.status).JSON(api.HTTPResponse{
		Status:  r.code.String(),
		Code:    r.code,
		Data:    r.data,
		Link:    r.link,
		Message: r.message,
	})
}

// NoContent ...
func (r response) NoContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

var _ IResponse = response{}
