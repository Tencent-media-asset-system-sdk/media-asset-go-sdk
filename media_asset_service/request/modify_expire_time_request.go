package request

// ModifyExpireTimeRequest
type ModifyExpireTimeRequest struct {
	CommonRequest
	MediaID uint64 `json:"MediaID" validate:"required"`
	Days    int32  `json:"Days" validate:"required,gt=0"`
}
