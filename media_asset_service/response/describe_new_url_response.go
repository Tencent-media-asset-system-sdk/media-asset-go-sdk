package response

// DescribeNewURLResponse
type DescribeNewURLResponse struct {
	Response struct {
		CommonResponse
		NewURLSet    []string `json:"NewURLSet"`
		FailedReason []string `json:"FailedReason"`
	} `json:"Response"`
}
