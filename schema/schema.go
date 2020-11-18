package schema

import "gopkg.in/mgo.v2/bson"

// "gopkg.in/mgo.v2/bson"

type TrackingDataStruct struct {
	Recommendation_id string `json:"recommendation_id" `
	Type              string `json:"type" `
}

type RibbonDataObjectStruct struct {
	Id           string
	Slug         string
	Name         string
	Type         int
	Ribbon_items []VODDataObjectStruct
}

type PageRibbonsObjectStruct struct {
	Id                 string
	Name               string
	Slug               string
	Type               int
	Status             int
	Display_image_type int
	Menus              []struct {
		Id        string `json:"id" `
		Slug      string `json:"slug" `
		Odr       int    `json:"odr" `
		Platforms []int  `json:"platforms" `
	} `json:"menus" `
	Tags []struct {
		Id   string
		Name string
		Type string
		Slug string
		Seo  SeoObjectStruct
	} `json:"tags" `
	Odr          int                   `json:"odr" `
	Platforms    []int                 `json:"platforms" `
	Ribbon_items []VODDataObjectStruct `json:"items" `
	Seo          SeoObjectStruct
}

type VODDataObjectStruct struct {
	Id                string
	Rib_item_id       string
	Rib_ids           []string
	Group_id          string
	R_id              string
	Type              int
	Is_premium        int
	Is_new            int
	Is_trailer        int
	Title             string
	Known_as          string
	Type_name         string
	Ads               []AdsOutputStruct
	Is_new_update     bool
	Short_description string
	Long_description  string
	Caption_langs     string
	Caption           string
	Resolution        int
	Runtime           int
	Released          bool
	Is_downloadable   int
	Season            []struct {
		Episode_name      string
		Long_description  string
		Link_play         string
		Title             string
		Current_episode   string
		Id                string
		Short_description string
		Type              int
		Season            int
		Share_urls        struct{}
		Known_as          string
		Released          bool
		Slug              string
		Isdvr             int
		Vod_schedule      struct {
			Day     string
			Channel string
			Hour    string
		}
		Group_id     string
		Episode      int
		Release_year int
		Release_date string
		Season_nane  string
		Publish_date string
		Runtime      int
		Resolution   int
	}
	People struct {
		Director []PeopleOutputStruct `json:"director" `
		Actor    []PeopleOutputStruct `json:"actor" `
	} `json:"people" `
	Tags []struct {
		Id   string
		Name string
		Type string
		Seo  SeoObjectStruct
		Slug string
	}
	Subtitles []struct {
		Id         string
		Code_name  string
		Is_default int64
		Uri        string
		Title      string
	}
	Audios []struct {
		Id         string
		Code_name  string
		Is_default int64
		Index      string
		Title      string
	}
	Player_logo    string
	Player_logo_4k string
	Release_year   int
	Publish_date   string
	Release_date   string
	Link_play      struct {
		Dash_asset_id       string
		Hls_asset_id        string
		Dash_link_play      string
		Hls_link_play       string
		H265_dash_link_play string
		H265_hls_link_play  string
		Map_profile         struct {
			Full_hd struct {
				Width     int
				Bandwidth int
				Name      string
				Height    int
			}
			Hd struct {
				Width     int
				Bandwidth int
				Name      string
				Height    int
			}
			Sd struct {
				Width     int
				Bandwidth int
				Name      string
				Height    int
			}
		}
	}
	Trailer struct {
		Dash_link_play string
		Hls_link_play  string
	}
	Link_live string
	Metadata  struct {
		Genre    []string
		Country  []string
		Subtitle []string
	}
	Image_soucre    []ImageSoucreObjectStruct `json:"image_soucre" `
	Images          ImageAllPlatformStruct
	Progress        int64
	Vod_schedule    VodScheduleStruct
	Is_notify       bool
	Nofity_sub      bool
	Avg_rate        float64
	Odr             float64
	Views           int
	Total_rate      int
	Is_watchlater   bool
	User_rating     int
	Episode         int
	Current_episode string
	Geo_check       int
	Platforms       []int
	Movie           struct {
		Title string
	}
	Season_slug_seo      string
	Content_provider_id  string
	Custom_ads           string
	Enable_ads           int
	Intro_start          int64
	Intro_end            int64
	Outtro_start         int64
	Outtro_end           int64
	Seo                  SeoObjectStruct
	Slug                 string
	Slug_seo             string
	Slug_seo_v5          string
	Share_url            string
	Share_url_seo        string
	More_info            interface{}
	Status               int
	Start_time           int64
	End_time             int64
	Subtitle             string
	No_seeker            bool
	Canonical_url        string
	Direct_url           string
	Direct_content       string
	Seo_content          string
	Drm_service_name     string
	Label_subtitle_audio string
	Label_public_day     string
	Is_vip               int
	Category             int
	Is_coming_soon       int
	Created_at           int64
	Season_name          string
	Show_name            string
	Display_name         string
	Drm_vieon_key        string
	Min_rate             float64
	Rating_value         string
	Author_name          string
	Review_body          string
	Trailer_link_play    struct {
		Hls_link_play  string
		Dash_link_play string
	}
	Tags_display []string
	Ranking      int
}

