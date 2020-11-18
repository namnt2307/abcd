package schema

const (
	VERSION_KEY = "v3_"

	COLLECTION_LIVE_TV       = "livetv"
	COLLECTION_LIVE_TV_GROUP = "livetv_group"

	PREFIX_REDIS_LIVE_TV_GROUP               = VERSION_KEY + "live_tv_group"
	PREFIX_REDIS_HASH_LIVE_TV                = VERSION_KEY + "hash_livetv"
	PREFIX_REDIS_HASH_LIVE_TV_SLUG           = VERSION_KEY + "hash_livetv_slug"
	PREFIX_REDIS_ZRANGE_USC_FAVORITE_V4      = VERSION_KEY + "usc_favorite_content_zrange_v4_"
	PREFIX_REDIS_ZRANGE_USC_WATCHED_V4       = VERSION_KEY + "usc_watched_content_zrange_v4_"
	PREFIX_REDIS_USC_FAVORITE_V4             = VERSION_KEY + "usc_favorite_content_v4_"
	PREFIX_REDIS_LIST_ID_LIVE_TV_ID_BY_GROUP = VERSION_KEY + "list_id_live_tv_id_by_group"

	KV_REDIS_LIVE_TV_BY_GROUP  = VERSION_KEY + "kv_live_tv_by_group"
	KV_REDIS_EPG_BY_LIVE_TV_ID = VERSION_KEY + "kv_epg_by_live_tv_id"
	KV_REDIS_DETAIL_LIVE_TV    = VERSION_KEY + "kv_detail_live_tv"
	KV_REDIS_DETAIL_EPG        = VERSION_KEY + "kv_detail_epg"
	KV_REDIS_LIVETV_DRM_KEY    = VERSION_KEY + "kv_livetv_drm_key"

	LOCAL_LIVE_TV_BY_GROUP = VERSION_KEY + "local_live_tv_by_group"
	LOCAL_LIVE_TV_GROUP    = VERSION_KEY + "local_live_tv_group"
	LOCAL_DETAIL_LIVE_TV   = VERSION_KEY + "local_detail_live_tv"
	LOCAL_LIVE_TV_EPG      = VERSION_KEY + "local_live_tv_epg"

	TOTAL_LIVE_TV_BY_GROUP                = VERSION_KEY + "total_live_tv_by_group"
	STATUS_LIVETV_GROUP_FOR_SUPER_PREMIUM = 99
	STATUS_LIVETV_GROUP_PUBLIC            = 1
)

// type PackageGroupObjectStruct struct {
// 	Id           int    `json:"id" `
// 	Name         string `json:"name" `
// 	Price        int    `json:"price" `
// 	Period       string `json:"period" `
// 	Period_value int    `json:"period_value" `
// }

// type SeoObjectStruct struct {
// 	Url         string `json:"url" `
// 	Description string `json:"description" `
// 	Title       string `json:"title" `
// }

type LiveTVGroupOutputObject struct {
	Id     string `json:"id" `
	Title  string `json:"title" `
	Status int    `json:"status" `
}

type LiveTVGroupStruct struct {
	Id  string
	Odr int
}

type LiveTVMongoObjectStruct struct {
	Id               string
	Livetv_group     []LiveTVGroupStruct
	Group_id         string
	Vtvcab_id        string
	Vtvcab_drm_id    string
	Main_channel_id  string
	Title            string
	Long_description string
	Slug             string
	Hls_link_play    string
	Dash_link_play   string
	Time_start       int
	Time_end         int
	Duration         int
	Type             int
	Status           int
	Is_catch_up      bool
	Max_time_catchup int
	Is_premium       int
	Geo_check        int
	Odr              int
	Seo              SeoObjectStruct
	Image_link       string
	Platforms        []int
	Drm_service_name string
	Drm_key          string
	Player_logo      string
}

type LiveTVObjectOutputStruct struct {
	Items    []LiveTVObjectStruct `json:"items" `
	Metadata struct {
		Total int `json:"total" `
		Limit int `json:"limit" `
		Page  int `json:"page" `
	} `json:"metadata" `
}

type DetailLiveTVObjectOutputStruct struct {
	Id               string                     `json:"id" `
	Livetv_group     []string                   `json:"categorys" `
	Vtvcab_drm_id    string                     `json:"vtvcab_drm_id" `
	Main_channel_id  string                     `json:"main_channel_id" `
	Title            string                     `json:"title" `
	Long_description string                     `json:"description" `
	Slug             string                     `json:"slug" `
	Link_play        string                     `json:"link_play" `
	Dash_link_play   string                     `json:"dash_link_play" `
	Hls_link_play    string                     `json:"hls_link_play" `
	Type             int                        `json:"type" `
	IsFavorite       bool                       `json:"is_favorite" `
	Is_catch_up      bool                       `json:"is_catch_up" `
	Seo              SeoObjectStruct            `json:"seo" `
	Image_link       string                     `json:"image_link" `
	Permission       int                        `json:"permission" `
	PackageGroup     []PackageGroupObjectStruct `json:"packages" `
	Programme        EpgObjectOutputStruct      `json:"programme" `
	Drm_service_name string                     `json:"drm_service_name" `
	Player_logo      string                     `json:"player_logo" `
	Is_premium       int                        `json:"is_premium" `
	Geo_check        int                        `json:"geo_check" `
	Usi              string                     `json:"usi" `
	Asset_id         string                     `json:"asset_id,omitempty" `
}

type LiveTVObjectStruct struct {
	Id string `json:"id" `
	// Livetv_group []string `json:"categorys" `
	// Vtvcab_drm_id    string                     `json:"vtvcab_drm_id" `
	// Main_channel_id  string                     `json:"main_channel_id" `
	Title            string          `json:"title" `
	Long_description string          `json:"description" `
	Slug             string          `json:"slug" `
	Type             int             `json:"type" `
	IsFavorite       bool            `json:"is_favorite" `
	Is_catch_up      bool            `json:"is_catch_up" `
	Is_premium       int             `json:"is_premium" `
	Geo_check        int             `json:"geo_check" `
	Seo              SeoObjectStruct `json:"seo" `
	Image_link       string          `json:"image_link" `
	// Permission       int             `json:"permission" `
	// PackageGroup     []PackageGroupObjectStruct `json:"packages" `
}

type EpgObjectOutputStruct struct {
	Id               string          `json:"id" `
	Title            string          `json:"title" `
	Long_description string          `json:"description" `
	Slug             string          `json:"slug" `
	Type             int             `json:"type" `
	Time_start       int             `json:"time_start" `
	Time_end         int             `json:"time_end" `
	Duration         int             `json:"duration" `
	Is_catch_up      bool            `json:"is_catch_up" `
	Hls_link_play    string          `json:"hls_link_play" `
	Dash_link_play   string          `json:"dash_link_play" `
	Seo              SeoObjectStruct `json:"seo" `
	Ott_disabled     int             `json:"Ott_disabled" ` //0 => cho phép (mặc định)| 1 => không cho phép
}

type UscUserLivetvFavoriteContentObjStruct struct {
	User_id      string
	Content_id   string
	Entity_type  int
	Updated_date int64
	Status       int
}
