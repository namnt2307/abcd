package middleware

const (
	PREFIX_LOG_VOD_VIEW 				= "LOG_VOD_VIEW_%s_%s"
	PREFIX_LOG_LIVE_TV_VIEW 			= "LOG_LIVE_TV_VIEW_%s_%s_%s"
	PREFIX_LOG_LIVE_STREAM_VIEW 		= "LOG_LIVE_STREAM_VIEW_%s_%s"
	PREFIX_LOG_USC_ONLINE_PER_HOUS      = "LOG_USC_ONLINE_PER_HOUS_%s_%s"
	PREFIX_LOG_TOTAL_ONLINE_PER_HOUS    = "LOG_TOTAL_ONLINE_PER_HOUS_%s"
)

type DataInputLogStruct struct {
	Action     string  `json:"action" `
	Uri        string  `json:"uri" `
	Method     string  `json:"method" `
	User_agent string  `json:"user_agent" `
	Ip         string  `json:"ip" `
	Time       string  `json:"time" `
	Params     string  `json:"params" `
	Data       string  `json:"data" `
	Data_response       string  `json:"data_response" `
	Data_response_str   string  `json:"data_response_str" `
	Duration   float64 `json:"duration" `
	ClientIP   string  `json:"clientIP" `
	Status     int     `json:"status" `
	Platform 	string  `json:"platform" `
	User_id 	string  `json:"user_id" `
	User_is_premium 	int  `json:"user_is_premium" `
	Token_access string `json:"token_access" `
}


type DataInputLog_VOD_Struct struct {
	Id 						string `json:"id" `
	Group_id 				string `json:"group_id" `
	Is_premium 				int `json:"is_premium" `
	Type 					int `json:"type" `
	Title 					string `json:"title" `
	Episode 				int `json:"episode" `
	Content_provider_id 	string `json:"content_provider_id" `
	Movie struct{
		Title string `json:"title" `
	} `json:"movie" `
	Tags 					[]struct{
		Id 		string `json:"id" `
		Name 	string `json:"name" `
		Type 	string `json:"type" `
	} `json:"tags" `
}

type DataInputLog_LIVETV_Struct struct {
	Technical struct {
		Id                  string `json:"id" `
		MainChannelId       string `json:"mainChannelId" `
		Title               string
		Categories       	[]string `json:"categorys" `
	} `json:"technical" `
}