type ImageAllPlatformStruct struct {
	Web struct {
		Home_vod_hot      string
		Vod_thumb_big     string
		Banner            string
		Home_carousel_tv  string
		Home_carousel_web string
		Thumbnail         string
		Vod_thumb         string
		Poster            string
		Big_thumb         string

		Thumbnail_hot_v4 string
		Thumbnail_big_v4 string
		Carousel_tv_v4   string
		Carousel_app_v4  string
		Carousel_web_v4  string
		Thumbnail_v4     string
		Poster_v4        string
		Promotion_banner string
		Title_card_light string
		Title_card_dark  string
		Poster_original  string
		Thumb_original   string
	}
	App struct {
		Home_vod_hot      string
		Vod_thumb_big     string
		Banner            string
		Home_carousel_tv  string
		Home_carousel_web string
		Thumbnail         string
		Vod_thumb         string
		Poster            string
		Big_thumb         string

		Thumbnail_hot_v4 string
		Thumbnail_big_v4 string
		Carousel_tv_v4   string
		Carousel_app_v4  string
		Carousel_web_v4  string
		Thumbnail_v4     string
		Poster_v4        string
		Promotion_banner string
		Title_card_light string
		Title_card_dark  string
		Poster_original  string
		Thumb_original   string
	}

	Smarttv struct {
		Home_vod_hot      string
		Vod_thumb_big     string
		Banner            string
		Home_carousel_tv  string
		Home_carousel_web string
		Thumbnail         string
		Vod_thumb         string
		Poster            string
		Big_thumb         string

		Thumbnail_hot_v4 string
		Thumbnail_big_v4 string
		Carousel_tv_v4   string
		Carousel_app_v4  string
		Carousel_web_v4  string
		Thumbnail_v4     string
		Poster_v4        string
		Promotion_banner string
		Title_card_light string
		Title_card_dark  string
		Poster_original  string
		Thumb_original   string
	}

	Tablet struct {
		Home_vod_hot      string
		Vod_thumb_big     string
		Banner            string
		Home_carousel_tv  string
		Home_carousel_web string
		Thumbnail         string
		Vod_thumb         string
		Poster            string
		Big_thumb         string

		Thumbnail_hot_v4 string
		Thumbnail_big_v4 string
		Carousel_tv_v4   string
		Carousel_app_v4  string
		Carousel_web_v4  string
		Thumbnail_v4     string
		Poster_v4        string
		Promotion_banner string
		Title_card_light string
		Title_card_dark  string
		Poster_original  string
		Thumb_original   string
	}

	Mobile struct {
		Home_vod_hot      string
		Vod_thumb_big     string
		Banner            string
		Home_carousel_tv  string
		Home_carousel_web string
		Thumbnail         string
		Vod_thumb         string
		Poster            string
		Big_thumb         string

		Thumbnail_hot_v4 string
		Thumbnail_big_v4 string
		Carousel_tv_v4   string
		Carousel_app_v4  string
		Carousel_web_v4  string
		Thumbnail_v4     string
		Poster_v4        string
		Promotion_banner string
		Title_card_light string
		Title_card_dark  string
		Poster_original  string
		Thumb_original   string
	}
}

