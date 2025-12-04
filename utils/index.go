package utils

import "fmt"

func String_or_num_to_string(value any) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return fmt.Sprintf("%d", value)
	default:
		panic("only string and int are supported")
	}
}
