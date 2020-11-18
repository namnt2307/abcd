package seo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "cm-v5/schema"
	. "cm-v5/serv/module"

	"github.com/gosimple/slug"
)

const (
	SEO_TYPE_VOD        = 0
	SEO_TYPE_LIVETV     = 1
	SEO_TYPE_TAG        = 2
	SEO_TYPE_COLLECTION = 3

	FORMAT_SLUG_TAG        = "/%s/%s-tag"
	FORMAT_SLUG_COLLECTION = "/%s/%s-rib"
)

func GetYear(current_only bool) string {
	currentTime := time.Now()
	year := currentTime.Year()

	if !current_only {
		return strconv.Itoa(year) + ", " + strconv.Itoa(year-1)

	}
	return strconv.Itoa(year)
}

func FormatSeoByTag(tagId, slug, name string, totalContent int, contentStr string, cache bool) SeoObjectStruct {
	var SeoObject SeoObjectStruct

	SeoObject.Url = fmt.Sprintf(TAG_URL, slug)
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = fmt.Sprintf(TAG_TITLE, totalContent, name)
	SeoObject.Title_seo_tag = SeoObject.Title
	SeoObject.Description = fmt.Sprintf(TAG_DESCRIPTION, name, contentStr)
	SeoObject.Seo_text = SeoObject.Description
	SeoObject.Canonical_tag = "https://vieon.vn/phim-hay/"
	SeoObject.Meta_robots = "index,follow"

	infoSEO, err := GetSEO(tagId, cache)
	if err != nil {
		return SeoObject
	}
	SeoObject.Url = infoSEO.Slug
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = infoSEO.Title
	SeoObject.Title_seo_tag = infoSEO.Title_seo_tag
	SeoObject.Description = infoSEO.Meta_description
	SeoObject.Seo_text = infoSEO.Seo_text
	SeoObject.Canonical_tag = infoSEO.Canonical_tag
	SeoObject.Meta_robots = infoSEO.Meta_robots
	SeoObject.Alternate = infoSEO.Alternate
	SeoObject.Deeplink = infoSEO.Deeplink

	return SeoObject
}

func FormatSeoByMutilTags(dataTags interface{}) SeoObjectStruct {
	var TagArr []TagObjectStruct

	dataByte, _ := json.Marshal(dataTags)
	json.Unmarshal(dataByte, &TagArr)

	var SeoObject SeoObjectStruct
	if len(TagArr) > 0 {

		SeoObject = TagArr[0].Seo
		return SeoObject
	}

	return SeoObject
}

func FormatSeoVOD(slug_seo_v5 string, seo SeoObjectStruct) SeoObjectStruct {
	seo.Url = fmt.Sprintf(SEO_VOD_URL, slug_seo_v5)
	seo.Share_url = HandleSeoShareUrlVieON(seo.Url)
	return seo
}

func FormatSeoVODDetail(Content ContentObjOutputStruct, cache bool) SeoObjectStruct {
	var seo SeoObjectStruct = Content.Seo
	var peopleStr, countryStr, categoryStr string
	var peopleArr []string

	for k, vPeople := range Content.People.Actor {
		if k >= 5 {
			break
		}
		peopleArr = append(peopleArr, vPeople.Name)
	}
	peopleStr = strings.Join(peopleArr, ", ")
	for _, vTag := range Content.Tags {
		if vTag.Type == "country" {
			countryStr = vTag.Name
		}

		if vTag.Type == "category" {
			categoryStr = vTag.Name
		}
	}

	if VOD_TYPE_EPISODE == Content.Type {
		seo.Title = fmt.Sprintf(SEO_VOD_EPISODE_TITLE, Content.Title, Content.Movie.Title, Content.Movie.Episode)
		seo.Title_seo_tag = seo.Title
		seo.Description = fmt.Sprintf(SEO_VOD_EPISODE_DESCRIPTION, Content.Title, Content.Movie.Title, Content.Movie.Episode, countryStr, peopleStr, categoryStr)
		seo.Seo_text = seo.Description
		seo.Canonical_tag = Content.Movie.Seo.Share_url
		seo.Meta_robots = "index,follow"
	} else if VOD_TYPE_SEASON == Content.Type {
		seo.Title = fmt.Sprintf(SEO_SEASION_TITLE, Content.Title, Content.Episode)
		seo.Title_seo_tag = seo.Title
		seo.Description = fmt.Sprintf(SEO_SEASION_DESCRIPTION, Content.Title, Content.Episode, countryStr, peopleStr, categoryStr)
		seo.Seo_text = seo.Description
		seo.Canonical_tag = seo.Share_url
		seo.Meta_robots = "index,follow"
	} else {
		seo.Title = fmt.Sprintf(SEO_VOD_TITLE, Content.Title)
		seo.Title_seo_tag = seo.Title
		seo.Description = fmt.Sprintf(SEO_VOD_DESCRIPTION, Content.Title, countryStr, peopleStr, categoryStr)
		seo.Seo_text = seo.Description
		seo.Canonical_tag = seo.Share_url
		seo.Meta_robots = "index,follow"
	}
	infoSEO, err := GetSEO(Content.Id, cache)
	if err == nil {
		seo.Title = infoSEO.Title
		seo.Title_seo_tag = infoSEO.Title_seo_tag
		seo.Description = infoSEO.Meta_description
		seo.Seo_text = infoSEO.Seo_text
		seo.Canonical_tag = infoSEO.Canonical_tag
		seo.Meta_robots = infoSEO.Meta_robots
		seo.Alternate = infoSEO.Alternate
		seo.Deeplink = infoSEO.Deeplink
	}
	return seo
}

