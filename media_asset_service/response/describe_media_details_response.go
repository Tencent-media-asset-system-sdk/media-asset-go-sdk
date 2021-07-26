package response

// DescribeMediaDetailsResponse
type DescribeMediaDetailsResponse struct {
	CommonResponse
	MediaInfoSet []*MediaInfo `json:"MediaInfoSet"`
}
