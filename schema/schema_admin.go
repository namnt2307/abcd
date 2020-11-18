package schema

import "gopkg.in/mgo.v2/bson"

type RibbonsV3ObjectStruct struct {
	Id                 string             `json:"id" `
	Name               string             `json:"name" `
	Name_filter        string             `json:"name_filter" `
	Slug               string             `json:"slug" `
	Slug_filter        string             `json:"slug_filter" `
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
	Image_soucre []ImageSoucreObjectStruct `json:"image_soucre" `
	Images       ImageAllPlatformStruct    `json:"images" `
	Description  string                    `json:"description" `
}

type PropertiesStruct struct {
	Line        int `json:"line" `
	Is_title    int `json:"is_title" `
	Is_slide    int `json:"is_slide" `
	Is_refresh  int `json:"is_refresh" `
	Is_view_all int `json:"is_view_all" `
	Thumb       int `json:"thumb" `
}

type RibbonsV3ObjectSyncStruct struct {
	Id                 string `json:"id" `
	Name               string `json:"name" `
	Slug               string `json:"slug" `
	Type               int    `json:"type" `
	Is_premium         int    `json:"is_premium" `
	Is_new             int    `json:"is_new" `
	Display_image_type int    `json:"display_image_type" `
	Status             int    `json:"status" `
	Geo_check          int    `json:"geo_check" `
	Odr                int    `json:"odr" `
	Menus              []struct {
		Id        string `json:"id" `
		Slug      string `json:"slug" `
		Odr       int    `json:"odr" `
		Platforms []int  `json:"platforms" `
	} `json:"menus" `
	Tags struct {
		Id   string `json:"id" `
		Name string `json:"name" `
		Type string `json:"type" `
		Slug string `json:"slug" `
		Seo  SeoObjectStruct
	} `json:"tags" `
	Platforms []int           `json:"platforms" `
	Seo       SeoObjectStruct `json:"seo" `
}

type MenuObjectStruct struct {
	Id        string `json:"id" `
	Slug      string `json:"slug" `
	Odr       int    `json:"odr" `
	Platforms []int  `json:"platforms" `
}

type TagsObjectStruct struct {
	Id   string `json:"id" `
	Name string `json:"name" `
	Type string `json:"type" `
	Slug string `json:"slug" `
	Seo  SeoObjectStruct
}

type ImageSoucreObjectStruct struct {
	Height          string `json:"height" `
	Width           string `json:"width" `
	Image_type      string `json:"image_type" `
	Image_name      string `json:"image_name" `
	Resolution_type string `json:"resolution_type" `
	Url             string `json:"url" `
}

type RibbonItemV3PageStruct struct {
	Items    []RibbonItemV3ObjectStruct `json:"items" `
	Metadata MetadataStruct             `json:"metadata" `
}
type MetadataStruct struct {
	Page  int `json:"page" `
	Limit int `json:"limit" `
	Total int `json:"total" `
}

type RibbonItemV3ObjectStruct struct {
	Rib_item_id        string                    `json:"rib_item_id" `
	Group_id           string                    `json:"group_id" `
	Ref_id             string                    `json:"ref_id" `
	Rib_ids            []string                  `json:"rib_ids" `
	Geo_check          int                       `json:"geo_check" `
	Display_image_type int                       `json:"display_image_type" `
	Title              string                    `json:"title" `
	Short_description  string                    `json:"short_description" `
	Status             int                       `json:"status" `
	Is_premium         int                       `json:"is_premium" `
	Is_new             int                       `json:"is_new" `
	Odr                float64                   `json:"odr" `
	Type               int                       `json:"type" `
	Type_name          string                    `json:"type_name" `
	Platforms          []int                     `json:"platforms" `
	Start_time         int                       `json:"start_time" `
	End_time           int                       `json:"end_time" `
	Subtitle           string                    `json:"subtitle" `
	More_info          interface{}               `json:"more_info" `
	Seo                SeoObjectStruct           `json:"seo" `
	Resolution         int                       `json:"resolution" `
	Avg_rate           float64                   `json:"avg_rate" `
	Total_rate         int                       `json:"total_rate" `
	Is_watchlater      bool                      `json:"is_watchlater" `
	Episode            int                       `json:"episode" `
	Current_episode    string                    `json:"current_episode" `
	Slug               string                    `json:"slug" `
	R_id               string                    `json:"r_id" `
	Share_url          string                    `json:"Share_url" `
	Share_url_seo      string                    `json:"Share_url_seo" `
	Image_soucre       []ImageSoucreObjectStruct `json:"image_soucre" `
	Images             ImageAllPlatformStruct    `json:"images" `
	Updated_at         int64                     `json:"updated_at" `
	Created_at         int64                     `json:"created_at" `
	Info_admin         struct {
		Admin_id_create string `json:"admin_id_create" `
		Admin_id        string `json:"admin_id" `
		Admin_note      string `json:"admin_note" `
	} `json:"info_admin" `
	Info_Ribbon []RibbonsV3ObjectStruct `json:"info_ribbon" `
	Link_play   struct {
		Hls_link_play  string `json:"hls_link_play" `
		Dash_link_play string `json:"dash_link_play" `
	} `json:"link_play" `
	Tags_display []string `json:"tags_display"`
	External_url       string                    `json:"external_url" `
	Allow_click           int                       `json:"allow_click" `
}

