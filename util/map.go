package util

import (
	"encoding/json"
)

func StructToMap(in any) map[string]any {
	var out map[string]interface{}
	bytes, _ := json.Marshal(in)
	_ = json.Unmarshal(bytes, &out)
	return out
}

func MapToStruct[T any](in map[string]any) T {
	var out T
	bytes, _ := json.Marshal(in)
	_ = json.Unmarshal(bytes, &out)
	return out
}
