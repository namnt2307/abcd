package tracking

const (
	TRACKING_WATCH_NEXT_LOG = 120
	TRACKING_WATCH_RECORD_TIME = 5
)

type DataResponseStruct struct {
	Success bool `json:"success" `
}

type TrackingWatchResultStruct struct {
	Message string `json:"message" `
	Data    struct {
		Next_log    int    `json:"next_log" `
		Token       string `json:"token" `
		Record_time int    `json:"record_time" `
	} `json:"data" `
	Content_info    struct {
		Id                string `json:"id" `
		Group_id          string `json:"group_id" `
		Type              int `json:"type" `
		Is_premium        int `json:"is_premium" `
		Title             string `json:"title" `
		Movie           struct {
			Title string `json:"title" `
		} `json:"movie" `
		Content_provider_id string `json:"content_provider_id" `
		Tags []struct {
			Id   string          `json:"id" `
			Name string          `json:"name" `
			Type string          `json:"type" `
		} `json:"tags" `
	} `json:"content_info" `
	Success bool `json:"success" `
}

type TrackingWatchRequestDataStruct struct {
	Content_id   string
	Content_type int
	Content_name string
	Data         []struct {
		Action    string
		Duration  int64
		Progress  int64
		Timestamp int64
	}
	Usi string
}