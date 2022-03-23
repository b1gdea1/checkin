package admin_service

import (
	"context"
	"github.com/b1gdea1/checkin_admin/mysqlDao"
	redisdriver "github.com/b1gdea1/checkin_admin/redis"
	"github.com/b1gdea1/checkin_admin/util"
	"strings"
)

type CheckinServiceAsAdministrator struct {
}

func (receiver CheckinServiceAsAdministrator) UpdateQuestions(_ context.Context,
	detail *QuestionUpdateRequest) (result *QuestionUpdateResponse, err error) {
	result = &QuestionUpdateResponse{}
	token := detail.Token
	adminID := detail.AdminId
	adminExist := mysqlDao.AdminExist(adminID)
	(*result).Result = QuestionUpdateResult_QUESTION_UPDATE_SUCCESS
	if !adminExist {
		(*result).Result = QuestionUpdateResult_QUESTION_UPDATE_NO_ADMIN
		return result, nil
	}
	verifyToken := redisdriver.VerifyToken(adminID, token)
	if !verifyToken {
		(*result).Result = QuestionUpdateResult_QUESTION_UPDATE_NO_LOGIN
		return result, nil
	}
	questions := detail.GetQuestions()
	var questionsSQL []*mysqlDao.Question
	for i, question := range questions {
		questionsSQL = append(questionsSQL, &mysqlDao.Question{
			QuestionID: int32(i),
			Type:       int32(question.Type),
			Content:    question.Content,
			Answers:    strings.Join(question.Answers, "\t"),
			Depends:    strings.Join(question.Depends, "\t"),
		})
	}
	err = redisdriver.PrepareTomorrowQuestion(questions) // 暂存到redis中
	if err != nil {
		result = nil
		return result, err
	}
	newToken := util.GetToken()
	ret := redisdriver.UploadUserToken(adminID, token, 300)
	if !ret {
		(*result).Result = QuestionUpdateResult_QUESTION_UPDATE_TOKEN_FAIL
		return result, nil
	}
	(*result).Token = newToken
	return
}

func (receiver CheckinServiceAsAdministrator) PublishNotice(ctx context.Context, detail *NoticeRequest) (*NoticeResp,
	error) {
	result := &NoticeResp{
		Result: NoticeResult_NOTICE_SEND_SUCCESS,
	}
	adminID := detail.AdminID
	if adminExist := mysqlDao.AdminExist(adminID); !adminExist {
		(*result).Result = NoticeResult_NOTICE_ADMIN_DOSE_NOT_EXIST
		return result, nil
	}
	token := detail.AdminToken
	if tokenLegal := redisdriver.VerifyToken(adminID, token); !tokenLegal {
		(*result).Result = NoticeResult_NOTICE_ADMIN_DOST_NOT_LOGIN
		return result, nil
	}
	userID := detail.UserID
	msg := detail.Msg
	ttl := detail.TTL
	b := redisdriver.PublishNotice(userID, msg, int(ttl))
	if !b {
		(*result).Result = NoticeResult_NOTICE_SEND_FAIL
	}
	return result, nil
	
}
