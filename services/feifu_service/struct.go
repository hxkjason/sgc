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

	// CardMsg 卡片消息
	CardMsg struct {
		Config   Config    `json:"config"`
		Elements []Element `json:"elements"`
		Header   Header    `json:"header"`
	}

	Config struct {
		UpdateMulti bool `json:"update_multi"`
	}
	Header struct {
		Template string     `json:"template,omitempty"` // 标题主题颜色
		Title    ContentTag `json:"title"`
	}

	ContentTag struct {
		Content interface{} `json:"content"`
		Tag     string      `json:"tag"`
	}

	Element struct {
		Tag               string      `json:"tag"`
		Content           interface{} `json:"content,omitempty"`
		Text              ContentTag  `json:"text,omitempty"`
		TextAlign         string      `json:"text_align,omitempty"`
		TextSize          string      `json:"text_size,omitempty"`
		FlexMode          string      `json:"flex_mode,omitempty"`
		BackgroundStyle   string      `json:"background_style,omitempty"`
		HorizontalSpacing string      `json:"horizontal_spacing,omitempty"`
		HorizontalAlign   string      `json:"horizontal_align,omitempty"`
		Columns           []Column    `json:"columns,omitempty"`
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
