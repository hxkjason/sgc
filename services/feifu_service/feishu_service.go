package feishu_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"runtime"
	"sgc/constants"
	"sgc/redis"
	"sgc/utils"
	"strconv"
	"strings"
	"time"
)

const (
	SuccessEmoji    = " ✔"
	FailedEmoji     = " ❌"
)

var (
	varOncePeriodDuration int64 = 300000000000 // 纳秒
)

type (
	// FeiShuMsg 飞书消息体
	FeiShuMsg struct {
		Title   string                     `json:"title"`
		Content [][]map[string]interface{} `json:"content"`
	}

	// LogContent 消息应用内容
	LogContent struct {
		AppName    string `json:"AppName"` // 服务器
		Type       int    `json:"Type"`
		Level      int    `json:"Level"`
		Message    string `json:"Message"` // 描述
		StackTrace string `json:"StackTrace"`
		Source     string `json:"Source"`
		Server     string `json:"Server"`
		DateTime   string `json:"DateTime"`
		MessageId  string `json:"MessageId"`
	}
)

// SendDevopsMsg 发送 Devops Msg
func SendDevopsMsg(noticeMsg, webhookUrl, secret string) {

	hostname, _ := os.Hostname()
	logContent := LogContent{
		AppName:    os.Getenv(constants.RunEnv),
		Message:    noticeMsg,
		DateTime:   time.Now().Format(constants.DateTimeLayout),
		Server:     hostname,
		StackTrace: getCallInfo(),
		MessageId:  uuid.NewV4().String(),
	}

	// 几分钟内只发一次
	cacheKey := utils.SplicingStr(logContent.AppName, "_", logContent.Message)
	md5StrKey := utils.Md5(cacheKey)
	if RecentlyHasSend(md5StrKey, varOncePeriodDuration) {
		return
	}

	err := SendMsg(logContent, webhookUrl, secret)
	if err == nil {
		redis.SetKeyValue(md5StrKey, 1, time.Duration(varOncePeriodDuration))
	} else {
		// 如果发送失败，删除一段时间内已发送过或正在发送的缓存
		delCount, _ := redis.DelKey(md5StrKey)
		if delCount < 1 {
			fmt.Println("删除失败,md5StrKey:", md5StrKey)
		}
	}
}

func getCallInfo() string {
	_, file, line, _ := runtime.Caller(2)
	return utils.SplicingStr(file, ":", strconv.Itoa(line))
}

func RecentlyHasSend(md5CacheKey string, duration int64) bool {

	periodNotSend, _ := redis.SetKeyNX(md5CacheKey, 1, time.Duration(duration))
	if !periodNotSend {
		fmt.Println("该消息在", duration/1000000000, " 秒内已发送过或正在发送")
		return true
	}
	return false
}

// SendMsg 发送消息
func SendMsg(logContent LogContent, webhookUrl, secret string) error {
	if webhookUrl == "" {
		webhookUrl, secret = "https://open.feishu.cn/open-apis/bot/v2/hook/929608a0-f812-4d24-9579-8e7c6da4535e", "GRCIIXUKMFBnxQ1KCjIRYf"
	}
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	sign, err := GenSign(secret, timestampStr)
	if err != nil {
		log.Println("feiShu GenSign has err:" + err.Error())
		return err
	}
	// msgContent
	msgContent := getFeiShuMsgContent(logContent)
	// data
	sendData := `{"timestamp": "` + timestampStr + `","sign": "` + sign + `","msg_type": "post","content": {"post":{"zh_cn":` + msgContent + `}}}`
	// request
	result, err := http.Post(webhookUrl, "application/json", strings.NewReader(sendData))
	defer result.Body.Close()

	if err != nil {
		log.Println(utils.SplicingStr("发送告警消息到飞书群组出错,MessageId:", logContent.MessageId, " err:", err.Error()))
		return err
	}
	if result.StatusCode != 200 {
		err = errors.New(result.Status)
		log.Println(utils.SplicingStr("post msg to feiShu failed, messageId:", logContent.MessageId, ",status:", result.Status, "msgContent:", msgContent))
	}
	return err
}

// getFeiShuMsgContent 整理飞书消息的发送内容
func getFeiShuMsgContent(logContent LogContent) (feiShuMsgContent string) {
	msg := FeiShuMsg{Title: logContent.AppName}
	content := [][]map[string]interface{}{
		{},
		{
			map[string]interface{}{"tag": "text", "un_escape": true, "text": "服务器&nbsp;:&nbsp;&nbsp;"},
			map[string]interface{}{"tag": "text", "text": logContent.Server},
		},
		{
			map[string]interface{}{"tag": "text", "un_escape": true, "text": "时&nbsp;&nbsp;&nbsp;&nbsp;间&nbsp;:&nbsp;&nbsp;"},
			map[string]interface{}{"tag": "text", "text": logContent.DateTime},
		},
		{
			map[string]interface{}{"tag": "text", "un_escape": true, "text": "资&nbsp;&nbsp;&nbsp;&nbsp;源&nbsp;:&nbsp;&nbsp;"},
			map[string]interface{}{"tag": "text", "text": logContent.Source},
		},
		{
			map[string]interface{}{"tag": "text", "un_escape": true, "text": "跟&nbsp;&nbsp;&nbsp;&nbsp;踪&nbsp;:&nbsp;&nbsp;"},
			map[string]interface{}{"tag": "text", "text": logContent.StackTrace},
		},
		{
			map[string]interface{}{"tag": "text", "un_escape": true, "text": "描&nbsp;&nbsp;&nbsp;&nbsp;述&nbsp;:&nbsp;&nbsp;"},
			map[string]interface{}{"tag": "text", "text": logContent.Message},
		},
		{
			map[string]interface{}{"tag": "text", "un_escape": true, "text": "标&nbsp;&nbsp;&nbsp;&nbsp;识&nbsp;:&nbsp;&nbsp;"},
			map[string]interface{}{"tag": "text", "text": logContent.MessageId},
		},
	}

	msg.Content = content
	jsonBytes, _ := json.Marshal(&msg)
	return string(jsonBytes)
}

// GenSign 生成签名
func GenSign(secret, timestampStr string) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	forSignStr := utils.SplicingStr(timestampStr, "\n", secret)

	var data []byte
	h := hmac.New(sha256.New, []byte(forSignStr))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
