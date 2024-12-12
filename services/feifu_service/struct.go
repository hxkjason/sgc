package feishu_service

const (
	SuccessEmoji = " ✔"
	FailedEmoji  = " ❌"

	WebhookUrlDefault    = "https://open.feishu.cn/open-apis/bot/v2/hook/929608a0-f812-4d24-9579-8e7c6da4535e"
	WebhookSecretDefault = "GRCIIXUKMFBnxQ1KCjIRYf"
)

var (
	varOncePeriodDuration int64 = 300000000000 // 纳秒
)

type (
	// SendFeiShuMsgRes 发送飞书消息响应
	SendFeiShuMsgRes struct {
		Code          int         `json:"code"`
		Data          interface{} `json:"data"`
		Msg           string      `json:"msg"`
		StatusCode    int         `json:"StatusCode"`
		StatusMessage string      `json:"StatusMessage"`
	}

	CardMsgItem struct {
		Name    string
		Content string
	}

	CardMsgAttr struct {
		Server      string // 服务器
		RequestTime string // 当前时间
		Cost        string // 耗时
		Trace       string // 跟踪
		Desc        string // 描述
		MessageId   string // 消息标识
		Items       []CardMsgItem
	}

	// CardMsg 卡片消息
	CardMsg struct {
		Config   Config    `json:"config"`
		Elements []Element `json:"elements"`
		Header   Header    `json:"header"`
	}

	Config struct {
		EnableForward bool `json:"enable_forward"` // 是否支持转发卡片。默认值 true
		UpdateMulti   bool `json:"update_multi"`   // 是否为共享卡片。为 true 时即更新卡片的内容对所有收到这张卡片的人员可见。默认值 false。
	}
	Header struct {
		Template string     `json:"template,omitempty"` // 标题主题颜色
		Title    ContentTag `json:"title"`              // 标题
		SubTitle ContentTag `json:"subtitle,omitempty"` // 副标题
		UdIcon   UdIcon     `json:"ud_icon"`            // 图标
	}
	ContentTag struct {
		Content interface{} `json:"content"`
		Tag     string      `json:"tag"`
	}
	UdIcon struct {
		Tag   string `json:"tag"`
		ToKen string `json:"token"`
	}

	/*
		Element 元素
		分栏组建文档: https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-components/containers/column-set
	*/
	Element struct {
		Tag               string      `json:"tag"` // 组件的标签。分栏组件的固定值为 column_set
		Content           interface{} `json:"content,omitempty"`
		Text              ContentTag  `json:"text,omitempty"`
		TextAlign         string      `json:"text_align,omitempty"`
		TextSize          string      `json:"text_size,omitempty"`
		FlexMode          string      `json:"flex_mode,omitempty"`          // 移动端和 PC 端的窄屏幕下，各列的自适应方式。默认值 none。
		BackgroundStyle   string      `json:"background_style,omitempty"`   // 分栏的背景色样式。默认值 default
		HorizontalSpacing string      `json:"horizontal_spacing,omitempty"` // 分栏中列容器之间的间距。默认值 default（为 8px）[0,28px]：自定义间距
		HorizontalAlign   string      `json:"horizontal_align,omitempty"`   // 列容器水平对齐的方式。默认值 left。
		Columns           []Column    `json:"columns,omitempty"`
		Margin            string      `json:"margin,omitempty"` // 列容器的外边距。
	}
	Column struct {
		Tag             string    `json:"tag"`
		Width           string    `json:"width"`
		VerticalAlign   string    `json:"vertical_align"`
		VerticalSpacing string    `json:"vertical_spacing"`
		BackgroundStyle string    `json:"background_style"`
		Elements        []Element `json:"elements"`
		Weight          int       `json:"weight"`
	}
)
