package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTodayQuestion(t *testing.T) {
	questions, err := GetTodayQuestion(time.Now().AddDate(0, 0, 1), time.Now())
	if err != nil {
		return
	}
	for _, question := range questions {
		fmt.Println(question.Type)
	}
}
