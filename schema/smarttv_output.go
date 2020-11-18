package schema

type LinkPlayStruct struct {
	Dash_asset_id       string `json:"dash_asset_id" `
	Hls_asset_id        string `json:"hls_asset_id" `
	Dash_link_play      string `json:"dash_link_play" `
	Hls_link_play       string `json:"hls_link_play" `
	H265_dash_link_play string `json:"h265_dash_link_play,omitempty" `
	H265_hls_link_play  string `json:"h265_hls_link_play,omitempty" `
	Map_profile         struct {
		Four_k struct {
			Width      int    `json:"width" `
			Bandwidth  int    `json:"bandwidth" `
			Name       string `json:"name" `
			Height     int    `json:"height" `
			Is_premium int    `json:"is_premium" `
			Permission int    `json:"permission" `
		} `json:"four_k" `
		Full_hd struct {
			Width      int    `json:"width" `
			Bandwidth  int    `json:"bandwidth" `
			Name       string `json:"name" `
			Height     int    `json:"height" `
			Is_premium int    `json:"is_premium" `
			Permission int    `json:"permission" `
		} `json:"full_hd" `
		Hd struct {
			Width      int    `json:"width" `
			Bandwidth  int    `json:"bandwidth" `
			Name       string `json:"name" `
			Height     int    `json:"height" `
			Is_premium int    `json:"is_premium" `
			Permission int    `json:"permission" `
		} `json:"hd" `
		Sd struct {
			Width      int    `json:"width" `
			Bandwidth  int    `json:"bandwidth" `
			Name       string `json:"name" `
			Height     int    `json:"height" `
			Is_premium int    `json:"is_premium" `
			Permission int    `json:"permission" `
		} `json:"sd" `
	} `json:"map_profile" `
}

type AdsOutputStruct struct {
	Url    string `json:"url" `
	Repeat int    `json:"repeat" `
	Type   string `json:"type" `
	Group  string `json:"group" `
}

type PeopleOutputStruct struct {
	Id     string `json:"id" `
	Name   string `json:"name" `
	Images struct {
		Avatar string `json:"avatar" `
	} `json:"images" `
	Seo SeoObjectStruct `json:"seo" `
}

type SubtitlesOuputStruct struct {
	Id         string `json:"id" `
	Code_name  string `json:"code_name" `
	Is_default int64  `json:"is_default" `
	Uri        string `json:"uri" `
	Title      string `json:"title" `
	Is_premium int    `json:"is_premium" `
	Permission int    `json:"permission" `
}

type AudiosOutputStruct struct {
	Id         string `json:"id" `
	Code_name  string `json:"code_name" `
	Is_default int64  `json:"is_default" `
	Index      string `json:"index" `
	Title      string `json:"title" `
	Is_premium int    `json:"is_premium" `
	Permission int    `json:"permission" `
}

