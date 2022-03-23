package main

import (
	"context"
	"fmt"
	. "github.com/b1gdea1/checkin_client/admin_service"
	. "github.com/b1gdea1/checkin_client/base_service"
	_ "github.com/b1gdea1/checkin_client/common_service"
	"google.golang.org/grpc"
	"testing"
)

func TestAdminUpdateQuestion(t *testing.T) {
	conn, err := grpc.Dial(":8002", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	checkinServiceClient := NewAdminServiceClient(conn)
	response, err := checkinServiceClient.UpdateQuestions(context.Background(), &QuestionUpdateRequest{
		AdminId: "332526",
		Token:   "81731109891445270481646892449",
		Questions: []*Question{
			&Question{
				Type:    QuestionType_NORMAL_QUESTION,
				Content: "how are you?",
				Answers: []string{},
				Depends: nil,
			},
			&Question{
				Type:    QuestionType_SINGLE_SELECTION,
				Content: "where are you from?",
				Answers: []string{"1", "2", "3"},
				Depends: nil,
			},
			&Question{
				Type:    QuestionType_MULTI_SELECTION,
				Content: "112212",
				Answers: []string{"1", "2", "3", "4"},
				Depends: []string{"2"},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.Token)
	fmt.Println(response.Result)
}
