package mediaassetsdk

// 媒体类型
const (
	MediaTypeVideo = "视频"
	MediaTypeLive  = "直播流"
	MediaTypeImage = "图片"
	MediaTypeAudio = "音频"
	MediaTypeText  = "文稿"
)

// 媒体标签
const (
	MediaLabelNews          = "新闻"
	MediaLabelEntertainment = "综艺"
	MediaLabelInternetInfo  = "互联网资讯"
	MediaLabelMovie         = "电影"
	MediaLabelSeries        = "电视剧"
	MediaLabelSpecial       = "专题"
	MediaLabelSport         = "体育"
)

// 媒体二级标签
const (
	MediaSecondLabelEvening = "晚会"
	MediaSecondLabelOther   = "其他"
)

// 媒体语言
const (
	MediaLangMandarin  = "普通话"
	MediaLangCantonese = "粤语"
)

const (
	MediaStateUploading     = "上传中"
	MediaStateWaitingVerify = "等待验证"
	MediaStateCompleted     = "上传完成"
	MediaStateFailed        = "上传失败"
	MediaStateDownloading   = "下载素材中"
	MediaStateVerifying     = "验证素材中"
	MediaStateDeleted       = "素材已删除"
	MediaStateCleaned       = "素材已清理"
)
