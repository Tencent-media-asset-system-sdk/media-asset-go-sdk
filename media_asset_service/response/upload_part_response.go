package response

// UploadPartResponse
type UploadPartResponse struct {
	Response struct {
		CommonResponse
		ETag string `json:"ETag"`
	} `json:"Response"`
}
