package request

// RemoveMediasRequest
type RemoveMediasRequest struct {
	CommonRequest
	MediaIDSet []uint64 `json:"MediaIDSet" validate:"required,lte=100"`
}