func FormatSeoArtist(ArtistSeo ArtistObjectStruct, totalContent int, contentStr string) SeoObjectStruct {
	var SeoObject SeoObjectStruct
	SeoObject.Url = fmt.Sprintf(ARTIST_URL, ArtistSeo.Slug)
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = fmt.Sprintf(ARTIST_TITLE, ArtistSeo.Name)
	SeoObject.Title_seo_tag = fmt.Sprintf(ARTIST_TITLE, ArtistSeo.Name)
	SeoObject.Description = fmt.Sprintf(ARTIST_DESCRIPTION, ArtistSeo.Name, ArtistSeo.Country.Name, totalContent, contentStr)
	SeoObject.Canonical_tag = "https://vieon.vn/"
	SeoObject.Meta_robots = "noidex,follow"
	SeoObject.Seo_text = SeoObject.Description

	return SeoObject
}

func FormatSeoForSearch(keyword string, SearchResult SearchResultStruct) SeoObjectStruct {
	var SeoObject SeoObjectStruct
	if keyword != "" {
		SeoObject.Url = fmt.Sprintf(SEARCH_URL, slug.Make(keyword))
		var resultArr []string
		for k, v := range SearchResult.Items {
			if k >= 5 {
				break
			}
			resultArr = append(resultArr, v.Title)
		}

		var resultStr string = strings.Join(resultArr, ", ")
		SeoObject.Title = fmt.Sprintf(SEARCH_TITLE, keyword)
		SeoObject.Description = fmt.Sprintf(SEARCH_DESCRIPTION, keyword, resultStr)

		SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
		SeoObject.Title_seo_tag = SeoObject.Title
		SeoObject.Seo_text = SeoObject.Description
		SeoObject.Canonical_tag = "https://vieon.vn/"
		SeoObject.Meta_robots = "index,follow"
	}
	return SeoObject
}

func FormatSeoLiveEvent(slug, title, description string) SeoObjectStruct {
	var SeoObject SeoObjectStruct
	SeoObject.Url = fmt.Sprintf(LIVE_EVENT_URL, slug)
	SeoObject.Title = fmt.Sprintf(LIVE_EVENT_TITLE, title)
	SeoObject.Description = description
	return SeoObject
}

func FormatSeoMenu(menuId, slug, title string, cache bool) SeoObjectStruct {
	var titleUpFirst = strings.Title(title)
	var titleToLower = strings.ToLower(title)
	var SeoObject SeoObjectStruct

	SeoObject.Url = fmt.Sprintf(SEO_MENU_URL, slug)
	SeoObject.Title = fmt.Sprintf(SEO_MENU_TITLE, titleUpFirst)
	SeoObject.Description = fmt.Sprintf(SEO_MENU_DESCRIPTION, titleToLower, titleToLower)

	infoSEO, err := GetSEO(menuId, cache)
	if err != nil {
		return SeoObject
	}
	SeoObject.Url = infoSEO.Slug
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = infoSEO.Title
	SeoObject.Title_seo_tag = infoSEO.Title_seo_tag
	SeoObject.Description = infoSEO.Meta_description
	SeoObject.Seo_text = infoSEO.Seo_text
	SeoObject.Canonical_tag = infoSEO.Canonical_tag
	SeoObject.Meta_robots = infoSEO.Meta_robots
	SeoObject.Alternate = infoSEO.Alternate
	SeoObject.Deeplink = infoSEO.Deeplink

	return SeoObject
}