type ContentObjOutputStruct struct {
	Id                string            `json:"id" `
	Group_id          string            `json:"group_id" `
	Type              int               `json:"type" `
	Title             string            `json:"title" `
	Ads               []AdsOutputStruct `json:"ads" `
	Custom_ads        string            `json:"custom_ads" `
	Enable_ads        int               `json:"enable_ads" `
	Short_description string            `json:"short_description" `
	Long_description  string            `json:"long_description" `
	Resolution        int               `json:"resolution" `
	Runtime           int               `json:"runtime" `
	Is_premium        int               `json:"is_premium" `
	Is_new            int               `json:"is_new" `
	Link_play         LinkPlayStruct    `json:"link_play" `
	People            struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags []struct {
		Id   string          `json:"id" `
		Name string          `json:"name" `
		Type string          `json:"type" `
		Seo  SeoObjectStruct `json:"seo" `
		Slug string          `json:"slug" `
	} `json:"tags" `
	Subtitles           []SubtitlesOuputStruct   `json:"subtitles,omitempty" `
	Audios              []AudiosOutputStruct     `json:"audios,omitempty" `
	Release_year        int                      `json:"release_year" `
	Player_logo         string                   `json:"player_logo" `
	Player_logo_4k      string                   `json:"player_logo_4k" `
	Images              ImagesOutputObjectStruct `json:"images" `
	Progress            int64                    `json:"progress" `
	Avg_rate            float64                  `json:"avg_rate" `
	Total_rate          int                      `json:"total_rate" `
	Is_watchlater       bool                     `json:"is_watchlater" `
	User_rating         int                      `json:"user_rating" `
	Episode             int                      `json:"episode" `
	Current_episode     string                   `json:"current_episode" `
	Content_provider_id string                   `json:"content_provider_id" `
	Seo                 SeoObjectStruct          `json:"seo" `
	Slug                string                   `json:"slug" `
	Season_slug_seo     string                   `json:"season_slug_seo"`
	Slug_seo            string                   `json:"slug_seo" `
	Is_end              bool                     `json:"is_end" `
	Is_downloadable     int                      `json:"is_downloadable" `
	// Default_episode VODDataObjectStruct
	Default_episode struct {
		Id                  string                   `json:"id" `
		Group_id            string                   `json:"group_id" `
		Type                int                      `json:"type" `
		Title               string                   `json:"title" `
		Ads                 []AdsOutputStruct        `json:"ads" `
		Link_play           LinkPlayStruct           `json:"link_play" `
		Is_premium          int                      `json:"is_premium" `
		Is_new              int                      `json:"is_new" `
		Short_description   string                   `json:"short_description" `
		Long_description    string                   `json:"long_description" `
		Resolution          int                      `json:"resolution" `
		Runtime             int                      `json:"runtime" `
		Progress            int64                    `json:"progress" `
		Avg_rate            float64                  `json:"avg_rate" `
		Total_rate          int                      `json:"total_rate" `
		Is_watchlater       bool                     `json:"is_watchlater" `
		User_rating         int                      `json:"user_rating" `
		Episode             int                      `json:"episode" `
		Current_episode     string                   `json:"current_episode" `
		Content_provider_id string                   `json:"content_provider_id" `
		Seo                 SeoObjectStruct          `json:"seo" `
		Slug                string                   `json:"slug" `
		Slug_seo            string                   `json:"slug_seo" `
		Season_slug_seo     string                   `json:"season_slug_seo"`
		Permission          int                      `json:"permission" `
		Is_end              bool                     `json:"is_end" `
		No_seeker           bool                     `json:"no_seeker" `
		Images              ImagesOutputObjectStruct `json:"images" `
		Range_page_index    int64                    `json:"range_page_index" `
		Canonical_url       string                   `json:"canonical_url" `
		Direct_url          string                   `json:"direct_url" `
		Direct_content      string                   `json:"direct_content" `
		Share_url           string                   `json:"share_url" `
		Share_url_seo       string                   `json:"share_url_seo" `
		Seo_content         string                   `json:"seo_content" `
		Drm_service_name    string                   `json:"drm_service_name" `
		Intro_start         int64                    `json:"intro_start" `
		Intro_end           int64                    `json:"intro_end" `
		Outtro_start        int64                    `json:"outtro_start" `
		Outtro_end          int64                    `json:"outtro_end" `
	} `json:"default_episode" `
	Permission       int                        `json:"permission" `
	PackageGroup     []PackageGroupObjectStruct `json:"packages" `
	No_seeker        bool                       `json:"no_seeker" `
	Geo_check        int                        `json:"geo_check" `
	Range_page_index int64                      `json:"range_page_index" `
	Movie            struct {
		Title           string          `json:"title" `
		Episode         int             `json:"episode" `
		Current_episode string          `json:"current_episode" `
		Seo             SeoObjectStruct `json:"seo" `
		People          struct {
			Director []PeopleOutputStruct `json:"director" `
			Actor    []PeopleOutputStruct `json:"actor" `
		} `json:"people" `
		Tags []struct {
			Id   string          `json:"id" `
			Name string          `json:"name" `
			Type string          `json:"type" `
			Seo  SeoObjectStruct `json:"seo" `
			Slug string          `json:"slug" `
		} `json:"tags" `
	} `json:"movie" `
	Canonical_url     string                   `json:"canonical_url" `
	Direct_url        string                   `json:"direct_url" `
	Direct_content    string                   `json:"direct_content" `
	Share_url         string                   `json:"share_url" `
	Share_url_seo     string                   `json:"share_url_seo" `
	Seo_content       string                   `json:"seo_content" `
	Drm_service_name  string                   `json:"drm_service_name" `
	Have_trailer      int                      `json:"have_trailer" `
	Is_vip            int                      `json:"is_vip" `
	Intro_start       int64                    `json:"intro_start" `
	Intro_end         int64                    `json:"intro_end" `
	Outtro_start      int64                    `json:"outtro_start" `
	Outtro_end        int64                    `json:"outtro_end" `
	Related_season    []RelatedSeasonObjStruct `json:"related_season,omitempty" `
	Category          int                      `json:"category" `
	Vod_schedule      VodScheduleStruct        `json:"vod_schedule"`
	Season_name       string                   `json:"season_name"`
	Show_name         string                   `json:"show_name"`
	Min_rate          float64                  `json:"min_rate,omitempty"` //not display
	Created_at        int64                    `json:"created_at"`
	Rating_value      string                   `json:"rating_value"`
	Author_name       string                   `json:"author_name"`
	Review_body       string                   `json:"review_body"`
	Trailer_link_play struct {
		Hls_link_play  string `json:"hls_link_play"`
		Dash_link_play string `json:"dash_link_play"`
	} `json:"trailer_link_play"`
}

type PackageGroupObjectStruct struct {
	Id           int    `json:"id" `
	Name         string `json:"name" `
	Price        int    `json:"price" `
	Period       string `json:"period" `
	Period_value int    `json:"period_value" `
}

