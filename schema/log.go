package schema

const (
	COLLECTION_ENTITY_LOG          	= "entity_log"
	COLLECTION_USC_LOG          	= "usc_log"
	COLLECTION_USC_ONLINE_LOG       = "usc_online_log"
)

type EntityLogObjectStruct struct {
	Entity_id         	string 
	Entity_name         string
	Group_id          	string 
	Group_name         	string 
	Entity_type 		int
	Web_view 			int64
	App_view 			int64
	Tv_view 			int64
	Num_view 			int64
	Time_view 			int64
	Date 				string
}

type USCLogObjectStruct struct {
	Num_phone 			int64
	Num_email 			int64
	Num_facebook 		int64
	Num_google 			int64
	Num_total 			int64
	Date 				string
}

type USCOnlineLogObjectStruct struct {
	Date 				string // 2019-05-12
	Day_str 			string // Sun Mon Tue Wed Thu Fri Sat
	Detail_hour 		map[string]int64 // "11am: 100   4pm: 1000"
}

