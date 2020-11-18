package schema

type PlatformOutputStruct struct {
	Id   int
	Type string
}

type VodTipDataOutputObjectStruct struct {
	Id                string `json:"id" `
	Title             string `json:"title" `
	Short_description string `json:"short_description" `
	Resolution        int    `json:"resolution" `
	Is_watchlater     bool   `json:"is_watchlater" `
	Audios            []struct {
		Id         string `json:"id" `
		Code_name  string `json:"code_name" `
		Is_default int64  `json:"is_default" `
		Index      string `json:"index" `
		Title      string `json:"title" `
	} `json:"audios" `
	Seo        SeoObjectStruct `json:"seo" `
	Is_premium int             `json:"is_premium" `
	Is_new     int             `json:"is_new" `
}

type WatchlaterContentObjStruct struct {
	Items    []WatchlaterContentDetailObjStruct `json:"items" `
	Metadata struct {
		Total int64 `json:"total" `
		Limit int   `json:"limit" `
		Page  int   `json:"page" `
	} `json:"metadata" `
}

type WatchlaterContentDetailObjStruct struct {
	Avg_rate        float64                  `json:"avg_rate" `
	Current_episode string                   `json:"current_episode" `
	Episode         int                      `json:"episode" `
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
	Images          ImagesOutputObjectStruct `json:"images" `
	Is_watchlater   bool                     `json:"is_watchlater" `
	Resolution      int                      `json:"resolution" `
	Runtime         int                      `json:"runtime" `
	Slug            string                   `json:"slug" `
	Type            int                      `json:"type" `
	Title           string                   `json:"title" `
	Total_rate      int                      `json:"total_rate" `
	Movie           struct {
		Title string `json:"title" `
	} `json:"movie" `
	Seo        SeoObjectStruct `json:"seo" `
	Is_premium int             `json:"is_premium" `
	Is_new     int             `json:"is_new" `
}

type TagInfoOutputObjectStruct struct {
	Id   string          `json:"id" `
	Name string          `json:"name" `
	Type string          `json:"type" `
	Seo  SeoObjectStruct `json:"seo,omitempty" `
	Slug string          `json:"slug" `
}

type TagsItemsObjectStruct struct {
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
	Title           string                   `json:"title" `
	Type            int                      `json:"type" `
	Resolution      int                      `json:"resolution" `
	Slug            string                   `json:"slug" `
	Avg_rate        float64                  `json:"avg_rate" `
	Total_rate      int                      `json:"total_rate" `
	Is_watchlater   bool                     `json:"is_watchlater" `
	Images          ImagesOutputObjectStruct `json:"images" `
	Seo             SeoObjectStruct          `json:"seo" `
	Slug_seo        string                   `json:"slug_seo" `
	Is_premium      int                      `json:"is_premium" `
	Is_new          int                      `json:"is_new" `
	Current_episode string                   `json:"current_episode" `
	Episode         int                      `json:"episode" `
	Geo_check       int                      `json:"geo_check" `
	People          struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags_display []string `json:"tags_display" `
	Ranking      int      `json:"ranking" `
}

type TagsOutputObjectStruct struct {
	Items []TagsItemsObjectStruct `json:"items" `

	Tracking_data TrackingDataStruct `json:"tracking_data" `
	Metadata      struct {
		Tags  []TagInfoOutputObjectStruct `json:"tags" `
		Limit int                         `json:"limit" `
		Seo   SeoObjectStruct             `json:"seo" `
		Total int                         `json:"total" `
		Page  int                         `json:"page" `
	} `json:"metadata" `
}

type ArtistOutputObjectStruct struct {
	Id       string `json:"id" `
	Name     string `json:"name" `
	Slug     string `json:"slug" `
	Gender   int    `json:"gender" `
	Info     string `json:"info" `
	Birthday string `json:"birthday" `
	Images   struct {
		Avatar string `json:"avatar" `
	} `json:"images" `
	Country struct {
		Id   string `json:"id" `
		Name string `json:"name" `
	} `json:"country" `
	Job string `json:"job" `
	// Seo SeoObjectStruct `json:"seo" `
}

