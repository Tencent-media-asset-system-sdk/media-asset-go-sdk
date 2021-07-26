package request

// DescribeMediaDetailsRequest
type DescribeMediaDetailsRequest struct {
	CommonRequest
	MediaIDSet []uint64 `json:"MediaIDSet" validate:"required"`
}
