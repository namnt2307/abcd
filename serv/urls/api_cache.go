package urls

import (
	// . "../module"
	"cm-v5/serv/module/cache"
	"cm-v5/serv/module/subscription"

	"github.com/gin-gonic/gin"
)

func LogUpdateCacheRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sentry_log_request_cache(c)
		c.Next()
	}
}

func InitUrlsWarmUp(router *gin.Engine) {

	ver3 := router.Group("/backend/cm/v5/update-cache")
	{
		// ver3.GET("/update-cache/home", cache.UpdateCacheHome)
		// ver3.GET("/update-cache/detail", cache.UpdateCacheDetail)

		ver3.Use(LogUpdateCacheRequest())
		ver3.GET("/ribbon/:ribbon_id", cache.UpdateCacheRibbonV3)
		ver3.GET("/page_ribbons/:page_id", cache.UpdateCachePageRibbon)
		ver3.GET("/page_banners/:page_id", cache.UpdateCachePageBanner)

		//Khanh
		ver3.GET("/ribbon-by-menuid/:menu_id", cache.ClearCacheRibbonByMenuID)
		ver3.GET("/ribbon-by-ribbonid/:ribbon_id", cache.ClearCacheRibbonByRibbonId)
		ver3.GET("/ribbon-by-itemid/:ribbon_item_id", cache.ClearCacheRibbonByRibbonItemId)

		ver3.GET("/menu", cache.UpdateCacheMenu)
		ver3.GET("/content/:content_id", cache.UpdateCacheContentById)
		ver3.DELETE("/content/:content_id", cache.UpdateCacheContentById)
		ver3.GET("/episode/:group_id", cache.UpdateCacheContentEpisodeById)
		ver3.DELETE("/episode/:group_id", cache.UpdateCacheContentEpisodeById)
		ver3.GET("/artist/:people_id", cache.UpdateCacheArtistById)
		ver3.GET("/artist_related/:people_id", cache.UpdateCacheArtistRelatedById)
		ver3.GET("/artist_contents/:people_id", cache.UpdateCacheContentsByID)
		ver3.POST("/tags", cache.UpdateCacheTagListID)
		ver3.GET("/message", cache.UpdateListMessage)
		ver3.GET("/setting", cache.UpdateCacheSetting)
		ver3.GET("all/live-event", cache.UpdateCacheAllLiveEvent)
		ver3.GET("/live-event/:event_id", cache.UpdateCacheAllLiveEvent)
		ver3.GET("/packages/:content_id", cache.UpdateCachePackagesByContent)
		// ver3.POST("/slug/content", cache.GetContentBySlug)
		// ver3.GET("/episode/:content_id", cache.GetEpisodeById)
		// ver3.POST("/slug/episode", cache.GetEpisodeBySlug)
		// ver3.GET("/related/:content_id", cache.GetRelatedById)
		// ver3.POST("/slug/related", cache.GetRelatedBySlug)

		//KHANH clear cache V4 livetv
		ver3.GET("/all/livetv-group", cache.UpdateLiveTVGroup)
		ver3.GET("/livetv", cache.UpdateCacheLiveTV)
		ver3.GET("/livetv/epg", cache.UpdateCacheLiveTVEpg)
		ver3.GET("/search/seach_keywords", cache.UpdateCacheSearchKeyword)
		// DT-12264 add cache view
		ver3.GET("/add_view/:content_id/:views", cache.AddViewByContentId)

		//add cache subscription
		ver3.POST("subscription", subscription.ClearCacheSubscription)

		//cache config
		ver3.GET("config-key", cache.UpdateCacheConfigByKey)

		//cache content provider ads
		ver3.GET("/content_provider_ads/:provider_id", cache.UpdateAdsForCP)
	}
}
