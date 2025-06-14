package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Int64Array []int64

func (a *Int64Array) UnmarshalJSON(b []byte) error {
	var temp []string
	// 这里的 json.Unmarshal 会调用 json-iterator 的解析器
	if err := json.Unmarshal(b, &temp); err != nil {
		// 同样，为了健壮性，我们尝试直接解析为整数数组
		var tempInts []int64
		if err2 := json.Unmarshal(b, &tempInts); err2 != nil {
			return err // 返回原始错误
		}
		*a = tempInts
		return nil
	}

	result := make(Int64Array, len(temp))
	for i, s := range temp {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("无法将字符串 '%s' 解析为 int64: %w", s, err)
		}
		result[i] = val
	}

	*a = result
	return nil
}

func (a Int64Array) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	if len(a) == 0 {
		return []byte("[]"), nil
	}

	var builder strings.Builder
	for i, v := range a {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(`"`)
		builder.WriteString(strconv.FormatInt(v, 10))
		builder.WriteString(`"`)
	}

	return []byte("[" + builder.String() + "]"), nil
}
