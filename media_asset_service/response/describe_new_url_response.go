package response

// DescribeNewURLResponse
type DescribeNewURLResponse struct {
	Response struct {
		CommonResponse
		NewURLSet    []string `json:"UploadMediaInfoSet"`
		FailedReason []string `json:"FailedReason"`
	} `json:"Response"`
}
