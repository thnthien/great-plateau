package api

import (
	"net/http"

	cerrors "github.com/thnthien/great-plateau/errors"
)

func DefaultStatusMapping(code cerrors.Code) int {
	if code > 999 {
		c := code
		for c > 999 {
			c /= 10
		}
		return int(c)
	}
	if code >= 100 {
		return int(code)
	}
	return http.StatusInternalServerError
}