type MenuItemObjStruct struct {
	Id        string
	Name      string
	Slug      string
	Parent_id string
	Required  bool
	Icon      string
	Icon_text string
	Tag       struct {
		Id   string
		Name string
		Type string
		Slug string
		Seo  SeoObjectStruct
	}
	List_ribbon     string
	Title_ribbon    string
	Quantity_ribbon int
	Sub_menu_ribbon []MenuRibbonObjectStruct
	Layout_types    string
	Seo             SeoObjectStruct
	Have_banner     int
}

type MenuRibbonObjectStruct struct {
	Id          string          `json:"id"`
	Name_filter string          `json:"name"`
	Slug_filter string          `json:"slug"`
	Seo         SeoObjectStruct `json:"seo"`
}

type VersionDataObjStruct struct {
	Version_data_smarttv string
}

type LiveTVVTVCabObjStruct struct {
	Services      []LiveTVVTVCabObjDetailChannelStruct `json:"services" `
	Version       string                               `json:"version" `
	Total_records int                                  `json:"total_records" `
}

type LiveTVVTVCabObjDetailChannelStruct struct {
	Technical struct {
		Id                  string `json:"id" `
		MainChannelId       string `json:"mainChannelId" `
		DrmId               string `json:"drmId" `
		Title               string
		ShortName           string `json:"shortName" `
		LongName            string `json:"longName" `
		Active              bool   `json:"active" `
		IsFavorite          bool   `json:"is_favorite" `
		PrivateMetadata     string
		PrivateMetadataInfo struct {
			Thumb string `json:"thumb" `
		}
		Categories       []string
		NetworkLocation  string
		PromoImages      []string
		StartOverSupport bool                           `json:"startOverSupport" `
		CatchUpSupport   bool                           `json:"catchUpSupport" `
		Programme        LiveTVVTVCabObjDetailEPGStruct `json:"programme" `
		Is_ready         bool                           `json:"is_ready" `
	} `json:"technical" `
	Programme    LiveTVVTVCabObjDetailEPGStruct `json:"programme" `
	Seo          SeoObjectStruct                `json:"seo" `
	Permission   int                            `json:"permission" `
	PackageGroup []PackageGroupObjectStruct     `json:"packages" `
	Socket_info  struct {
		Tcp_port string `json:"tcp_port" `
		Host     string `json:"host" `
		Port     string `json:"port" `
	} `json:"socket_info" `
}

type LiveTVVTVCabEPGObjStruct struct {
	Programmes    []LiveTVVTVCabObjDetailEPGStruct `json:"programmes" `
	Version       string                           `json:"version" `
	Total_records int                              `json:"total_records" `
}

type LiveTVVTVCabObjDetailEPGStruct struct {
	Period struct {
		Duration int64   `json:"duration" `
		Start    float64 `json:"start" `
		End      float64 `json:"end" `
		Provider string  `json:"provider" `
	} `json:"period" `
	Title       string `json:"title" `
	IsStartOver bool   `json:"isStartOver" `
	IsCatchUp   bool   `json:"isCatchUp" `
	Description string `json:"description" `
	Episode     string `json:"episode" `
	Synopsis    string `json:"synopsis" `
	PromoImages []string
	Id          string          `json:"id" `
	Seo         SeoObjectStruct `json:"seo" `
}

type UscUserWatchlaterContentObjStruct struct {
	User_id      string
	Content_id   string
	Entity_type  int
	Updated_date int64
	Status       int
}

type UscUserChannelFavoriteContentObjStruct struct {
	User_id      string
	Content_id   string
	Entity_type  int
	Updated_date int64
	Status       int
}

type TagObjectStruct struct {
	Id   string          `json:"id"`
	Name string          `json:"name"`
	Type string          `json:"type"`
	Seo  SeoObjectStruct `json:"seo"`
	Slug string          `json:"slug"`
}

type ArtistObjectStruct struct {
	Id       string
	Name     string
	Slug     string
	Gender   int
	Birthday string
	Status   int
	Info     string
	Country  struct {
		Id   string
		Name string
	}
	Job    string
	Seo    SeoObjectStruct `json:"seo" `
	Images struct {
		Avatar string `json:"avatar" `
	} `json:"images" `
}

