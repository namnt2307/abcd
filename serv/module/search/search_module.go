package search

import (
	"regexp"
	"strings"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	recommendation "cm-v5/serv/module/recommendation"
	"github.com/gosimple/slug"
)

func RemoveSpecialCharacters(txt string) string {
	str := strings.Trim(txt, `~!@#$%^&*()_+/*{}Ơ|:"<>?ơ]\;',./[]`)
	str = strings.TrimSpace(str)
	return str
}

func BuildQuery(keyword, tags, platform, version string, entityType int, cacheActive bool) string {

	//Lay danh sach tap trong keyword
	re := regexp.MustCompile("[0-9]+")
	epsNum := re.FindAllString(keyword, -1)
	if epsNum == nil {
		epsNum = []string{}
	}

	//Xoa cac chu va ky tu dac biet
	keyword = RemoveSpecialCharacters(keyword)
	if keyword == "" {
		return ""
	}

	// re = regexp.MustCompile(SEARCHREGEXEPS)
	// keyword = strings.ToLower(keyword)
	// keyword = re.ReplaceAllString(keyword, "")

	//Chi cho phep dia chi ip viet nam tim kiem\
	var mustFilter []string
	// mustFilter = append(mustFilter, `{"term": {"geo_check": 0}},`)

	//Filter theo tags
	tagsFilter := TagsFilter(tags)
	if tagsFilter != "" {
		mustFilter = append(mustFilter, tagsFilter)
	}

	//Filter theo platform
	// platform_filter = _platform_filter(platform_id)
	// must_filter.append(platform_filter)
	var platformInfo = Platform(platform)
	// mustFilter = append(mustFilter, platformFilter)

	var mustFilterString = ""
	for _, val := range mustFilter {
		mustFilterString += val
	}
	//should
	var shouldFilters []string
	var types []int
	var qString = ""
	if entityType == 0 {
		types = []int{1, 3, 5}
		typesEpS := []int{4}

		if len(epsNum) > 0 {
			//Build query search theo episode
			var mustFilterEps []string
			mustFilterEps = append(mustFilterEps, mustFilterString)
			mustFilterEps = append(mustFilterEps, TermsFilter("type", typesEpS)+",")
			mustFilterEps = append(mustFilterEps, TermsFilter("episode", epsNum)+",")
			mustFilterEps = append(mustFilterEps, BuildMatchKeywordsQuery(keyword, `["keywords.keyword^5"]`, true))
			qString = ""
			for _, val := range mustFilterEps {
				qString += val
			}
			shouldFilters = append(shouldFilters, `{"bool": {"must":[ `+qString+`]}},`)
		}
	} else {
		types = []int{0}
	}

	if platformInfo.Type != "smarttv" && entityType == 0 {
		types = append(types, -2)
	}

	typeFilter := TermsFilter("type", types)

	//Build query search theo artist
	var mustFilterArt []string
	mustFilterArt = append(mustFilterArt, mustFilterString)
	mustFilterArt = append(mustFilterArt, typeFilter+",")
	mustFilterArt = append(mustFilterArt, BuildMatchKeywordsQuery(keyword, `["keywords.keyword^2", "keywords.artist^2", "keywords.description"]`, true))
	qString = ""
	for _, val := range mustFilterArt {
		qString += val
	}
	shouldFilters = append(shouldFilters, `{"bool": {"must": [`+string(qString)+`]}},`)

	//Build query search theo ext
	var mustFilterExt []string
	mustFilterExt = append(mustFilterExt, mustFilterString)
	mustFilterExt = append(mustFilterExt, typeFilter+",")
	mustFilterExt = append(mustFilterExt, BuildContainKeywordsQuery(keyword, `["keywords.keyword"]`))
	qString = ""
	for _, val := range mustFilterExt {
		qString += val
	}
	shouldFilters = append(shouldFilters, `{"bool": {"must": [`+string(qString)+`]}}`)

	qString = ""
	for _, val := range shouldFilters {
		qString += val
	}

	qString = strings.Replace(qString, `]"`, "]", -1)
	qString = strings.Replace(qString, `"[`, "[", -1)
	var query = `{
        "bool": {
            "minimum_should_match": 1,
            "should":[ ` + qString + `]
        }
	}`

	return query
}

func TermsFilter(field string, types interface{}) string {
	dataByte, err := json.Marshal(types)
	if err != nil {
		return ""
	}

	return `{
		"terms":{
			"` + field + `":"` + string(dataByte) + `"
		}
	}`
}

func BuildMatchKeywordsQuery(keyword, fields string, isPhrase bool) string {
	var pharase string
	if isPhrase {
		pharase = "phrase"
	} else {
		pharase = "best_fields"
	}
	return `{
        "nested": {
            "path": "keywords",
            "query": {
                "multi_match": {
                    "query": "` + keyword + `",
                    "type": "` + pharase + `",
                    "fields": "` + fields + `"
                }
            }
        }
    }`
}

func BuildContainKeywordsQuery(keyword, fields string) string {
	return `{
        "nested": {
            "path": "keywords",
            "query": {
                "query_string": {
                    "query": "*` + keyword + `*",
                    "fields": ` + fields + `
                }
            }
        }
    }`
}

func TagsFilter(tags string) string {
	if tags == "" {
		return ""
	}
	var arrTag interface{}
	arrTag = strings.Split(tags, ",")
	dataByte, err := json.Marshal(arrTag)
	if err != nil {
		return ""
	}
	// strTag := strings.Replace(string(dataByte), `"`, `'`, -1)
	return `{
		"terms":{
			"tags":` + string(dataByte) + `
		}
	},`
}

func GetTrackingDataSearch(platform, Type string) TrackingDataStruct {
	var TrackingData TrackingDataStruct
	trackingType := platform + "_" + Type
	trackingType = slug.Make(trackingType)
	TrackingData = recommendation.GetRandomDefaultTrackingData(strings.ToUpper(trackingType))
	return TrackingData
}
