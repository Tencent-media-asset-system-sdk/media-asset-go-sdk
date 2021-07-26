package response

// DescribeMediasResponse
type DescribeMediasResponse struct {
	CommonResponse
	MediaInfoSet []*MediaInfo `json:"MediaInfoSet"`
	TotalCount   int32        `json:"TotalCount"`
}