type ImageArttistObjectStruct struct {
	Id              string
	Image_name      string
	Image_type      string
	Resolution_type int
	Url             string
}

type MessageObjectStruct struct {
	Id          string      `json:"id" `
	Title       string      `json:"title" `
	Content     string      `json:"content" `
	Expire      int64       `json:"expire" `
	Status      int         `json:"status" `
	Image       string      `json:"image" `
	More_info   interface{} `json:"more_info" `
	Entity_type int         `json:"entity_type" `
	Entity_id   string      `json:"entity_id" `
	Episode_id  string      `json:"episode_id" `
	Created_at  int64       `json:"created_at" `
	Is_push_all int         `json:"is_push_all" `
	User_ids    []string    `json:"user_ids" `
}

type LiveEventObjectStruct struct {
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
	Images      ImageAllPlatformStruct
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
	Start_time      string `json:"start_time" `
	Str_to_time     int64
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
type NexteventInfoStruct struct {
	Title       string  `json:"title" `
	Location    string  `json:"location" `
	Start_time  string  `json:"start_time" `
	Total_award float64 `json:"total_award" `
}
type PageRibbonsV3ObjectStruct struct {
	Id                 string `json:"id" `
	Name               string `json:"name" `
	Type               int    `json:"type" `
	Display_image_type int    `json:"display_image_type" `
	Menus              []struct {
		Id string `json:"id" `
	} `json:"menus" `
	Tags struct {
		Id   string          `json:"id" `
		Name string          `json:"name" `
		Type string          `json:"type" `
		Slug string          `json:"slug" `
		Seo  SeoObjectStruct `json:"seo" `
	} `json:"tags" `
	Ribbon_items  []RibbonItemV3ObjectStruct `json:"items" `
	Tracking_data TrackingDataStruct         `json:"tracking_data" `
}

type QuestionStatisticsStruct struct {
	Id         bson.ObjectId            `bson:"_id,omitempty" json:"id"`
	Channel_id string                   `json:"channel_id"`
	Answers    []AnswerStatisticsStruct `json:"answers"`
}
type AnswerStatisticsStruct struct {
	Answer_id     string  `json:"answer_id"`
	Answer_title  string  `json:"answer_title"`
	Total_users   int64   `json:"total_users"`
	Percent_users float64 `json:"percent_users"`
}

type ReportTypeStruct struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UserReportContentStruct struct {
	Id            string `json:"id"`
	Message       string `json:"message"`
	Created_at    string `json:"Created_at"`
	Updated_at    string `json:"updated_at"`
	Status        string `json:"status"`
	Entity_id     string `json:"entity_id"`
	Report_id     string `json:"report_id"`
	User_id       string `json:"user_id"`
	User_agent    string `json:"user_agent"`
	Access_token  string `json:"access_token"`
	Video_profile string `json:"video_profile"`
	Audio_profile string `json:"audio_profile"`
	Subtitle      string `json:"subtitle"`
}
type CommentObjectStruct struct {
	Id         string                `json:"id"`
	Content_id string                `json:"content_id,omitempty"`
	Message    string                `json:"message,omitempty"`
	Status     int                   `json:"status,omitempty"`
	Created_at int64                 `json:"created_at,omitempty" `
	Updated_at int64                 `json:"updated_at,omitempty" `
	Reply      []CommentObjectStruct `json:"reply,omitempty"`
	User       UserCommentStruct     `json:"user,omitempty"`
	Platforms  PlatformOutputStruct  `json:"platforms" `
	Oldmessage string                `json:"oldmessage,omitempty"`
	Pin        int                   `json:"pin"`
}

type UserCommentStruct struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int    `json:"gender"`
}

// Thu DT-9922

type BadwordObjectStruct struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Content    string        `json:"content,omitempty"`
	Key        string        `json:"key,omitempty"`
	Status     int           `json:"status,omitempty"`
	Type       int           `json:"type,omitempty"`
	Created_at string        `json:"created_at,omitempty" `
	Updated_at string        `json:"updated_at,omitempty" `
}

type StructAdvWord struct {
	Content string
	Type    int
}

type RatingDetailContentStruct struct {
	Avg_rate    float64
	Total_rate  int
	Total_point float64
}
