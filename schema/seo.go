package schema

const (
	URL                       = "url"
	TITLE                     = "title"
	DESCRIPTION               = "description"
	CATEGORY_TITLE            = "%s mới nhất %s, hay nhất %s - VieON"
	CATEGORY_DECRIPTION       = "Xem %s mới nhất %s, tuyển tập %s hay nhất %s"
	GEN_CON_TITLE             = "Phim %s mới nhất %s, phim %s hay nhất %s - VieON"
	GEN_CON_DESCRIPTION       = "Xem phim %s mới nhất %s, tuyển tập phim %s hay nhất %s - VieON"
	MULTIPLE_TAGS_TITLE       = "Tag %s - VieON"
	MULTIPLE_TAGS_DESCRIPTION = "Danh sách video thuộc Tag %s"
	SEO_TITLE                 = "Tìm Kiếm '%s' - VieON"
	SEO_TITLE_EPISODE         = "%s - Tập %s - VieON"
	SEO_DESCRIPTION_EPISODE   = "%s - tap-%s %s"
	CHANNEL_URL               = "/truyen-hinh-truc-tuyen/kenh-%s"
	CHANNEL_TITLE             = "Xem kênh %s - VieON"
	SEARCH_URL                = "/tim-kiem/?q=%s"
	SEARCH_TITLE              = "Tìm kiếm %s và các nội dung liên quan | VieON"
	SEARCH_DESCRIPTION        = "Kết quả tìm kiếm cho từ khóa %s, %s"
	ACTOR_URL                 = "/dien-vien/%s"
	ACTOR_TITLE               = "Diễn viên %s - VieON"
	ACTOR_DESCRIPTION         = "Thông tin diễn viên %s"
	DIRECTOR_URL              = "/dao-dien/%s"
	DIRECTOR_TITLE            = "Đạo diễn %s - VieON"
	DIRECTOR_DESCRIPTION      = "Thông tin đạo diễn %s"

	TAG_URL         = "/phim-hay/t/%s/"
	TAG_TITLE       = "Tuyển tập top %d %s | VieON"
	TAG_DESCRIPTION = "%s: %s"

	ARTIST_URL         = "/nghe-si/%s"
	ARTIST_TITLE       = "%s: Xem Tuyển tập nội dung chương trinh đặc sắc | VieON"
	ARTIST_DESCRIPTION = "Nghệ sĩ %s ở %s tham gia %d phim: %s ..."

	LIVE_EVENT_URL   = "/truc-tiep/%s"
	LIVE_EVENT_TITLE = "Trực tiếp %s - VieON"

	SEO_MENU_URL         = "/%s/"
	SEO_MENU_TITLE       = "Danh Sách %s Mới Nhất, Hay Nhất 2020 - VieON"
	SEO_MENU_DESCRIPTION = "Danh sách %s mới nhất, hay nhất năm 2020. Xem %s chất lượng cao Full HD Vietsub, Thuyết Minh, Lồng Tiếng"

	SEO_RIBBON_TYPE_LIVETV_URL     = "/truyen-hinh-truc-tuyen/"
	SEO_RIBBON_TYPE_LIVESTREAM_URL = "/live-streaming/"

	SEO_RIBBON_URL         = "/phim-hay/r/%s/"
	SEO_RIBBON_CHILD_URL   = "/%s-col/rib-%s"
	SEO_RIBBON_TITLE       = "Tuyển tập top %d %s | VieON"
	SEO_RIBBON_DESCRIPTION = "%s: %s"

	SEO_LIVETV_URL         = "/truyen-hinh-truc-tuyen/%s"
	SEO_LIVETV_EPG_URL     = "/truyen-hinh-truc-tuyen/%s/%s"
	SEO_LIVETV_TITLE       = "Xem truyền hình trực tuyến %s | VieON"
	SEO_LIVETV_DESCRIPTION = "Xem trực tiếp và xem lại các chương trình đã phát sóng của kênh truyền hình %s trên VieON với chất lượng Full HD, không giật lag, không quảng cáo."

	SEO_EPG_URL         = "/truyen-hinh-truc-tuyen/%s/%s"
	SEO_EPG_TITLE       = "Xem %s | VieON"
	SEO_EPG_DESCRIPTION = "Xem %s kênh truyền hình %s trên VieON với chất lượng Full HD, không giật lag, không quảng cáo. "

	SEO_VOD_URL = "/%s.html"
	SEO_VOD_TITLE = "%s | VieON"
	SEO_VOD_DESCRIPTION = "Xem %s của %s có sự tham gia của %s. Thuộc thể loại: %s"

	SEO_SEASION_TITLE = "%s - %d Tập | VieON"
	SEO_SEASION_DESCRIPTION = "Xem %s - %d Tập của %s có sự tham gia của %s. Thuộc thể loại: %s"


	SEO_VOD_EPISODE_TITLE = "%s %s - %d Tập | VieON"
	SEO_VOD_EPISODE_DESCRIPTION = "Xem %s %s - %d Tập của %s có sự tham gia của %s. Thuộc thể loại: %s"


	SEO_PREFIX_EPISODE = "--eps-"
	SEO_PREFIX_TRAILER = "--rel-"
	SEO_EPISODE_URL = "/%s"+SEO_PREFIX_EPISODE+"%s.html"
	SEO_TRAILER_URL = "/%s"+SEO_PREFIX_TRAILER+"%s.html"

	SEO_EPISODE_LIST_URL         = "%s/danh-sach-tap"
	SEO_EPISODE_LIST_TITLE       = "Danh Sách Tập %s - VieON"
	SEO_EPISODE_LIST_DESCRIPTION = "Danh Sách Tập %s : Trọn bộ, Chất lượng Full HD, Không giật lag."

	SEO_RELATED_URL         = "%s/de-xuat-cho-ban"
	SEO_RELATED_TITLE       = "Video Liên Quan - %s - VieON"
	SEO_RELATED_DESCRIPTION = "Danh sách video đề xuất mới nhất, hay nhất năm 2020 liên quan %s"
)

type SeoObjectStruct struct {
	Slug 		  string `json:"slug" `
	Url           string `json:"url" `
	Share_url     string `json:"share_url" `
	Title         string `json:"title" `
	Title_seo_tag string `json:"title_seo_tag" `
	Description   string `json:"description" `
	Canonical_tag string `json:"canonical_tag" `
	Meta_robots   string `json:"meta_robots" `
	Seo_text      string `json:"seo_text" `
	Alternate     string `json:"alternate" `
	Deeplink      string `json:"deeplink" `
}
