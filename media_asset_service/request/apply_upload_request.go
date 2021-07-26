package request

// ApplyUploadRequest
type ApplyUploadRequest struct {
	CommonRequest
	Name       string    `json:"Name" validate:"required"`
	MediaMeta  MediaMeta `json:"MediaMeta" validate:"required,dive"`
	Size       string    `json:"Size" validate:"required"`
	ContentMD5 string    `json:"ContentMD5"`
	Inner      bool      `json:"Inner"`
}
