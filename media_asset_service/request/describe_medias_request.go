package request

// DescribeMediasRequest
type DescribeMediasRequest struct {
	CommonRequest
	PageNumber int32    `json:"PageNumber"`
	PageSize   int32    `json:"PageSize"`
	Inner      bool     `json:"Inner"`
	FilterBy   FilterBy `json:"FilterBy"`
}
