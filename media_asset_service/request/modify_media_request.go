package request

// ModifyMediaRequest
type ModifyMediaRequest struct {
	CommonRequest
	MediaID        uint64 `json:"MediaID" validate:"required"`
	MediaTag       string `json:"MediaTag"`
	MediaSecondTag string `json:"MediaSecondTag"`
}
