package media

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type MediaImage struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}
//convert byte , then string and send to database ,when it's Send to database(It works automatically)
func (m MediaImage) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil   // ← string return করো, []byte না
}
//convert byte , then string and send to server ,when it's fetched from database (It works automatically)

func (m *MediaImage) Scan(value interface{}) error {
	if value == nil {
		*m = MediaImage{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		// কখনো string-ও আসতে পারে
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return errors.New("failed to scan MediaImage")
		}
	}
	return json.Unmarshal(bytes, m)
}

type MediaImageList []MediaImage
//convert byte , then string and send to database ,when it's Send to database(It works automatically)

func (m MediaImageList) Value() (driver.Value, error) {
	if m == nil {
		return "[]", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil   // ← string
}
//convert byte , then string and send to server ,when it's fetched from database (It works automatically)

func (m *MediaImageList) Scan(value interface{}) error {
	if value == nil {
		*m = MediaImageList{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return errors.New("failed to scan MediaImageList")
		}
	}
	return json.Unmarshal(bytes, m)
}