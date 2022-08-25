package utils

import (
	"fmt"
	"strconv"
	"strings"
)

//转为int
func GetInt(v interface{}) int {
	switch result := v.(type) {
	case int:
		return result
	case int32:
		return int(result)
	case int64:
		return int(result)
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.Atoi(d)
			return value
		}
	}
	return 0
}

//转为string
func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []string:
		return strings.Join(result, "")
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprint(result)
		}
	}
	return ""
}
