package i18n;

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type LocalizedString map[string]string

func (ls LocalizedString) Value() (driver.Value, error) {
	if ls == nil {
		return nil, nil
	}
	return json.Marshal(ls)
}

func (ls *LocalizedString) Scan(value interface{}) error {
	if value == nil {
		*ls = make(LocalizedString)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan LocalizedString")
	}
	return json.Unmarshal(bytes, ls)
}

func (ls LocalizedString) Get(lang string) string {
	if val, ok := ls[lang]; ok && val != "" {
		return val
	}
	if val, ok := ls["en"]; ok && val != "" {
		return val
	}
	for _, v := range ls {
		if v != "" {
			return v
		}
	}
	return ""
}
