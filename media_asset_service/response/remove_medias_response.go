package response

// RemoveMediasResponse
type RemoveMediasResponse struct {
	Response struct {
		CommonResponse
		FailedMediaSet []*FailedMediaInfo `json:"FailedMediaSet"`
	} `json:"Response"`
}
