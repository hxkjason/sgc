package feishu_service

func SendCardMessage(rb CardMsgAttr) error {
	return SendCardMsg(GenCardMsgJson(rb), rb.WebhookUrl, rb.WebhookSecret)
}

func GenCardMsgJson(rb CardMsgAttr) (cardMsg CardMsg) {
	cardMsg = CardMsg{
		Config: Config{
			EnableForward: true,
			UpdateMulti:   true,
		},
		Elements: nil,
		Header: Header{
			Template: "",
			Title:    ContentTag{Content: rb.Title, Tag: "plain_text"},
			SubTitle: ContentTag{Content: rb.SubTitle, Tag: "plain_text"},
			UdIcon:   UdIcon{Tag: "standard_icon"},
		},
	}

	switch rb.MsgType {
	case MsgTypeWarning:
		cardMsg.Header.Template = "orange"
		cardMsg.Header.UdIcon.ToKen = "meeting-ai_filled"
	case MsgTypeSuccess:
		cardMsg.Header.Template = "turquoise"
		cardMsg.Header.UdIcon.ToKen = "sheet-iconsets-check_filled"
	case MsgTypeDanger:
		cardMsg.Header.Template = "red"
		cardMsg.Header.UdIcon.ToKen = "sheet-iconsets-cross_filled"
	default:
		cardMsg.Header.Template = "carmine"
		cardMsg.Header.UdIcon.ToKen = "meeting-ai_filled"
	}

	if rb.AppName != "" {
		cardMsg.AppendColumnSetElement("应 用", rb.AppName)
	}
	if rb.Server != "" {
		cardMsg.AppendColumnSetElement("服 务", rb.Server)
	}
	if rb.RequestTime != "" {
		cardMsg.AppendColumnSetElement("时 间", rb.RequestTime)
	}
	if rb.Cost != "" {
		cardMsg.AppendColumnSetElement("耗 时", rb.Cost)
	}
	if rb.Trace != "" {
		cardMsg.AppendColumnSetElement("跟 踪", rb.Trace)
	}
	if rb.Desc != "" {
		cardMsg.AppendColumnSetElement("描 述", rb.Desc)
	}
	for _, item := range rb.Items {
		if item.Name != "" {
			cardMsg.AppendColumnSetElement(item.Name, item.Content)
		}
	}
	if rb.MessageId != "" {
		cardMsg.AppendColumnSetElement("标 识", rb.MessageId)
	}
	return cardMsg
}

// AppendColumnSetElement 添加分栏元素
func (cm *CardMsg) AppendColumnSetElement(title, content string) *CardMsg {
	cm.Elements = append(cm.Elements, Element{
		Tag:               "column_set",
		Content:           nil,
		Text:              ContentTag{},
		TextAlign:         "",
		TextSize:          "",
		FlexMode:          "none",
		BackgroundStyle:   "",
		HorizontalSpacing: "4px",
		HorizontalAlign:   "left",
		Columns: []Column{
			{
				Tag:             "column",
				Width:           "weighted",
				VerticalAlign:   "top",
				VerticalSpacing: "6px",
				BackgroundStyle: "",
				Elements: []Element{
					{
						Tag:       "markdown",
						Content:   "**" + title + ":**",
						TextAlign: "center",
						TextSize:  "normal",
					},
				},
				Weight: 1,
			},
			{
				Tag:             "column",
				Width:           "weighted",
				VerticalAlign:   "top",
				VerticalSpacing: "6px",
				BackgroundStyle: "",
				Elements: []Element{
					{
						Tag:       "markdown",
						Content:   content,
						TextAlign: "left",
						TextSize:  "normal",
					},
				},
				Weight: 5,
			},
		},
		Margin: "0px",
	})
	return cm
}
