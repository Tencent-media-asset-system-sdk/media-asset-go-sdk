package response

// UploadPartResponse
type UploadPartResponse struct {
	Response struct {
		CommonResponse
		ETag string `json:"ETag"`
	} `json:"Response"`
}

// PutObjectResponse
type PutObjectResponse struct {
	Response struct {
		CommonResponse
		Key  string `json:"Key"`
		ETag string `json:"ETag"`
	} `json:"Response"`
}