func FormatSeoRibbon(ribId, slug, title string, total int, contentStr string, cache bool) SeoObjectStruct {
	var SeoObject SeoObjectStruct

	SeoObject.Slug = slug
	SeoObject.Url = fmt.Sprintf(SEO_RIBBON_URL, slug)
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = fmt.Sprintf(SEO_RIBBON_TITLE, total, title)
	SeoObject.Title_seo_tag = SeoObject.Title
	SeoObject.Description = fmt.Sprintf(SEO_RIBBON_DESCRIPTION, title, contentStr)
	SeoObject.Seo_text = SeoObject.Description
	SeoObject.Canonical_tag = "https://vieon.vn/phim-hay/"
	SeoObject.Meta_robots = "noindex,follow"

	infoSEO, err := GetSEO(ribId, cache)
	if err != nil || infoSEO.Slug == "" {
		return SeoObject
	}
	SeoObject.Url = infoSEO.Slug
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = infoSEO.Title
	SeoObject.Title_seo_tag = infoSEO.Title_seo_tag
	SeoObject.Description = infoSEO.Meta_description
	SeoObject.Seo_text = infoSEO.Seo_text
	SeoObject.Canonical_tag = infoSEO.Canonical_tag
	SeoObject.Meta_robots = infoSEO.Meta_robots
	SeoObject.Alternate = infoSEO.Alternate
	SeoObject.Deeplink = infoSEO.Deeplink

	return SeoObject
}

func FormatSeoLiveTV(livetvId, slug, title string, cache bool) SeoObjectStruct {
	var titleUpFirst = strings.Title(title)
	var SeoObject SeoObjectStruct
	SeoObject.Url = fmt.Sprintf(SEO_LIVETV_URL, slug)
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = fmt.Sprintf(SEO_LIVETV_TITLE, titleUpFirst)
	SeoObject.Description = fmt.Sprintf(SEO_LIVETV_DESCRIPTION, title)
	SeoObject.Title_seo_tag = SeoObject.Title
	SeoObject.Seo_text = SeoObject.Description
	SeoObject.Canonical_tag = "https://vieon.vn/truyen-hinh-truc-tuyen/"
	SeoObject.Meta_robots = "index,follow"

	infoSEO, err := GetSEO(livetvId, cache)
	if err != nil || infoSEO.Slug == "" {
		return SeoObject
	}
	SeoObject.Url = infoSEO.Slug
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = infoSEO.Title
	SeoObject.Title_seo_tag = infoSEO.Title_seo_tag
	SeoObject.Description = infoSEO.Meta_description
	SeoObject.Seo_text = infoSEO.Seo_text
	SeoObject.Canonical_tag = infoSEO.Canonical_tag
	SeoObject.Meta_robots = infoSEO.Meta_robots
	SeoObject.Alternate = infoSEO.Alternate
	SeoObject.Deeplink = infoSEO.Deeplink
	return SeoObject
}

func FormatSeoEPG(slugChannel, slugEpg, titleChannel, titleEpg string) SeoObjectStruct {
	var SeoObject SeoObjectStruct
	SeoObject.Url = fmt.Sprintf(SEO_EPG_URL, slugChannel, slugEpg)
	SeoObject.Share_url = HandleSeoShareUrlVieON(SeoObject.Url)
	SeoObject.Title = fmt.Sprintf(SEO_EPG_TITLE, titleEpg)
	SeoObject.Description = fmt.Sprintf(SEO_EPG_DESCRIPTION, titleEpg, titleChannel)
	SeoObject.Title_seo_tag = SeoObject.Title
	SeoObject.Seo_text = SeoObject.Description
	SeoObject.Canonical_tag = "https://vieon.vn/truyen-hinh-truc-tuyen/"
	SeoObject.Meta_robots = "index,follow"
	return SeoObject
}

func FormatSeoEpisodeList(title, url string) SeoObjectStruct {
	var SeoObject SeoObjectStruct
	SeoObject.Url = fmt.Sprintf(SEO_EPISODE_LIST_URL, url)
	SeoObject.Title = fmt.Sprintf(SEO_EPISODE_LIST_TITLE, title)
	SeoObject.Description = fmt.Sprintf(SEO_EPISODE_LIST_DESCRIPTION, title)
	return SeoObject
}

func FormatSeoRelated(title, url string) SeoObjectStruct {
	var SeoObject SeoObjectStruct
	SeoObject.Url = fmt.Sprintf(SEO_RELATED_URL, url)
	SeoObject.Title = fmt.Sprintf(SEO_RELATED_TITLE, title)
	SeoObject.Description = fmt.Sprintf(SEO_RELATED_DESCRIPTION, title)
	return SeoObject
}
