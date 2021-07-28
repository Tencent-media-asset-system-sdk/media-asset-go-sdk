package response

// DescribeMediasResponse
type DescribeMediasResponse struct {
	Response struct {
		CommonResponse
		MediaInfoSet []*MediaInfo `json:"MediaInfoSet"`
		TotalCount   int32        `json:"TotalCount"`
	} `json:"Response"`
}
