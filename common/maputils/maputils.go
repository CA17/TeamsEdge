package maputils

import (
	"errors"
	"fmt"
	"time"

	"github.com/ca17/teamsedge/common"
)

func GetStringValue(d map[string]interface{}, key string, defval string) string {
	val, ok := d[key]
	if !ok || val == nil || val == "" {
		return defval
	}
	val2, err := common.ParseString(val)
	if err != nil {
		return defval
	}
	return val2
}

func GetStringValueWithErr(d map[string]interface{}, key string) (string, error) {
	val, ok := d[key]
	if !ok || val == nil || val == "" {
		return "", fmt.Errorf("%s is empty", key)
	}
	val2, err := common.ParseString(val)
	if err != nil {
		return "", err
	}
	return val2, nil
}

var errval = errors.New("err int64 value")

func GetInt64ValueWithErr(d map[string]interface{}, key string) (int64, error) {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, errval
}

func GetIntValue(d map[string]interface{}, key string, defval int) int {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return defval
		}
		return int(v)
	}
	return defval
}

func GetInt64Value(d map[string]interface{}, key string, defval int64) int64 {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return defval
		}
		return v
	}
	return defval
}

func GetFloat64Value(d map[string]interface{}, key string, defval float64) float64 {
	val, ok := d[key]
	if ok {
		v, err := common.ParseFloat64(val)
		if err != nil {
			return defval
		}
		return v
	}
	return defval
}

func GetDateObject(d map[string]interface{}, key string, defval time.Time) time.Time {
	val, ok := d[key]
	if ok {
		var result = defval
		val, err := common.ParseString(val)
		if err != nil {
			return defval
		}
		if len(val) == 19 {
			result, err = time.Parse("2006-01-02 15:04:05", val)
		} else {
			result, err = time.Parse("2006-01-02 15:04:05 Z0700 MST", val)
		}
		if err != nil {
			return defval
		}
		return result
	}
	return defval
}

func GetSIntValue(d map[string]string, key string, defval int) int {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return defval
		}
		return int(v)
	}
	return defval
}

func GetSInt64Value(d map[string]string, key string, defval int64) int64 {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return defval
		}
		return v
	}
	return defval
}
