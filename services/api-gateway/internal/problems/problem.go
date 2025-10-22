package problems

import (
	"fmt"
	"net/http"
)

type Problem struct {
	Type           string `json:"type,omitempty"`
	Title          string `json:"title"`
	Status         int    `json:"status"`
	Detail         string `json:"detail,omitempty"`
	Instance       string `json:"instance,omitempty"`
	InternalDetail string `json:"-"`
}

func (p *Problem) Error() string {
	return fmt.Sprintf("%s: %s", p.Title, p.Detail)
}

func NewInternal(detail, InternalMsg string) *Problem {
	return &Problem{
		Type:           "https://example.com/probs/internal",
		Title:          "Internal Server Error",
		Status:         http.StatusInternalServerError,
		Detail:         detail,      // safe for clients
		InternalDetail: InternalMsg, // detailed for logs
	}
}
