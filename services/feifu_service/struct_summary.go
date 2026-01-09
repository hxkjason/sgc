package feishu_service

const (
	MessageTypeInteractive = "interactive" // 卡片消息
	MessageTypeSummary     = "summary"     // 汇总消息

	HeaderColorTypeWarning = 0
	HeaderColorTypeSuccess = 1
	HeaderColorTypeFail    = 2
	HeaderColorTypeNotice  = 3
)

type (
	SendFSMsgReq struct {
		AppId          string         `json:"AppId"`                                                    // 应用ID
		ReceiveId      string         `json:"ReceiveId" binding:"required"`                             // 接收对象ID
		ReceiveIdType  string         `json:"ReceiveIdType" binding:"required,oneof=chat_id open_id"`   // 接收对象类型 chat_id:群组ID，open_id:用户ID
		MessageType    string         `json:"MessageType" binding:"required,oneof=interactive summary"` // 消息类型 interactive:卡片消息
		CardMsgAttr    CardMsgAttr    `json:"CardMsgAttr,omitempty"`                                    // 消息属性
		SummaryMsgAttr SummaryMsgAttr `json:"SummaryMsgAttr,omitempty"`                                 // 汇总消息属性
	}

	// SummaryMsgAttr 汇总统计消息
	SummaryMsgAttr struct {
		MsgType  int    // 消息类型 0:warning 1:success 2:fail 3:notice
		Title    string // 标题
		SubTitle string // 副标题
		Elements []SummaryElement
	}

	SummaryElement struct {
		Tag       string    `binding:"required,oneof=hr title head detail"`
		Content   string    `json:"Content,omitempty"`
		TextAlign string    `json:"TextAlign,omitempty"`
		RowList   []RowItem `json:"RowList,omitempty"`
	}

	RowItem struct {
		TextList []interface{}
	}
)
