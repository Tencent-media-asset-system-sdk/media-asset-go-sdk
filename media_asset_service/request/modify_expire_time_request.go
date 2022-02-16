package request

// ModifyExpireTimeRequest
type ModifyExpireTimeRequest struct {
	CommonRequest
	MediaID string `json:"MediaID" validate:"required"`
	Days    int32  `json:"Days" validate:"required,gt=0"`
}
