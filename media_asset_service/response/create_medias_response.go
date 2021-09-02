package response

// CreateMediasResponse
type CreateMediasResponse struct {
	Response struct {
		CommonResponse
		UploadMediaInfoSet []UploadMediaInfo `json:"UploadMediaInfoSet"`
	} `json:"Response"`
}

// UploadMediaInfo
type UploadMediaInfo struct {
	FailedReason string `json:"FailedReason"`
	MediaID      uint64 `json:"MediaID"`
}
