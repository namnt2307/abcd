package schema

const (
	COLLECTION_RIB                    = "rib"
	COLLECTION_RIB_V3                 = "rib_v3"
	COLLECTION_RIB_ITEM_V3            = "rib_item_v3"
	COLLECTION_VOD                    = "vod"
	COLLECTION_NOTIFICATION           = "notification"
	COLLECTION_USC_USER_FAVORITE      = "usc_user_favorite_content"
	COLLECTION_USC_USER_WATCHLATER    = "usc_user_watchlater_content"
	COLLECTION_PLATFORM               = "platform"
	COLLECTION_SETTING                = "setting"
	COLLECTION_LIVE_EVENT             = "live_event"
	COLLECTION_MENU                   = "menu"
	COLLECTION_DM_QUESTIONS           = "dm_questions"
	COLLECTION_DM_RANKING             = "dm_ranking"
	COLLECTION_DM_USER_ANSWER         = "dm_user_answers"
	COLLECTION_COMMENT                = "comment"
	COLLECTION_BLACKLIST              = "blacklist"
	COLLECTION_TOPVIEWS_CONTENT_FINAL = "top_view_content_final"
	// DT-12265
	COLLECTION_VIEW = "vod_views"
	// DT-13431
	COLLECTION_RATE          = "vod_rate"
	TYPE_SEARCH_ITEM_SEASON  = 3
	TYPE_SEARCH_ITEM_EPISODE = 4
	TYPE_SEARCH_ITEM_ARTIST  = -2

	RIB_TYPE_MASTER_BANNER = 0
	RIB_TYPE_POSTER        = 1
	RIB_TYPE_THUMBNAIL     = 2
	RIB_TYPE_BANNER        = 3
	RIB_TYPE_LIVE_TV       = 4
	RIB_TYPE_LIVE_STREAM   = 5
	RIB_TYPE_LIVE_EPG      = 6

	RIB_ITEM_TYPE_NAME_LIVETV     = "LIVETV"
	RIB_ITEM_TYPE_NAME_EPG        = "EPG"
	RIB_ITEM_TYPE_NAME_LIVESTREAM = "LIVESTREAM"
	RIB_ITEM_TYPE_NAME_VOD        = "VOD"

	DOMAIN_API      = "API_DOMAIN"
	DOMAIN_ENDPOINT = "DOMAIN_ENDPOINT"

	//Type TAGS
	GENRE    = 1
	COUNTRY  = 2
	TAG      = 3
	CATEGORY = 4
	LIVETV   = 5
	AUDIO    = 6

	//Permission

	PERMISSION_VOD = 207

	//Type VOD
	VOD_TYPE_LIVESTREAM          = 0
	VOD_TYPE_MOVIE               = 1
	VOD_TYPE_SHOW                = 2
	VOD_TYPE_SEASON              = 3
	VOD_TYPE_EPISODE             = 4
	VOD_TYPE_LIVETV              = 5
	VOD_TYPE_TRAILER             = 6
	VOD_TYPE_EPG                 = 7
	VOD_TYPE_LIVE_EVENT_FINISHED = 10

	//Key cache Redis
	KV_PROVIDER_ADS_DETAIL            = "KV_PROVIDER_ADS_DETAIL_"
	PREFIX_REDIS_RIB                  = "rib_"
	PREFIX_REDIS_HASH_CONTENT         = "vod_hash"
	PREFIX_REDIS_HASH_CONTENT_SLUG    = "vod_by_slug_hash"
	PREFIX_REDIS_HASH_CONTENT_BY_TAGS = "vod_by_tags_hash"
	PREFIX_REDIS_HASH_RIB_ITEM_V3     = "rib_item_hash"
	PREFIX_REDIS_CONTENT_RATING       = "vod_rating_info_"

	PREFIX_REDIS_USC_FAVORITE        = "usc_favorite_content_"
	PREFIX_REDIS_ZRANGE_USC_FAVORITE = "usc_favorite_content_zrange_"
	PREFIX_REDIS_ZRANGE_USC_WATCHED  = "usc_watched_content_zrange_"

	PREFIX_REDIS_USC_RATING            = "usc_rating_content_"
	PREFIX_REDIS_USC_WATCHLATER        = "usc_watchlater_content_"
	PREFIX_REDIS_ZRANGE_USC_WATCHLATER = "usc_watchlater_content_zrange_"

	PREFIX_REDIS_ZRANGE_USC_WATCHING       = "usc_watching_content_zrange_v2_"
	PREFIX_REDIS_HASH_USC_WATCHING         = "usc_watching_content_hash_v2_"
	PREFIX_REDIS_HASH_USC_WATCHING_SESSION = "usc_watching_session_hash_v2_"

	PREFIX_REDIS_DOWNLOAD_PROFILE_VOD = "vod_profile_download_"

	PREFIX_REDIS_USC_MESSAGE_STATUS        = "usc_message_status_"
	PREFIX_REDIS_ZRANGE_USC_HISTORY_SEARCH = "usc_history_search_zrange_"

	PREFIX_REDIS_TAGS_SLUG_HASH             = "tags_slug_hash"
	REFIX_REDIS_TAGS_SEO_HASH               = "tags_seo_hash"
	REFIX_REDIS_ARTIST_HASH                 = "artist_hash"
	PREFIX_REDIS_ARTIST                     = "artist"
	PREFIX_REDIS_SEARCH                     = "search"
	PREFIX_REDIS_SUGGEST                    = "search_suggest"
	PREFIX_REDIS_SEARCH_TOP_KEYWORD         = "search_top_keyword"
	PREFIX_REDIS_ZRANGE_LIVE_EVENT          = "live_event_zrange"
	PREFIX_REDIS_HASH_LIVE_EVENT            = "live_event_hash"
	PREFIX_REDIS_ZRANGE_LIVE_EVENT_FINISHED = "live_event_finished_zrange"
	PREFIX_REDIS_ZRANGE_RELATED_VIDEOS      = "zrange_related_videos"
	PREFIX_REDIS_LIVE_EVENT_FINISHED_TOTAL  = "live_event_finished_total"
	PREFIX_REDIS_USER_REPORT                = "user_report"
	PREFIX_REDIS_ZRANGE_COMMENT             = "comment_zrange"
	PREFIX_REDIS_HASH_COMMENT               = "comment_hash"

	PREFIX_REDIS_SETING = "seting"

	LIST_RIB_ITEM            = "list_rib_item"
	LIST_RIB_ITEM_V3         = "list_rib_item_v3"
	CONTENT_LIST_EPS         = "content_list_eps_"
	CONTENT_LIST_RELATED     = "content_list_related"
	CONTENT_LIST_BY_SLUG_EPS = "content_list_by_slug_eps"
	EPISODE_GROUP_ID_BY_SLUG = "episode_group_id_by_slug"
	EPISODE_VOD_TOTAL        = "episode_vod_total"
	RELATED_VIDEOS           = "related_videos"
	RELATED_VIDEOS_TOTAL     = "related_videos_total"
	RELATED_VOD              = "related_vod_v2"
	RELATED_VOD_TOTAL        = "related_vod_total_v2"
	LIST_RIBBONS_BANNER      = "list_ribbons_banner"
	LIST_RIBBONS             = "list_ribbons"
	LIST_RIBBONS_BANNER_V3   = "list_ribbons_banner_v3"
	LIST_RIBBONS_V3          = "list_ribbons_v3"
	LIST_MENU                = "list_menu"
	TAGS_VOD                 = "tags_vod"
	TAGS_VOD_TOTAL           = "tags_vod_total"
	ARTIST_IMAGES            = "artist_images"
	ARTIST_SLUG              = "artist_slug"
	ARTIST_VOD               = "artist_vod"
	ARTIST_VOD_ID            = "artist_vod_id"
	ARTIST_VOD_TOTAL         = "artist_vod_total"
	ARTIST_RELATED_TOTAL     = "artist_related_total"
	ARTIST_RELATED           = "artist_related"
	TAGS_VOD_ID              = "tags_vod_id"
	MESSAGE_LIST             = "message_list"
	ARTIST_IMAGES_SIZE       = "artist_images_size"

	// Key local cache

	LOCAL_CONTENT_PERMISSION        = "local_content_permission"
	LOCAL_CONTENT                   = "local_content"
	LOCAL_VOD_V3                    = "local_vod_v3"
	LOCAL_LIST_EPISODE              = "local_list_episode"
	LOCAL_LIST_EPISODE_RANGE        = "local_list_episode_range"
	LOCAL_RELATED                   = "local_related"
	LOCAL_DETAIL_RIBBONS_BANNERS    = "local_detail_ribbons_banner"
	LOCAL_DETAIL_RIBBONS            = "local_detail_ribbons"
	LOCAL_DETAIL_RIBBONS_BANNERS_V3 = "local_detail_ribbons_banner_v3"
	LOCAL_DETAIL_RIBBONS_V3         = "local_detail_ribbons_v3"
	LOCAL_DETAIL_RIBBON_INFO        = "local_detail_ribbon_info"
	LOCAL_PAGE_RIBBONS_V3           = "local_page_ribbons_v3"
	LOCAL_TAGS                      = "local_tags"
	LOCAL_TAGS_CATEGORY             = "local_tags_category"
	LOCAL_ARTIST                    = "local_artist"
	LOCAL_ARTIST_RELEATED           = "local_artist_releated"
	LOCAL_ARTIST_VOD                = "local_artist_vod"
	LOCAL_SEARCH                    = "local_search"
	LOCAL_LIVE_EVENT                = "local_live_event"
	LOCAL_LIVE_EVENT_ID             = "local_live_event_id"
	LOCAL_LIVE_EVENT_SLUG           = "local_live_event_slug"
	LOCAL_LIVE_EVENT_FINISHED       = "local_live_event_finished"
	LOCAL_REPORT_TYPE               = "local_report_type"
	LOCAL_LIST_COMMENT_IN_ENTITY    = "local_list_comment_in_entity"
	LOCAL_TOPVIEWS_VOD              = "local_topviews_vod"

	// Key KV Redis cache

	KV_REDIS_CONTENT                  = "KV_content"
	KV_REDIS_VOD_V3                   = "KV_vod_v3"
	KV_VOD_LIST_PACKAGE_GROUP         = "KV_vod_list_package_group"
	KV_REDIS_LIST_EPISODE             = "KV_list_episode"
	KV_RELATED                        = "KV_related_v2"
	KV_DETAIL_RIBBONS_BANNERS         = "KV_detail_ribbons_banners"
	KV_DETAIL_RIBBONS                 = "KV_detail_ribbons"
	KV_DETAIL_RIBBONS_BANNERS_V3      = "KV_detail_ribbons_banners_v3"
	KV_DETAIL_RIBBONS_V3              = "KV_detail_ribbons_v3"
	KV_DETAIL_RIBBON_INFO             = "KV_detail_ribbon_info"
	KV_DETAIL_RIBBON_INFO_V3          = "KV_detail_ribbon_info_v3"
	KV_ARTIST_RELATED                 = "KV_artist_related"
	KV_REDIS_LIVE_EVENT_ID            = "KV_live_event_id"
	KV_REDIS_LIVE_EVENT_SLUG          = "KV_live_event_slug"
	KV_REDIS_CATEGORY_TAGS            = "KV_category_tags"
	KV_REDIS_SLUG_EPG                 = "KV_slug_epg_%s_%s"
	KV_REDIS_LIVE_EVENT_FINISHED      = "KV_live_event_finished"
	KV_REDIS_LIVE_EVENT_FINISHED_ID   = "KV_live_event_finished_id"
	KV_REDIS_LIVE_EVENT_FINISHED_SLUG = "KV_live_event_finished_slug"
	KV_REDIS_REPORT_TYPE              = "KV_report_type"
	KV_USER_POST_COMMENT              = "KV_user_post_comment"
	KV_CHECK_DRM_VIEON                = "KV_check_drm_vieon"
	KV_DATA_TOPVIEWS_CONTENT          = "KV_data_topviews_content"
	//DT-12264
	KV_VIEWS                  = "VIEWS_"
	KV_LIST_ID_SEASON_BY_SHOW = "KV_list_id_season_by_show"

	PREFIX_CONFIG       = "config:"
	TTL_REDIS_LV1       = 1 * 24 * 3600
	TTL_REDIS_2_HOURS   = 2 * 3600
	TTL_REDIS_6_HOURS   = 6 * 3600
	TTL_REDIS_1_HOURS   = 3600
	TTL_REDIS_30_MINUTE = 1800
	TTL_REDIS_5_MINUTE  = 300

	// TTL_LOCALCACHE       = 30
	// TTL_KVCACHE          = 60
	TTL_KVCACHE_1_MINUTE = 5

	//Recommendation
	RECOMMENDATION_MIN_RANDOM_STRING = 12
	RECOMMENDATION_MAX_RANDOM_STRING = 20
	RECOMMENDATION_PARAMS_TYPE       = "type"
	RECOMMENDATION_PARAMS_ID         = "recommendation_id"
	RECOMMENDATION_CONTENT_LISTING   = "CONTENT_LISTING"
	TOPIC_REC_CLICK                  = "rec_click"
	RECOMMENDATION_VIEON             = "viePlay"
	RECOMMENDATION_GROUP             = "RECOMMENDATION"
	RECOMMENDATION_GROUP_SCENARIOS   = "RECOMMENDATION_SCENARIOS"
	RECOMMENDATION_NAME_YUSP         = "YUSP"
	RECOMMENDATION_NAME_VIEPLAY      = "VIEPLAY"
	REMOVE_FROM_PLAYLIST             = "REMOVE_FROM_PLAYLIST"
	ADD_TO_PLAYLIST                  = "ADD_TO_PLAYLIST"
	SCENARIO_BECAUSE_YOU_WATCHED     = "YMAL"

	// SETTING
	GROUP_DOMAIN_ENDPOINT_CONFIG            = "DOMAIN_ENDPOINT"
	GROUP_DOMAIN_ENDPOINT_API_DOMAIN_CONFIG = "API_DOMAIN"

	//YUSP
	GROUP_YUSP_CONFIG            = "YUSP_CONFIG"
	NAME_YUSP_CALL_API           = "YUSP_CALL_API"
	NAME_GET_ITEM_RECOMMENDATION = "GET_ITEM_RECOMMENDATION"
	CONTENT_TYPE                 = "contentType"
	AVAILABLE_EPISODE            = "availableEpisode"
	TOTAL_EPISODE                = "totalEpisode"
	SEASON                       = "season"
	SERIES                       = "series"
	VALUE                        = "value"
	REC_CLICK                    = "REC_CLICK"
	VIEW                         = "VIEW"
	WATCH                        = "WATCH"
	COMMENT                      = "COMMENT"

	// Permission
	GROUP_PERMISSION_CONFIG = "PERMISSION_CONFIG"

	// DeviceKey
	DEVICEKEY_WEB       = ""
	DEVICEKEY_MOBILEWEB = ""
	DEVICEKEY_ANDROID   = "AD_"
	DEVICEKEY_IOS       = "IOS_"
	DEVICEKEY_SMARTTV   = "STV_"

	//Search
	SEARCH                        = "SEARCH"
	SEARCHREGEXEPS                = "(tập|tap|ep|eps|episode|episodes|tập |tap |ep |eps |episode |episodes )([0-9])"
	SPECIAL_CHARACTERS            = `(~|!|{|}|'|"|\|/|,|:|;|@|#|$|%|[|]|tập|tap|ep|eps|episode|episodes|tập |tap |ep |eps |episode |episodes )([0-9])`
	SEARCH_RESULT                 = "SEARCH_RESULT"
	SEARCH_KEYWORD                = "SEARCH_KEYWORD"
	ZERO_RESULT                   = "ZERO_RESULT"
	ELASTICSEARCH                 = "ELASTIC"
	ELASTICSEARCH_SV2_INDEX       = "SEARCHING_V2_INDEX"
	ELASTICSEARCH_CONTENT_DOCTYPE = "CONTENT_DOCTYPE"
	ELASTICSEARCH_ARTIST_DOCTYPE  = "ARTIST_DOCTYPE"

	//Kafka
	GROUP_KAFKA_CONFIG                               = "KAFKA_CONFIG"
	GROUP_KAFKA_CONFIG_SERVER_LIST                   = "SERVER_LIST"
	GROUP_KAFKA_CONFIG_PREFIX_TOPIC                  = "PREFIX_TOPIC"
	GROUP_KAFKA_CONFIG_TOPIC_NAME_OTT_CACHING_QUEUES = "TOPIC_NAME_OTT_CACHING_QUEUES"
	NAME_TOPIC_REC_CLICK                             = "TOPIC_REC_CLICK"
	NAME_TOPIC_LOG_REQUEST                           = "TOPIC_REQUEST_LOG"

	// USC
	PROVIDER_MOBILE   = 0
	PROVIDER_EMAIL    = 1
	PROVIDER_GOOGLE   = 2
	PROVIDER_FACEBOOK = 3

	LIMIT_EPS_PER_RANGE = 30

	LIVEEVENT_LIMIT_RANKING = 10

	//SOCKET
	SOCKET_MESSAGE  = "SOCKET_MESSAGE"
	SOCKET_HOST     = "HOST"
	SOCKET_HOST_OLD = "HOST_OLD"
	SOCKET_PORT     = "PORT"
	SOCKET_TCP_PORT = "TCP_PORT"
	SOCKET_PROTOCOL = "PROTOCOL"

	//QUESTION
	QUESTION_STATUS_PUBLISH_WAIT      = 1
	QUESTION_STATUS_AUTO_PUBLISH_WAIT = 2
	QUESTION_STATUS_PUBLISH_DONE      = 3

	//PERMISSION CONTENT
	PERMISSION_REQUIRE_LOGIN   = 0
	PERMISSION_VALID           = 206
	PERMISSION_REQUIRE_PREMIUM = 207
	PERMISSION_REQUIRE_PACKAGE = 208
	PERMISSION_REQUEST_LIMIT   = 429
	PERMISSION_NOT_ALLOW       = 405

	//Page + Limit => max
	PAGE_MAX  = 50
	LIMIT_MAX = 100
	//STATUS USER
	USER_SUPER_PREMIUM = 99

	//
	CONFIG_KEY_PACKAGE_EXPIRED_DAYS = "package_expired_days"
)
