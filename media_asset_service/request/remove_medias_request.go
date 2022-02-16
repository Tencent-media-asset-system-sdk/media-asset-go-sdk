package request

// RemoveMediasRequest
type RemoveMediasRequest struct {
	CommonRequest
	MediaIDSet []string `json:"MediaIDSet" validate:"required,lte=100"`
}
