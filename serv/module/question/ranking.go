package question

import (
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetRankingByChannel(channel_id string) (RankingQuestionStruct, error) {
	var Ranking RankingQuestionStruct
	// Connect MongoDB
	session, db, err := GetCollectionDmDB()
	if err != nil {
		return Ranking, err
	}
	defer session.Close()

	err = db.C(COLLECTION_DM_RANKING).Find(bson.M{"channel_id": channel_id}).Sort("created_at").One(&Ranking)
	if err != nil && err.Error() != "not found" {
		return Ranking, err
	}
	return Ranking, nil
}