type ContentDetailPersonalObjStruct struct {
	Id         string         `json:"id" `
	Group_id   string         `json:"group_id" `
	Permission int            `json:"permission" `
	Link_play  LinkPlayStruct `json:"link_play" `
	Thumbs     struct {
		Vtt      string `json:"vtt" `
		Vtt_lite string `json:"vtt_lite" `
		Image    string `json:"image" `
	} `json:"thumbs" `
	Ads              []AdsOutputStruct          `json:"ads" `
	Subtitles        []SubtitlesOuputStruct     `json:"subtitles" `
	Audios           []AudiosOutputStruct       `json:"audios" `
	Progress         int64                      `json:"progress" `
	Is_watchlater    bool                       `json:"is_watchlater" `
	User_rating      int                        `json:"user_rating" `
	Drm_service_name string                     `json:"drm_service_name" `
	PackageGroup     []PackageGroupObjectStruct `json:"packages" `
	Views            int64                      `json:"views"`
	DownloadProfile  []ContentProfileObjStruct  `json:"download_profile" `
	Intro            struct {
		Start int64 `json:"start" `
		End   int64 `json:"end" `
	} `json:"intro" `
	Outtro struct {
		Start int64 `json:"start" `
		End   int64 `json:"end" `
	} `json:"outtro" `
	Is_vip               int     `json:"is_vip" `
	UserInteractionCount float64 `json:"user_interaction_count" `
}
type ContentProfileObjStruct struct {
	Name       string `json:"name" `
	Resolution string `json:"resolution" `
	Size       int64  `json:"size" `
}

type ItemsEpisodeObjOutputStruct struct {
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
	Type            int                      `json:"type" `
	Title           string                   `json:"title" `
	Resolution      int                      `json:"resolution" `
	Runtime         int                      `json:"runtime" `
	Images          ImagesOutputObjectStruct `json:"images" `
	Progress        int64                    `json:"progress" `
	Is_watchlater   bool                     `json:"is_watchlater" `
	Is_premium      int                      `json:"is_premium" `
	Is_new          int                      `json:"is_new" `
	Is_trailer      int                      `json:"is_trailer" `
	Episode         int                      `json:"episode" `
	Current_episode string                   `json:"current_episode" `
	// Intro_start     int64                    `json:"intro_start" `
	// Intro_end       int64                    `json:"intro_end" `
	// Outtro_start    int64                    `json:"outtro_start" `
	// Outtro_end      int64                    `json:"outtro_end" `
	Seo          SeoObjectStruct `json:"seo" `
	Slug         string          `json:"slug" `
	Slug_seo     string          `json:"slug_seo" `
	Display_name string          `json:"display_name"`
}

type EpisodeObjOutputStruct struct {
	Items []ItemsEpisodeObjOutputStruct `json:"items" `

	Tracking_data struct {
		Recommendation_id string `json:"recommendation_id" `
		Type              string `json:"type" `
	} `json:"tracking_data" `

	Metadata struct {
		User_id string `json:"user_id" `
		Total   int    `json:"total" `
		Limit   int    `json:"limit" `
		Page    int    `json:"page" `
	} `json:"metadata" `
	Seo SeoObjectStruct `json:"seo" `
}

type RelatedObjOutputStruct struct {
	Items    []ItemOutputObjStruct `json:"items" `
	Metadata struct {
		Total int `json:"total" `
		Limit int `json:"limit" `
		Page  int `json:"page" `
	} `json:"metadata" `
	Seo SeoObjectStruct `json:"seo" `
}

type ItemOutputObjStruct struct {
	Id              string                   `json:"id" `
	Group_id        string                   `json:"group_id" `
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
	// Slug      string `json:"slug" `
	// Slug_seo  string `json:"slug_seo" `
	Geo_check int `json:"geo_check" `
	People    struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags_display []string `json:"tags_display" `
	Ranking      int      `json:"ranking" `
}

type RelatedVideosObjOutputStruct struct {
	Items    []ItemOutputObjStruct `json:"items" `
	Metadata struct {
		Total int `json:"total" `
		Limit int `json:"limit" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type SamsungSmartHubPreview struct {
	Expires  int64               `json:"expires" `
	Sections []SamsungSHPSection `json:"sections" `
}
type SamsungSHPSection struct {
	Title string           `json:"title"`
	Tiles []SamsungSHPItem `json:"tiles"`
}
type SamsungSHPItem struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Image_ratio string `json:"image_ratio"`
	Image_url   string `json:"image_url"`
	Action_data string `json:"action_data"`
	Is_playable bool   `json:"is_playable"`
	Position    int    `json:"position"`
}

// DT-12265
type ViewObjectStruct struct {
	ContentId string `json:"contentId" `
	View      int64  `json:"view" `
}
