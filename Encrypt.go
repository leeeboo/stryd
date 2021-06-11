package stryd

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func ParseMapToStruct(raw interface{}, dst interface{}) error {

	jsonData, err := json.Marshal(raw)

	if err != nil {
		return errors.New("[ParseMapToStruct] " + err.Error())
	}

	d := json.NewDecoder(strings.NewReader(string(jsonData)))
	d.UseNumber()
	err = d.Decode(dst)
	if err != nil {
		return errors.New("[ParseMapToStruct] " + err.Error())
	}
	return nil
}

func Md5(raw string) string {

	m := md5.New() // #nosec
	_, err := m.Write([]byte(raw))

	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}

func Sha1(raw string) string {
	h := sha1.New() // #nosec
	_, err := h.Write([]byte(raw))

	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
