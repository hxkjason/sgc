package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/hxkjason/sgc/constants"
	"time"
)

type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = JSONTime{Time: time.Time{}}
		return
	}

	now, err := time.Parse(constants.DateTimeLayout, string(data))
	*t = JSONTime{Time: now.Local()}
	return
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(constants.DateTimeLayout))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// JsonMarshalDisableEscapeHtml JsonMarshal禁用escapeHtml
func JsonMarshalDisableEscapeHtml(data interface{}) ([]byte, error) {

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(data)
	return bf.Bytes(), err
}
