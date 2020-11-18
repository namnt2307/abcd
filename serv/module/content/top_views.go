package content

import (
	. "cm-v5/serv/module"
	"time"
)

func InitContentTopViews() {
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return
	}
	defer session.Close()

	for true {
		var TopViewsContent []struct {
			Id_content string
		}

		err = db.C("top_view_content").Find(nil).Limit(10).Sort("-ranking").All(&TopViewsContent)
		if err != nil {
			return
		}

		for i, val := range TopViewsContent {
			LIST_TOP_VIEW_CONTENT[val.Id_content] = i + 1
		}

		time.Sleep(300 * time.Second)
	}
}
