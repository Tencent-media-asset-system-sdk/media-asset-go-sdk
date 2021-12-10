package request

// DescribeNewURLRequest
type DescribeNewURLRequest struct {
	CommonRequest
	URLSet []string `json:"URLSet" validate:"required,dive"`
	Inner  bool     `json:"Inner"`
}
