package request

// CreateMediasRequest
type CreateMediasRequest struct {
	CommonRequest
	UploadMediaSet []UploadMedia `json:"UploadMediaSet" validate:"required,dive"`
	Inner          bool          `json:"Inner"`
}

// UploadMedia
type UploadMedia struct {
	Name       string    `json:"Name" validate:"required"`
	LocalPath  string    `json:"LocalPath"`
	MediaURL   string    `json:"MediaURL"`
	ContentMD5 string    `json:"ContentMD5"`
	MediaMeta  MediaMeta `json:"MediaMeta"`
	Persistent bool      `json:"Persistent"`
}
