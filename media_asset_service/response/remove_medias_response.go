package response

// RemoveMediasResponse
type RemoveMediasResponse struct {
	CommonResponse
	FailedMediaSet []*FailedMediaInfo `json:"FailedMediaSet"`
}
