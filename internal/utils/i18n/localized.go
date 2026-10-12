package i18n;

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type LocalizedString map[string]string

func (ls LocalizedString) Value() (driver.Value, error) {  //jakhon data insert kora hoy database e takhon aita hoy
	if ls == nil {
		return "{}", nil
	}
	b, err := json.Marshal(ls) // byte convertion jakhon database e jai
	if err != nil {
		return nil, err
	}
	return string(b), nil   // ← []byte না, string   oi byte k string e convertion
}

func (ls *LocalizedString) Scan(value interface{}) error {   // jakhon output pawa lage takkon ai scan ta hoy.  {name:"bn"} erakom type
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
