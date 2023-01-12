package util

import (
	"fmt"
	"testing"
)

func TestMapSetIfNotExist(t *testing.T) {
	values := map[string]string{
		"name": "张三",
	}

	cases := map[string]bool{
		"user": false,
		"name": true,
		"age":  false,
	}

	for key, exist := range cases {
		if MapSetIfNotExist(values, key, "==notExist==") != exist {
			t.Failed()
		}
	}
	fmt.Println(values)
}
