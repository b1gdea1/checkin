package main

import (
	"context"
	"fmt"
	"github.com/b1gdea1/checkin_client/common_service"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestStudentLogin(t *testing.T) {
	conn, err := grpc.Dial(":8001", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	serviceClient := common_service.NewStudentServiceClient(conn)
	response, err := serviceClient.LoginRPC(context.Background(),
		&common_service.LoginRequest{
			Type:     0,
			UserID:   "201930507008",
			Password: "123456",
		})
	if err != nil {
		log.Fatalf("request error:%v", err)
		return
	}
	fmt.Println(response)
}

func TestAdminLogin(t *testing.T) {
	conn, err := grpc.Dial(":8001", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	serviceClient := common_service.NewStudentServiceClient(conn)
	response, err := serviceClient.LoginRPC(context.Background(),
		&common_service.LoginRequest{
			Type:     1,
			UserID:   "332526",
			Password: "123456",
		})
	if err != nil {
		log.Fatalf("request error:%v", err)
		return
	}
	fmt.Println(response)
}

func TestRegister(t *testing.T) {
	conn, err := grpc.Dial(":8001", grpc.WithInsecure())
	if err != nil {
		log.Printf("request error:%v", err)
		return
	}
	defer conn.Close()
	serviceClient := common_service.NewStudentServiceClient(conn)
	response, err := serviceClient.RegisterUser(context.Background(), &common_service.RegisterRequest{
		UserID:     "201930507003",
		Password:   "661243",
		AdminID:    "332526",
		AdminToken: "35195569942945600801646838423",
	})
	if err != nil {
		log.Printf("request error:%v", err)
		return
	}
	if err != nil {
		log.Printf("request error:%v", err)
		return
	}
	fmt.Println(response.Result)
}

func TestCheckin(t *testing.T) {
	conn, err := grpc.Dial(":8001", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	serviceClient := common_service.NewStudentServiceClient(conn)
	response, err := serviceClient.Checkin(context.Background(), &common_service.CheckinRequest{
		UserID:  "201930507008",
		Token:   "48546958373163795951646839115",
		Answers: "nihao",
	})
	if err != nil {
		return
	}
	fmt.Println(response)
}
