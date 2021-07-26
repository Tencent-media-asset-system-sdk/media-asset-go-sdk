package request

// CommitUploadRequest
type CommitUploadRequest struct {
	CommonRequest
	MediaID  uint64   `json:"MediaID" validate:"required"`
	Key      string   `json:"Key" validate:"required"`
	Bucket   string   `json:"Bucket" validate:"required"`
	UploadId string   `json:"UploadId" validate:"required"`
	ETagSet  []string `json:"ETagSet"`
}
