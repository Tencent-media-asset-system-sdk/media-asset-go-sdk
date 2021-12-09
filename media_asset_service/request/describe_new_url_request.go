package request

// DescribeNewURLRequest
type DescribeNewURLRequest struct {
	CommonRequest
	URLSet []string `json:"url_set" validate:"required,dive"`
	Inner  bool     `json:"Inner"`
}