type ItemArtistOutputObjectStruct struct {
	Id     string `json:"id" `
	Name   string `json:"name" `
	Images struct {
		Avatar string `json:"avatar" `
	} `json:"images" `
	Seo SeoObjectStruct `json:"seo" `
}

type ArtistRelatedOutputObjectStruct struct {
	Seo      SeoObjectStruct                `json:"seo" `
	Items    []ItemArtistOutputObjectStruct `json:"items" `
	Metadata struct {
		Limit int `json:"limit" `
		Total int `json:"total" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type ItemsVodArtistOutputStruct struct {
	Id            string                   `json:"id" `
	Group_id      string                   `json:"group_id" `
	Type          int                      `json:"type" `
	Title         string                   `json:"title" `
	Resolution    int                      `json:"resolution" `
	Slug          string                   `json:"slug" `
	Slug_seo      string                   `json:"slug_seo" `
	Avg_rate      float64                  `json:"avg_rate" `
	Total_rate    int                      `json:"total_rate" `
	Images        ImagesOutputObjectStruct `json:"images" `
	Is_watchlater bool                     `json:"is_watchlater" `
	Seo           SeoObjectStruct          `json:"seo" `
	Is_premium    int                      `json:"is_premium" `
	Geo_check     int                      `json:"geo_check" `
	Is_new        int                      `json:"is_new" `
	People        struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags_display []string `json:"tags_display" `
	Ranking      int      `json:"ranking" `
}

type VodArtistOutputObjectStruct struct {
	Seo      SeoObjectStruct              `json:"seo" `
	Items    []ItemsVodArtistOutputStruct `json:"items" `
	Metadata struct {
		Limit int `json:"limit" `
		Total int `json:"total" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type ItemWatchingOutputObjectStruct struct {
	Avg_rate        float64                  `json:"avg_rate" `
	Episode         int                      `json:"episode" `
	Current_episode string                   `json:"current_episode" `
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
	Images          ImagesOutputObjectStruct `json:"images" `
	Is_watchlater   bool                     `json:"is_watchlater" `
	Is_premium      int                      `json:"is_premium" `
	Is_new          int                      `json:"is_new" `
	Movie           struct {
		Title string `json:"title" `
	} `json:"movie" `
	Progress     int64           `json:"progress" `
	Runtime      int             `json:"runtime" `
	Resolution   int             `json:"resolution" `
	Title        string          `json:"title" `
	Type         int             `json:"type" `
	Slug         string          `json:"slug" `
	Total_rate   int             `json:"total_rate" `
	Seo          SeoObjectStruct `json:"seo" `
	Tags_display []string        `json:"tags_display" `
	Ranking      int             `json:"ranking" `
}

type WatchingOutputObjectStruct struct {
	Items []ItemWatchingOutputObjectStruct `json:"items" `

	Tracking_data TrackingDataStruct `json:"tracking_data" `
	Metadata      struct {
		Limit int `json:"limit" `
		Total int `json:"total" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type MessageItemOutputObjectStruct struct {
	Id          string          `json:"id" `
	Title       string          `json:"title" `
	Content     string          `json:"content" `
	Expire      int64           `json:"expire" `
	Status      int             `json:"status" ` // 0: Chua xem , 1: da xem , 2: xoa
	Image       string          `json:"image" `
	More_info   interface{}     `json:"more_info" `
	Entity_type int             `json:"entity_type" `
	Entity_id   string          `json:"entity_id" `
	Episode_id  string          `json:"episode_id" `
	Created_at  int64           `json:"created_at" `
	Seo         SeoObjectStruct `json:"seo" `
	Is_push_all int             `json:"is_push_all" `
	User_ids    []string        `json:"user_ids" `
}

type MessageOutputObjectStruct struct {
	Items    []MessageItemOutputObjectStruct `json:"items" `
	Metadata struct {
		Limit int `json:"limit" `
		Total int `json:"total" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type SearchResultStruct struct {
	Seo           SeoObjectStruct          `json:"seo" `
	Items         []SearchItemResultStruct `json:"items" `
	Tracking_data TrackingDataStruct       `json:"tracking_data" `
	Metadata      struct {
		Limit   int    `json:"limit" `
		Total   int64  `json:"total" `
		Page    int    `json:"page" `
		Keyword string `json:"keyword" `
	} `json:"metadata" `
}

type SearchItemResultStruct struct {
	Is_artist       bool            `json:"is_artist" `
	Episode         int             `json:"episode,omitempty" `
	Avg_rate        float64         `json:"avg_rate,omitempty" `
	Title           string          `json:"title,omitempty" `
	Resolution      int             `json:"resolution,omitempty" `
	Total_rate      int             `json:"total_rate,omitempty" `
	Current_episode string          `json:"current_episode,omitempty" `
	Slug            string          `json:"slug,omitempty" `
	Images          interface{}     `json:"images,omitempty" `
	Is_watchlater   bool            `json:"is_watchlater,omitempty" `
	Seo             SeoObjectStruct `json:"seo,omitempty" `
	Group_id        string          `json:"group_id,omitempty" `
	Type            int             `json:"type,omitempty" `
	Id              string          `json:"id,omitempty" `
	Slug_seo        string          `json:"slug_seo,omitempty" `
	Name            string          `json:"name,omitempty" `
	Job             string          `json:"job,omitempty" `
	Is_premium      int             `json:"is_premium" `
	Geo_check       int             `json:"geo_check" `
	Is_new          int             `json:"is_new" `
	People          struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags_display []string `json:"tags_display" `
	Ranking      int      `json:"ranking" `
	Request_id   string   `json:"request_id"`
	Position     string   `json:"position"`
}

type SearchSuggestOutputStruct struct {
	Items         []string           `json:"items" `
	Tracking_data TrackingDataStruct `json:"tracking_data" `
}

type SearchTopKeywordOutputStruct struct {
	Items         []SearchTopKeywordResultStruct `json:"items" `
	Tracking_data TrackingDataStruct             `json:"tracking_data" `
	// Metadata      struct {
	// 	Limit   int    `json:"limit" `
	// 	Total   int64  `json:"total" `
	// 	Page    int    `json:"page" `
	// 	Keyword string `json:"keyword" `
	// } `json:"metadata" `
}

type SearchTopKeywordResultStruct struct {
	Id           string `json:"id" `
	Keyword      string `json:"keyword" `
	Search_count int    `json:"search_count" `
}

type LiveEventOutputObjectStruct struct {
	Id                string `json:"id" `
	Title             string `json:"title" `
	Short_description string `json:"short_description" `
	Long_description  string `json:"long_description" `
	Location          string `json:"location" `
	Link_play         struct {
		Hls_link_play  string `json:"hls_link_play" `
		Dash_asset_id  string `json:"dash_asset_id" `
		Hls_asset_id   string `json:"hls_asset_id" `
		Dash_link_play string `json:"dash_link_play" `
	} `json:"link_play" `
	Images      ImagesOutputObjectStruct `json:"images" `
	Interaction []struct {
		Id   string `json:"id" `
		Name string `json:"name" `
	} `json:"interaction" `
	Is_live     int `json:"is_live" `
	Socket_info struct {
		Tcp_port string `json:"tcp_port" `
		Host     string `json:"host" `
		Port     string `json:"port" `
	} `json:"socket_info" `
	Start_time      string                     `json:"start_time" `
	Str_to_time     int64                      `json:"str_to_time" `
	Share_url       string                     `json:"share_url" `
	Slug            string                     `json:"slug" `
	Seo             SeoObjectStruct            `json:"seo" `
	Share_url_seo   string                     `json:"share_url_seo" `
	Permission      int                        `json:"permission" `
	PackageGroup    []PackageGroupObjectStruct `json:"packages" `
	Type            int                        `json:"type" `
	Group_id        string                     `json:"group_id" `
	Total_award     float64                    `json:"total_award" `
	Nextevent_info  NexteventInfoStruct        `json:"nextevent_info" `
	Button_live     string                     `json:"button_live" `
	Button_cap_love string                     `json:"button_cap_love" `
}

type LiveEventFinishedOutputObjectStruct struct {
	Items    []VodDataOutputObjectStruct `json:"items" `
	Metadata struct {
		Page  int `json:"page" `
		Limit int `json:"limit" `
		Total int `json:"total" `
	} `json:"metadata" `
}

type TransactionUserObjectStruct struct {
	Items    []TrancactionObjectStruct `json:"items" `
	Metadata struct {
		Page  int `json:"page" `
		Limit int `json:"limit" `
		Total int `json:"total" `
	} `json:"metadata" `
}

type TrancactionObjectStruct struct {
	Name           string `json:"name" `
	Payment_method string `json:"payment_method" `
	Status         string `json:"status" `
	Created_date   string `json:"created_date" `
}

type TransactionUserSuccessObjectStruct struct {
	Items    []TransactionSuccessObjectStruct `json:"items" `
	Metadata struct {
		Total int `json:"total" `
		Limit int `json:"limit" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type TransactionSuccessObjectStruct struct {
	Id           int    `json:"id" `
	Name         string `json:"name" `
	Expired_date string `json:"expired_date" `
	Avatar       string `json:"image" `
}

type CommentOutputStruct struct {
	Id         string                  `json:"id"`
	Content_id string                  `json:"content_id"`
	Message    string                  `json:"message"`
	Created_at int64                   `json:"created_at" `
	Reply      []CommentOutputStruct   `json:"reply,omitempty"`
	User       UserCommentOutputStruct `json:"user"`
	// Thu add DT-9922
	Status int `json:"status"`
	Pin    int `json:"pin"`
}

type CommentOutputPageStruct struct {
	Items    []CommentOutputStruct `json:"items" `
	Metadata struct {
		Page  int `json:"page" `
		Limit int `json:"limit" `
		Total int `json:"total" `
	} `json:"metadata" `
}

type UserCommentOutputStruct struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int    `json:"gender"`
}

type BadwordOutputPageStruct struct {
	Items    []BadwordOutputStruct `json:"items" `
	Metadata struct {
		Page  int `json:"page" `
		Limit int `json:"limit" `
		Total int `json:"total" `
	} `json:"metadata" `
}

type BadwordOutputStruct struct {
	Id         string `json:"id"`
	Content_id string `json:"content_id"`
	Message    string `json:"message"`
	Created_at int64  `json:"created_at" `
}

type RibbonsV3OutputStruct struct {
	Id                 string             `json:"id" `
	Name               string             `json:"name" `
	Slug               string             `json:"slug" `
	Type               int                `json:"type" `
	Display_image_type int                `json:"display_image_type" `
	Status             int                `json:"status" `
	Geo_check          int                `json:"geo_check" `
	Is_premium         int                `json:"is_premium" `
	Is_new             int                `json:"is_new" `
	Odr                int                `json:"odr" `
	Menus              []MenuObjectStruct `json:"menus" `
	Tags               []TagsObjectStruct `json:"tags" `
	Platforms          []int              `json:"platforms" `
	Seo                SeoObjectStruct    `json:"seo" `
	Properties         struct {
		Web     PropertiesStruct `json:"web" `
		App     PropertiesStruct `json:"app" `
		Smarttv PropertiesStruct `json:"smarttv" `
	} `json:"properties" `
	Images ImagesOutputObjectStruct `json:"images" `
}

type ConfigOutputObjectStruct struct {
	Data struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"data"`
}

// DT-13431
type RateObjectStruct struct {
	ContentId   string  `json:"contentId" `
	Avg_rate    float64 `json:"avg_rate" `
	Total_rate  int     `json:"total_rate" `
	Total_point float64 `json:"total_point" `
}
