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
	MediaType      int `json:"MediaType" validate:"required"`
	MediaTag       int `json:"MediaTag" validate:"required"`
	MediaSecondTag int `json:"MediaSecondTag"`
	MediaLang      int `json:"MediaLang"` // 普通话, 粤语
}

// FilterBy
type FilterBy struct {
	MediaNameOrID string  `json:"MediaNameOrID"`
	MediaTypeSet  []int   `json:"MediaTypeSet"`
	MediaTagSet   []Label `json:"MediaTagSet"`
	StatusSet     []int   `json:"StatusSet"`
}

// Category
type Category struct {
	Type   int   `json:"Type"`
	TagSet []int `json:"TagSet"`
}

// Label
type Label struct {
	Type         int   `json:"Type"`
	Tag          int   `json:"Tag"`
	SecondTagSet []int `json:"SecondTagSet"`
}
