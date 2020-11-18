package schema

type VodOptimizeOutputStruct struct {
	Id                string `json:"id" `
	Group_id          string `json:"group_id" `
	Type              int    `json:"type" `
	Title             string `json:"title" `
	Short_description string `json:"short_description" `
	Long_description  string `json:"long_description" `
	Resolution        int    `json:"resolution" `
	Runtime           int    `json:"runtime" `
	Is_premium        int    `json:"is_premium" `
	Is_new            int    `json:"is_new" `
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
	Player_logo      string                     `json:"player_logo" `
	Player_logo_4k   string                     `json:"player_logo_4k" `
	Images           ImagesOutputObjectStruct   `json:"images" `
	Avg_rate         float64                    `json:"avg_rate" `
	Total_rate       int                        `json:"total_rate" `
	Episode          int                        `json:"episode" `
	Current_episode  string                     `json:"current_episode" `
	Seo              SeoObjectStruct            `json:"seo" `
	Is_end           bool                       `json:"is_end" `
	No_seeker        bool                       `json:"no_seeker" `
	Geo_check        int                        `json:"geo_check" `
	Range_page_index int64                      `json:"range_page_index" `
	Canonical_url    string                     `json:"canonical_url" `
	Direct_url       string                     `json:"direct_url" `
	Direct_content   string                     `json:"direct_content" `
	Seo_content      string                     `json:"seo_content" `
	Drm_service_name string                     `json:"drm_service_name" `
	Have_trailer     int                        `json:"have_trailer" `
	Category         int                        `json:"category" `
	Vod_schedule     VodScheduleStruct          `json:"vod_schedule"`
	PackageGroup     []PackageGroupObjectStruct `json:"packages" `
	Related_season   []RelatedSeasonObjStruct   `json:"related_season,omitempty" `
	Episode_info     interface{}                `json:"episode_info" `
}

type VodScheduleStruct struct {
	Day     string `json:"day" `
	Channel string `json:"channel" `
	Hour    string `json:"hour" `
}

type RelatedSeasonObjStruct struct {
	Id      string `json:"id" `
	Seo_url string `json:"seo_url" `
	Title   string `json:"title" `
}
type EpisodeOptimizeOutputStruct struct {
	Id                string                   `json:"id" `
	Group_id          string                   `json:"group_id" `
	Type              int                      `json:"type" `
	Title             string                   `json:"title" `
	Ads               []AdsOutputStruct        `json:"ads,omitempty"`
	Short_description string                   `json:"short_description" `
	Long_description  string                   `json:"long_description" `
	Resolution        int                      `json:"resolution" `
	Runtime           int                      `json:"runtime" `
	Is_premium        int                      `json:"is_premium" `
	Is_new            int                      `json:"is_new" `
	Player_logo       string                   `json:"player_logo" `
	Player_logo_4k    string                   `json:"player_logo_4k" `
	Images            ImagesOutputObjectStruct `json:"images" `
	Episode           int                      `json:"episode" `
	Current_episode   string                   `json:"current_episode" `
	Seo               SeoObjectStruct          `json:"seo" `
	Is_end            bool                     `json:"is_end" `
	No_seeker         bool                     `json:"no_seeker" `
	Geo_check         int                      `json:"geo_check" `
	Range_page_index  int64                    `json:"range_page_index" `
	Canonical_url     string                   `json:"canonical_url" `
	Direct_url        string                   `json:"direct_url" `
	Direct_content    string                   `json:"direct_content" `
	Seo_content       string                   `json:"seo_content" `
	Drm_service_name  string                   `json:"drm_service_name" `
	Is_watchlater     bool                     `json:"is_watchlater" `
}

type VodDetailOptimizePersonalStruct struct {
	Id               string                 `json:"id" `
	Group_id         string                 `json:"group_id" `
	Permission       int                    `json:"permission" `
	Link_play        LinkPlayStruct         `json:"link_play" `
	Ads              []AdsOutputStruct      `json:"ads" `
	Subtitles        []SubtitlesOuputStruct `json:"subtitles" `
	Audios           []AudiosOutputStruct   `json:"audios" `
	Progress         int64                  `json:"progress" `
	Is_watchlater    bool                   `json:"is_watchlater" `
	User_rating      int                    `json:"user_rating" `
	Drm_service_name string                 `json:"drm_service_name" `
	Intro            struct {
		Start int64 `json:"start" `
		End   int64 `json:"end" `
	} `json:"intro" `
	Outtro struct {
		Start int64 `json:"start" `
		End   int64 `json:"end" `
	} `json:"outtro" `
	Episode_info interface{} `json:"episode_info" `
}
