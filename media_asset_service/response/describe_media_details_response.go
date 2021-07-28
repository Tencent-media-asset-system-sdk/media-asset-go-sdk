package response

// DescribeMediaDetailsResponse
type DescribeMediaDetailsResponse struct {
	Response struct {
		CommonResponse
		MediaInfoSet []*MediaInfo `json:"MediaInfoSet"`
	} `json:"Response"`
}
