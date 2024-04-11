package paging

import "errors"

type Request struct {
	PageIndex int `json:"index"`
	PageSize  int `json:"size"`
}

func (req *Request) Validate() error {
	if req.PageIndex < 0 {
		return errors.New("invalid page index")
	}
	if req.PageSize < 0 {
		return errors.New("invalid page size")
	}
	return nil
}

func (req *Request) ApplyDefaults() {
	if req.PageIndex == 0 {
		req.PageIndex = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	return
}
