package paging

type Request struct {
	PageIndex int `json:"index"`
	PageSize  int `json:"size"`
}
