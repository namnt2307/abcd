package question

import (
	"math"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func QuestionStatistics(id string) (QuestionStatisticsStruct, error) {
	var Question QuestionStatisticsStruct
	// Connect MongoDB
	session, db, err := GetCollectionDmDB()
	if err != nil {
		return Question, err
	}
	defer session.Close()

	var where = bson.M{"_id": bson.ObjectIdHex(id)}

	err = db.C(COLLECTION_DM_QUESTIONS).Find(where).One(&Question)
	if err != nil {
		return Question, err
	}

	var StatisticsTotalAnswers []AnswerStatisticsStruct
	var query = []bson.M{
		{"$match": bson.M{"question_id": id, "channel_id": Question.Channel_id}},
		{
			"$group": bson.M{
				"_id":         bson.M{"client_answer_id": "$client_answer_id"},
				"answer_id":   bson.M{"$first": "$client_answer_id"},
				"total_users": bson.M{"$sum": 1}},
		},
	}
	pipe := db.C(COLLECTION_DM_USER_ANSWER).Pipe(query)
	err = pipe.All(&StatisticsTotalAnswers)
	if err != nil && err.Error() != "not found" {
		return Question, err
	}
	var total_answer float64 = 0
	for _, val := range StatisticsTotalAnswers {
		total_answer = total_answer + float64(val.Total_users)
	}

	for _, val := range StatisticsTotalAnswers {
		for i, answer := range Question.Answers {
			if answer.Answer_id == val.Answer_id {
				Question.Answers[i].Total_users = val.Total_users
				Question.Answers[i].Percent_users = math.RoundToEven(float64(val.Total_users) * 100 / total_answer)
				break
			}
		}
	}
	return Question, nil
}