type RibbonItemV3UpsertStruct struct {
	Rib_item_id        string                    `json:"rib_item_id" `
	Group_id           string                    `json:"group_id" `
	Ref_id             string                    `json:"ref_id" `
	Rib_ids            []string                  `json:"rib_ids" `
	Geo_check          int                       `json:"geo_check" `
	Display_image_type int                       `json:"display_image_type" `
	Title              string                    `json:"title" `
	Short_description  string                    `json:"short_description" `
	Status             int                       `json:"status" `
	Is_premium         int                       `json:"is_premium" `
	Is_new             int                       `json:"is_new" `
	Odr                float64                   `json:"odr" `
	Type               int                       `json:"type" `
	Type_name          string                    `json:"type_name" `
	Platforms          []int                     `json:"platforms" `
	Start_time         int                       `json:"start_time" `
	End_time           int                       `json:"end_time" `
	Subtitle           string                    `json:"subtitle" `
	More_info          interface{}               `json:"more_info" `
	Seo                SeoObjectStruct           `json:"seo" `
	Resolution         int                       `json:"resolution" `
	Avg_rate           float64                   `json:"avg_rate" `
	Total_rate         int                       `json:"total_rate" `
	Is_watchlater      bool                      `json:"is_watchlater" `
	Episode            int                       `json:"episode" `
	Current_episode    string                    `json:"current_episode" `
	Slug               string                    `json:"slug" `
	R_id               string                    `json:"r_id" `
	Share_url          string                    `json:"Share_url" `
	Share_url_seo      string                    `json:"Share_url_seo" `
	Image_soucre       []ImageSoucreObjectStruct `json:"image_soucre" `
	Images             ImageAllPlatformStruct    `json:"images" `
	Updated_at         int64                     `json:"updated_at" `
	Created_at         int64                     `json:"created_at" `
	Info_admin         struct {
		Admin_id_create string `json:"admin_id_create" `
		Admin_id        string `json:"admin_id" `
		Admin_note      string `json:"admin_note" `
	} `json:"info_admin" `
	Link_play struct {
		Hls_link_play  string `json:"hls_link_play" `
		Dash_link_play string `json:"dash_link_play" `
	} `json:"link_play" `
	External_url       string                    `json:"external_url" `
	Allow_click           int                       `json:"allow_click" `

}

type NotificationObjectStruct struct {
	Id             string      `json:"id" `
	Title          string      `json:"title" `
	Content        string      `json:"content" `
	Url            string      `json:"url" `
	Expire         int64       `json:"expire" `
	Status         int         `json:"status" `
	Image          string      `json:"image" `
	More_info      interface{} `json:"more_info" `
	Entity_id      string      `json:"entity_id" `
	Episode_id     string      `json:"episode_id" `
	Entity_type    int         `json:"entity_type" `
	Created_at     int64       `json:"created_at" `
	Update_at      int64       `json:"update_at" `
	Is_push_all    int         `json:"is_push_all" `
	Is_coming_soon int         `json:"is_coming_soon" `
	User_ids       []string    `json:"user_ids" `
	Msg_ids        []int       `json:"msg_ids" `
	Platforms      []string    `json:"platforms" `
	Admin_info     struct {
		Admin_id         string `json:"admin_id" `
		Admin_name       string `json:"admin_name" `
		Ad_approved_id   string `json:"ad_approved_id" `
		Ad_approved_name string `json:"ad_approved_name" `
		Ad_approved_at   int64  `json:"ad_approved_at" `
		Ad_push_id       string `json:"ad_push_id" `
		Ad_push_name     string `json:"ad_push_name" `
		Ad_push_at       int64  `json:"ad_push_at" `
	} `json:"admin_info" `
}

