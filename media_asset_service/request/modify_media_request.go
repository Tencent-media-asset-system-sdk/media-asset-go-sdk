package request

// ModifyMediaRequest
type ModifyMediaRequest struct {
	CommonRequest
	MediaID        string `json:"MediaID" validate:"required"`
	MediaTag       int    `json:"MediaTag"`
	MediaSecondTag int    `json:"MediaSecondTag"`
}
