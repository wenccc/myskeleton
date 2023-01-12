package util

import (
	"fmt"
	"testing"
)

func TestArrayColumn(t *testing.T) {
	users := []map[string]string{
		{
			"user": "张三",
			"age":  "10",
		},
		{
			"user": "李四",
			"age":  "10",
		},
		{},
	}

	fmt.Println(arrayColumn(users, func(m map[string]string) (target any, find bool) {
		user, ok := m["user"]
		return user, ok
	}))
}
