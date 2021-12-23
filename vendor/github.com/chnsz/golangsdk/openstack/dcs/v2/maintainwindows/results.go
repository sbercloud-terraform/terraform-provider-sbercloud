package maintainwindows

import (
	"github.com/chnsz/golangsdk"
)

// GetResponse response
type GetResponse struct {
	MaintainWindows []MaintainWindow `json:"maintain_windows"`
}

// MaintainWindow for dcs
type MaintainWindow struct {
	ID      int    `json:"seq"`
	Begin   string `json:"begin"`
	End     string `json:"end"`
	Default bool   `json:"default"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*GetResponse, error) {
	var s GetResponse
	err := r.Result.ExtractInto(&s)
	return &s, err
}
