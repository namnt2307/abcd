package content

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	. "cm-v5/serv/module"
)

const (
	LINK_GET_VIEON_INFO = `%s/video/%s/%s/%s/info` // https: //<domain>/video/<client_id>/<content_id>/<secure_hash>/info
)

type VideoInfoObjectStruct struct {
	Code   int
	Result struct {
		Thumbs struct {
			Image    string
			Vtt      string
			Vtt_lite string
		}

		Map_profile []struct {
			Bandwidth int
			Name      string
			Width     int
			Height    int
		}

		Audios []struct {
			Id         string
			Code_name  string
			Is_default int64
			Title      string
		}

		Subtitles []struct {
			Id         string
			Code_name  string
			Is_default int64
			Title      string
			Uri        string
		}
		Link_play struct {
			Hls_link_play  string
			Dash_link_play string
		}
	}
	Error struct {
		Message string
	}
}

func BuildLinkGetVideoInfo(contentId string, rules RulesObjStruct) string {
	// secure hash
	secureHash := SecureHashIsMD5(contentId, rules.Platform)

	// Build link_get_video_info
	link_get_video_info := fmt.Sprintf(LINK_GET_VIEON_INFO, PLAY_LIST_HOST_NAME, PLAY_LIST_CLIENT_ID, contentId, secureHash)

	// Build query param link play
	uri, err := url.Parse(link_get_video_info)
	if err != nil {
		return ""
	}
	parameters := url.Values{}
	parameters.Add("user_type", rules.User_type)
	parameters.Add("platform", rules.Platform)
	parameters.Add("model", rules.Platform)
	uri.RawQuery = parameters.Encode()

	return uri.String()
}

func SecureHashIsMD5(content_id, platform string) string {
	strHash := content_id + "|" + platform + "|" + PLAY_LIST_CLIENT_SECRET
	algorithm := md5.New()
	algorithm.Write([]byte(strHash))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func GetVideoInfo(contentId, ipUser string, rules RulesObjStruct) (VideoInfoObj VideoInfoObjectStruct, err error) {
	url := BuildLinkGetVideoInfo(contentId, rules)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return VideoInfoObj, err
	}
	req.Header.Set("X-FORWARDED-FOR", ipUser)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return VideoInfoObj, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return VideoInfoObj, err
	}

	if res.StatusCode != 200 {
		return VideoInfoObj, errors.New("resp with status: " + res.Status)
	}

	err = json.Unmarshal(body, &VideoInfoObj)
	if err != nil {
		return VideoInfoObj, err
	}

	return VideoInfoObj, err
}
