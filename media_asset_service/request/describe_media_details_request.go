package request

// DescribeMediaDetailsRequest
type DescribeMediaDetailsRequest struct {
	CommonRequest
	MediaIDSet []string `json:"MediaIDSet" validate:"required"`
}
