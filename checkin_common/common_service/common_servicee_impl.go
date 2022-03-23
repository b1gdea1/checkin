package common_service

import (
	"context"
	. "github.com/b1gdea1/checkin_common/base_service"
	"github.com/b1gdea1/checkin_common/mysqlDao"
	"github.com/b1gdea1/checkin_common/redis"
	"github.com/b1gdea1/checkin_common/util"
	"log"
	"strings"
	"time"
)

type Server struct {
}

func (u Server) GetQuestions(ctx context.Context, request *QuestionRequest) (detail *QuestionDetail, err error) {
	questions, err := mysqlDao.GetQuestions()
	if err != nil {
		return nil, err
	}
	detail = new(QuestionDetail)
	for _, question := range questions {
		detail.Questions = append(detail.Questions, &Question{
			Type:    QuestionType(question.Type),
			Content: question.Content,
			Answers: strings.Split(question.Answers, "\t"),
			Depends: strings.Split(question.Depends, "\t"),
		})
	}
	return detail, nil
}

func checkSchoolID(schoolID string) bool {
	// TODO 检查学号格式
	return true
}
func checkPassword(p string) bool {
	// TODO 检查密码格式
	return true
}
func (u Server) RegisterUser(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	log.Println(req)
	ID := req.UserID
	password := req.Password
	adminID := req.AdminID
	adminToken := req.AdminToken
	resp := &RegisterResponse{}
	// 验证管理员
	b1 := mysqlDao.AdminExist(adminID)
	if !b1 {
		resp.Result = RegisterResult_ADMIN_DOSE_NOT_EXIST
		return resp, nil
	}
	b2 := redis.VerifyToken(adminID, adminToken)
	if !b2 {
		resp.Result = RegisterResult_ADMIN_DOSE_NOT_LOGIN
		return resp, nil
	}
	id := mysqlDao.GetStudentID(ID, password)
	if id.Username != "" {
		resp.Result = RegisterResult_ALREADY_HAVE_SCHOOL_ID
		return resp, nil
	}
	if !checkSchoolID(ID) {
		resp.Result = RegisterResult_ILLEGAL_SCHOOL_ID
		return resp, nil
	}
	if !checkPassword(password) {
		resp.Result = RegisterResult_ILLEGAL_PASSWORD
	}
	role := req.Type
	res := false
	if role == 0 {
		res = mysqlDao.InsertStudentSimple(ID, password)
	} else {
		res = mysqlDao.InsertAdminSimple(ID, password)
	}
	
	if res {
		resp.Result = RegisterResult_REGISTER_SUCCESS
	} else {
		resp.Result = RegisterResult_REGISTER_FAIL
	}
	return resp, nil
}

func (u Server) LoginRPC(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// TODO:通过反射完成执行任何方法的前后输出日志
	log.Printf("login request:%v", req)
	role := req.Type
	ID := req.UserID
	password := req.Password
	resp := &LoginResponse{}
	(*resp).State = LoginResult_LOGIN_SUCCESS
	if role == 0 {
		s := mysqlDao.GetStudentID(ID, password)
		if s.ID == 0 {
			(*resp).State = LoginResult_LOGIN_WRONG_PASSWORD
		}
	} else {
		s := mysqlDao.GetAdministratorID(ID, password)
		if s.ID == 0 {
			(*resp).State = LoginResult_LOGIN_WRONG_PASSWORD
		}
	}
	if (*resp).State == LoginResult_LOGIN_WRONG_PASSWORD {
		return resp, nil
	}
	token := util.GetToken()
	ret := redis.UploadUserToken(ID, token, 300)
	if !ret {
		(*resp).State = LoginResult_FAIL_TO_UPLOAD_TOKEN
	}
	log.Println("token:", token)
	(*resp).Token = token
	return resp, nil
}

func (u Server) Checkin(ctx context.Context, request *CheckinRequest) (resp *CheckinResponse, err error) {
	token := request.Token
	userID := request.UserID
	resp = new(CheckinResponse)
	resp.Result = CheckinResult_CHECKIN_SUCCESS
	b := redis.VerifyToken(userID, token)
	if b {
		ret := mysqlDao.UpdateRecord(mysqlDao.Record{
			SID:    request.UserID,
			Answer: request.Answers,
			Date:   time.Now().Format("2006-01-02"),
		})
		if !ret {
			resp.Result = CheckinResult_CHECKIN_FAILED
		}
	}
	token = util.GetToken()
	ret := redis.UploadUserToken(userID, token, 300)
	if !ret {
		resp.Result = CheckinResult_CHECKIN_FAILED
		return resp, nil
	}
	resp.Token = token
	return resp, err
}
