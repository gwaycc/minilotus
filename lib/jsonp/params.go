package jsonp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gwaylib/errors"
)

var (
	UNIX_TIME_NO_SET = time.Time{}.Unix()
)

type Params map[string]interface{}

func ParseParams(data []byte) (Params, error) {
	params := Params{}
	if err := json.Unmarshal(data, &params); err != nil {
		return params, errors.As(err)
	}
	return params, nil
}

func ParseParamsByIO(r io.Reader) (Params, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.As(err)
	}
	return ParseParams(body)
}

func (p Params) Set(key string, val interface{}) {
	fmt.Println(p, key, val)
	//p[key] = val
}

func (p Params) JsonData() []byte {
	data, _ := json.Marshal(p)
	return data
}

func (p Params) String(key string) string {
	s, ok := p[key]
	if !ok {
		return ""
	}
	r, ok := s.(string)
	if ok {
		return r
	}
	return fmt.Sprint(s)
}
func (p Params) Int64(key string, def int64) int64 {
	s, ok := p[key]
	if !ok {
		return def
	}
	r, ok := s.(int64)
	if ok {
		return r
	}
	i, err := strconv.ParseInt(fmt.Sprint(s), 10, 64)
	if err != nil {
		return def
	}
	return i
}
func (p Params) Float64(key string, def float64) float64 {
	s, ok := p[key]
	if !ok {
		return def
	}
	r, ok := s.(float64)
	if ok {
		return r
	}
	i, err := strconv.ParseFloat(fmt.Sprint(s), 64)
	if err != nil {
		return def
	}
	return i
}

func (p Params) Time(key string, layoutOpt ...string) time.Time {
	layout := "2006-01-02 15:04:05" // default format, using the local locate.
	if len(layoutOpt) > 0 {
		layout = layoutOpt[0]
	}
	s, ok := p[key]
	if !ok {
		return time.Time{}
	}
	r, ok := s.(time.Time)
	if ok {
		return r
	}
	t, _ := time.Parse(layout, s.(string))
	return t
}

func (p Params) Email(key string) string {
	email := p.String(key)
	if strings.Index(email, "@") < 1 {
		return ""
	}
	for _, r := range email {
		if r > 255 {
			return ""
		}
	}
	return email
}

func (p Params) StringArray(key string) []string {
	s, ok := p[key]
	if !ok {
		return []string{}
	}
	arr, ok := s.([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, len(arr))
	for i, a := range arr {
		result[i] = fmt.Sprint(a)
	}
	return result
}
func (p Params) ParamsArray(key string) []Params {
	s, ok := p[key]
	if !ok {
		return []Params{}
	}
	arr, ok := s.([]interface{})
	if !ok {
		return []Params{}
	}
	result := []Params{}
	for _, a := range arr {
		p, ok := a.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, p)
	}
	return result
}