type NotificatiosObjectStruct struct {
	Items    []NotificationObjectStruct `json:"items" `
	Metadata MetadataStruct             `json:"metadata" `
}

type QuestionObjectStruct struct {
	Id            bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Start_time    int64                  `json:"start_time"`
	Duration      int64                  `json:"duration"`
	Question      string                 `json:"question"`
	Channel_id    string                 `json:"channel_id"`
	Odr           int                    `json:"odr"`
	Question_type int64                  `json:"question_type"`
	Answers       []AnswerObjectStruct   `json:"answers"`
	Result        string                 `json:"result"`
	Status        int                    `json:"status"`
	More_info     MoreInfoQuestionObject `json:"more_info"`
	Is_approved   int                    `json:"is_approved"`
	Created_at    int64                  `json:"created_at" `
	Updated_at    int64                  `json:"updated_at" `
}

type QuestionObjectUpsertStruct struct {
	Start_time    int64                  `json:"start_time"`
	Duration      int64                  `json:"duration"`
	Question      string                 `json:"question"`
	Channel_id    string                 `json:"channel_id"`
	Question_type int64                  `json:"question_type"`
	Odr           int                    `json:"odr"`
	Answers       []AnswerObjectStruct   `json:"answers"`
	Result        string                 `json:"result"`
	Status        int                    `json:"status"`
	Is_approved   int                    `json:"is_approved"`
	More_info     MoreInfoQuestionObject `json:"more_info"`
	Created_at    int64                  `json:"created_at" `
	Updated_at    int64                  `json:"updated_at" `
}
type MoreInfoQuestionObject struct {
	Ad_created_id    string `json:"ad_created_id" `
	Ad_created_name  string `json:"ad_created_name" `
	Ad_approved_id   string `json:"ad_approved_id" `
	Ad_approved_name string `json:"ad_approved_name" `
	Ad_approved_at   int64  `json:"ad_approved_at" `
	Ad_push_id       string `json:"ad_push_id" `
	Ad_push_name     string `json:"ad_push_name" `
	Ad_push_at       int64  `json:"ad_push_at" `
}
type AnswerObjectStruct struct {
	Answer_id    string `json:"answer_id"`
	Answer_title string `json:"answer_title"`
}

type QuestionPageStruct struct {
	Items    []QuestionObjectStruct `json:"items" `
	Metadata MetadataStruct         `json:"metadata" `
}

type RankingQuestionStruct struct {
	Channel_id    string                    `json:"channel_id"`
	Total_players int                       `json:"total_players"`
	Top_ranking   []PlayerRankingInfoStruct `json:"top_ranking"`
	Created_at    int64                     `json:"created_at" `
	Updated_at    int64                     `json:"updated_at" `
}

type PlayerRankingInfoStruct struct {
	User_id       string `json:"user_id"`
	Total_success int    `json:"total_success"`
	Mobile        string `json:"mobile"`
	User_name     string `json:"user_name"`
}

type SimpleInfoUser struct {
	Id        string `json:"id`
	Mobile    string `json:"mobile"`
	User_name string `json:"user_name"`
}

type CommentPageStruct struct {
	Items    []CommentObjectStruct `json:"items" `
	Metadata struct {
		Page  int `json:"page" `
		Limit int `json:"limit" `
		Total int `json:"total" `
	} `json:"metadata" `
}

// DT-10141
type BadwordPageStruct struct {
	Items    []BadwordObjectStruct `json:"items" `
	Metadata struct {
		Page  int `json:"page" `
		Limit int `json:"limit" `
		Total int `json:"total" `
	} `json:"metadata" `
}
