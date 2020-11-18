package middleware

import (
	"time"

	. "cm-v5/serv/module"
)

type InputLogHistoryPaymentStruct struct {
	User_id     string
	Api         string
	Method      string
	Data_input  interface{} `bson:"data_input " json:"data_input"`
	Data_output interface{} `bson:"data_output " json:"data_output"`
	Created     string
}

func SaveLogHistoryPayment(UserId, Api, Method string, Input interface{}, Output interface{}) error {
	Session, Db, Err := GetCollection()
	if Err != nil {
		return Err
	}

	defer Session.Close()

	var InputLogHistoryPayment InputLogHistoryPaymentStruct
	InputLogHistoryPayment.User_id = UserId
	InputLogHistoryPayment.Api = Api
	InputLogHistoryPayment.Method = Method
	InputLogHistoryPayment.Data_input = Input
	InputLogHistoryPayment.Data_output = Output
	InputLogHistoryPayment.Created = time.Now().Format("2006-01-02 15:04:05")
	Err = Db.C("log_api_payment").Insert(InputLogHistoryPayment)
	if Err != nil {
		Sentry_log(Err)
		return Err
	}

	return nil
}
