package response

// ApplyUploadResponse
type ApplyUploadResponse struct {
	Response struct {
		CommonResponse
		MediaID  string `json:"MediaID"`
		Key      string `json:"Key"`
		Bucket   string `json:"Bucket"`
		UploadId string `json:"UploadId"`
	} `json:"Response"`
}
