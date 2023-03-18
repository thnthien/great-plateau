package fiber

import (
	"unsafe"

	"github.com/gofiber/fiber/v2"

	"github.com/thnthien/great-plateau/api"
	cerrors "github.com/thnthien/great-plateau/errors"
)

func getString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ErrorHandler(env string, httpStatusMappingFunc func(code cerrors.Code) int) func(ctx *fiber.Ctx, err error) error {
	mappingFunc := httpStatusMappingFunc
	if mappingFunc == nil {
		mappingFunc = api.DefaultStatusMapping
	}
	return func(ctx *fiber.Ctx, err error) error {
		// Statuscode defaults to 500
		code := fiber.StatusInternalServerError

		rid := getString(ctx.Context().Response.Header.Peek(fiber.HeaderXRequestID))

		devMsg := err
		if env != "D" {
			devMsg = nil
		}

		if e, ok := err.(*fiber.Error); ok {
			errCode := cerrors.Code(e.Code)
			return ctx.Status(e.Code).JSON(api.HTTPResponse{
				Status:     errCode.String(),
				Code:       errCode,
				Message:    errCode.String(),
				DevMessage: devMsg,
				Errors:     nil,
				RID:        rid,
			})
		}

		clientError, ok := err.(*cerrors.CError)
		if !ok {
			return ctx.Status(code).JSON(api.HTTPResponse{
				Status:     cerrors.InternalServerError.String(),
				Code:       cerrors.InternalServerError,
				Message:    cerrors.InternalServerError.String(),
				DevMessage: devMsg,
				Errors:     nil,
				RID:        rid,
			})
		}

		if env == "D" {
			devMsg = clientError
		}

		code = mappingFunc(clientError.Code)
		return ctx.Status(code).JSON(api.HTTPResponse{
			Status:     clientError.Code.String(),
			Code:       clientError.Code,
			Message:    clientError.Message,
			DevMessage: devMsg,
			Errors:     nil,
			RID:        rid,
		})
	}
}
