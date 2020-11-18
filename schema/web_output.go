package schema

// "gopkg.in/mgo.v2/bson"

type RibbonDataOutputObjectStruct struct {
	Id           string
	Slug         string
	Name         string
	Type         int
	Is_premium   int `json:"is_premium" `
	Is_new       int `json:"is_new" `
	Ribbon_items []VODDataObjectStruct
}

type RibbonDetailOutputObjectStruct struct {
	PageRibbonsOutputObjectStruct
	Metadata struct {
		Total int `json:"total" `
		Limit int `json:"limit" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type PageRibbonsOutputObjectStruct struct {
	Id                 string `json:"id" `
	Name               string `json:"name" `
	Type               int    `json:"type" `
	Display_image_type int    `json:"display_image_type" `
	Is_premium         int    `json:"is_premium" `
	Menus              []struct {
		Id string `json:"id" `
	} `json:"menus" `
	Tags []struct {
		Id   string          `json:"id" `
		Name string          `json:"name" `
		Type string          `json:"type" `
		Slug string          `json:"slug" `
		Seo  SeoObjectStruct `json:"seo" `
	} `json:"tags" `
	// Tag struct {
	// 	Id   string          `json:"id" `
	// 	Name string          `json:"name" `
	// 	Type string          `json:"type" `
	// 	Slug string          `json:"slug" `
	// 	Seo  SeoObjectStruct `json:"seo" `
	// } `json:"tag" `
	Ribbon_items  []RibbonItemOutputObjectStruct `json:"items" `
	Tracking_data TrackingDataStruct             `json:"tracking_data" `
	Seo           SeoObjectStruct                `json:"seo" `
	Geo_check     int                            `json:"geo_check" `
	Properties    PropertiesStruct               `json:"properties" `
	Images        ImagesOutputObjectStruct       `json:"images" `
	Description   string                         `json:"description" `
}

type RibbonItemOutputObjectStruct struct {
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
	Geo_check       int                      `json:"geo_check" `
	Type            int                      `json:"type" `
	Title           string                   `json:"title" `
	Genre           string                   `json:"genre" `
	Resolution      int                      `json:"resolution" `
	Images          ImagesOutputObjectStruct `json:"images" `
	Avg_rate        float64                  `json:"avg_rate" `
	Total_rate      int                      `json:"total_rate" `
	Is_watchlater   bool                     `json:"is_watchlater" `
	Is_premium      int                      `json:"is_premium" `
	Is_new          int                      `json:"is_new" `
	Episode         int                      `json:"episode" `
	Current_episode string                   `json:"current_episode" `
	Seo             SeoObjectStruct          `json:"seo" `
	Slug            string                   `json:"slug" `
	// R_id               string `json:"r_id" `
	Display_image_type int `json:"display_image_type" `
	// Type_name          string                   `json:"type_name" `
	More_info interface{} `json:"more_info,omitempty" `

	Start_time           int64  `json:"start_time,omitempty" `
	End_time             int64  `json:"end_time,omitempty" `
	Subtitle             string `json:"subtitle,omitempty" `
	Share_url            string `json:"share_url" `
	Share_url_seo        string `json:"share_url_seo" `
	Label_subtitle_audio string `json:"label_subtitle_audio" `
	Label_public_day     string `json:"label_public_day" `
	Is_coming_soon       int    `json:"is_coming_soon" `
	Link_play            struct {
		Hls_link_play  string `json:"hls_link_play" `
		Dash_link_play string `json:"dash_link_play" `
	} `json:"link_play" `
	Tags   []TagObjectStruct `json:"tags"`
	People struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Release_year int      `json:"release_year" `
	Tags_display []string `json:"tags_display" `
	Ranking      int `json:"ranking" `
	External_url       string                    `json:"external_url" `
	Allow_click           int                    `json:"allow_click" `
}

type BannerRibbonsOutputObjectStruct struct {
	Id                string                   `json:"id" `
	Group_id          string                   `json:"group_id" `
	Title             string                   `json:"title" `
	Short_description string                   `json:"short_description" `
	Images            ImagesOutputObjectStruct `json:"images" `
	Seo               SeoObjectStruct          `json:"seo" `
	Type              int                      `json:"type" `
	More_info         interface{}              `json:"more_info,omitempty" `
	Start_time        int64                    `json:"start_time,omitempty" `
	End_time          int64                    `json:"end_time,omitempty" `
	Subtitle          string                   `json:"subtitle,omitempty" `
	Is_premium        int                      `json:"is_premium" `
	Is_new            int                      `json:"is_new" `
	Geo_check         int                      `json:"geo_check" `
	Link_play         struct {
		Hls_link_play  string `json:"hls_link_play" `
		Dash_link_play string `json:"dash_link_play" `
	} `json:"link_play" `
	Label_subtitle_audio string            `json:"label_subtitle_audio" `
	Label_public_day     string            `json:"label_public_day" `
	Resolution           int               `json:"resolution" `
	Tags                 []TagObjectStruct `json:"tags"`
	Release_year         int               `json:"release_year"`
	Avg_rate             float64           `json:"avg_rate" `
	Total_rate           int               `json:"total_rate" `

	People struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags_display []string `json:"tags_display" `
	Ranking      int      `json:"ranking" `
}

type VodDataOutputObjectStruct struct {
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
	Geo_check       int                      `json:"geo_check" `
	Type            int                      `json:"type" `
	Title           string                   `json:"title" `
	Resolution      int                      `json:"resolution" `
	Images          ImagesOutputObjectStruct `json:"images" `
	Avg_rate        float64                  `json:"avg_rate" `
	Total_rate      int                      `json:"total_rate" `
	Is_watchlater   bool                     `json:"is_watchlater" `
	Is_premium      int                      `json:"is_premium" `
	Is_new          int                      `json:"is_new" `
	Episode         int                      `json:"episode" `
	Current_episode string                   `json:"current_episode" `
	Seo             SeoObjectStruct          `json:"seo" `
	Slug            string                   `json:"slug" `
	// R_id               string `json:"r_id" `
	Display_image_type int `json:"display_image_type" `
	// Type_name          string                   `json:"type_name" `
	More_info interface{} `json:"more_info,omitempty" `

	Start_time           int64  `json:"start_time,omitempty" `
	End_time             int64  `json:"end_time,omitempty" `
	Subtitle             string `json:"subtitle,omitempty" `
	Share_url            string `json:"share_url" `
	Share_url_seo        string `json:"share_url_seo" `
	Label_subtitle_audio string `json:"label_subtitle_audio" `
	Label_public_day     string `json:"label_public_day" `
	Is_coming_soon       int    `json:"is_coming_soon" `
}

type ImagesOutputObjectStruct struct {
	Home_vod_hot      string `json:"home_vod_hot" `
	Vod_thumb_big     string `json:"vod_thumb_big" `
	Banner            string `json:"banner" `
	Home_carousel_tv  string `json:"home_carousel_tv" `
	Home_carousel_web string `json:"home_carousel_web" `
	Home_carousel     string `json:"home_carousel" `
	Thumbnail         string `json:"thumbnail" `
	Vod_thumb         string `json:"vod_thumb" `
	Poster            string `json:"poster" `
	Big_thumb         string `json:"big_thumb" `

	Thumbnail_hot_v4 string `json:"thumbnail_hot_v4" `
	Thumbnail_big_v4 string `json:"thumbnail_big_v4" `
	Carousel_tv_v4   string `json:"carousel_tv_v4" `
	Carousel_app_v4  string `json:"carousel_app_v4" `
	Carousel_web_v4  string `json:"carousel_web_v4" `

	Thumbnail_v4     string `json:"thumbnail_v4" `
	Poster_v4        string `json:"poster_v4" `
	Promotion_banner string `json:"promotion_banner" `
	Title_card_light string `json:"title_card_light" `
	Title_card_dark  string `json:"title_card_dark" `
	Poster_original  string `json:"poster_original" `
	Thumb_original   string `json:"thumb_original" `
}

type MenuOutputObjectStruct struct {
	Id   string `json:"id" `
	Name string `json:"name" `
	Slug string `json:"slug" `
	// Required bool   `json:"required" `
	Icon      string `json:"icon" `
	Icon_text string `json:"icon_text" `
	Tag       struct {
		Id   string          `json:"id" `
		Name string          `json:"name" `
		Type string          `json:"type" `
		Slug string          `json:"slug" `
		Seo  SeoObjectStruct `json:"seo" `
	} `json:"tag" `
	Sub_menu []MenuOutputObjectStruct `json:"sub_menu" `
	// Layout_types string                   `json:"layout_types" `
	Seo             SeoObjectStruct          `json:"seo" `
	Parent_id       string                   `json:"parent_id" `
	Have_banner     int                      `json:"have_banner"`
	Title_ribbon    string                   `json:"title_ribbon,omitempty"`
	Quantity_ribbon int                      `json:"quantity_ribbon,omitempty"`
	Sub_menu_ribbon []MenuRibbonObjectStruct `json:"sub_menu_ribbon,omitempty" `
}
