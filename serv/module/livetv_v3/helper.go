package livetv_v3

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

type LiveTVInfoObjectStruct struct {
	Code   int
	Result struct {
		Live_links struct {
			Hls  string
			Dash string
		}
		Epg_links struct {
			Hls  string
			Dash string
		}
	}
	Error struct {
		Message string
	}
}

func SecureHashIsMD5(content_id string) string {
	strHash := content_id + "|" + PLAY_LIST_CLIENT_SECRET
	algorithm := md5.New()
	algorithm.Write([]byte(strHash))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func GetLiveTVInfo(contentId, ipUser string, timeStart, duration int) (LiveTVInfo LiveTVInfoObjectStruct, err error) {
	secureHash := SecureHashIsMD5(contentId)
	link := fmt.Sprintf("%s/livetv/%s/%s/%s/info", PLAY_LIST_HOST_NAME, PLAY_LIST_CLIENT_ID, contentId, secureHash)
	// Build query param link play
	uri, err := url.Parse(link)
	if err != nil {
		return LiveTVInfo, err
	}
	parameters := url.Values{}
	parameters.Add("timestamp", fmt.Sprint(timeStart))
	parameters.Add("duration", fmt.Sprint(duration))
	uri.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return LiveTVInfo, err
	}

	req.Header.Set("X-FORWARDED-FOR", ipUser)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return LiveTVInfo, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return LiveTVInfo, err
	}

	if res.StatusCode != 200 {
		return LiveTVInfo, errors.New("resp with status: " + res.Status)
	}

	err = json.Unmarshal(body, &LiveTVInfo)
	if err != nil {
		return LiveTVInfo, err
	}

	return LiveTVInfo, err
}
