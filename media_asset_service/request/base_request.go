package request

// Request
type Request interface {
	GetRequestId() string
}

// CommonRequest
type CommonRequest struct {
	RequestID     string `json:"RequestId,omitempty"`
	AppID         uint64 `json:"AppID,omitempty"`
	Action        string `json:"Action,omitempty"`
	Uin           string `json:"Uin,,omitempty"`
	SubAccountUin string `json:"SubAccountUin,omitempty"`
	TIBusinessID  uint32 `json:"TIBusinessID,omitempty"`
	TIProjectID   uint32 `json:"TIProjectID,omitempty"`
}

// GetRequestId
func (c CommonRequest) GetRequestId() string {
	return c.RequestID
}

// MediaMeta
type MediaMeta struct {
	MediaType      string `json:"MediaType" validate:"required"`
	MediaTag       string `json:"MediaTag" validate:"required"`
	MediaSecondTag string `json:"MediaSecondTag"`
	MediaLang      string `json:"MediaLang"` // 普通话, 粤语
}

// FilterBy
type FilterBy struct {
	MediaNameOrID string   `json:"MediaNameOrID"`
	MediaTypeSet  []string `json:"MediaTypeSet"`
	MediaTagSet   []Label  `json:"MediaTagSet"`
	StatusSet     []string `json:"StatusSet"`
}

// Category
type Category struct {
	Type   string   `json:"Type"`
	TagSet []string `json:"TagSet"`
}

// Label
type Label struct {
	Type         string   `json:"Type"`
	Tag          string   `json:"Tag"`
	SecondTagSet []string `json:"SecondTagSet"`
}
