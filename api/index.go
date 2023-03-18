package api

import (
	cerrors "github.com/thnthien/great-plateau/errors"
)

type HTTPResponse struct {
	Status     string         `json:"status,omitempty"`
	Code       cerrors.Code   `json:"code,omitempty"`
	Message    string         `json:"message,omitempty"`
	DevMessage any            `json:"dev_message,omitempty"`
	Errors     map[string]any `json:"errors,omitempty"`
	RID        string         `json:"rid,omitempty"`
	Data       any            `json:"data,omitempty"`
	Link       *Links         `json:"link,omitempty"`
}

type Links struct {
	BeforeCount int      `json:"before_count,omitempty"`
	AfterCount  int      `json:"after_count,omitempty"`
	Count       int      `json:"count,omitempty"`
	Cursors     *Cursors `json:"cursors,omitempty"`
	Next        string   `json:"next,omitempty"`
	Prev        string   `json:"prev,omitempty"`
}

type Cursors struct {
	After  string `json:"after"`
	Before string `json:"before"`
}
