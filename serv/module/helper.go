package module

import (
	"bytes"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	. "cm-v5/schema"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func PrintConsole(obj interface{}) {
	log.Println(obj)
	return
}

func EncodeUUID4Hex() string {
	uuid := UUIDV4Generate()
	uuidHex := strings.Replace(uuid, "-", "", -1)

	return uuidHex
}

type DataResultStruct struct {
	Status int
	Data   interface{}
	Error  string
}

func HandlePanic() {
	if err := recover(); err != nil {
		log.Println("HandlePanic:", err)
	}
	return
}

func BuildImage(pathImg string) string {
	if pathImg == "" {
		return ""
	}
	return DOMAIN_IMAGE_CDN + "/" + pathImg
}

func MappingImagesV4(platform string, formatOutput ImagesOutputObjectStruct, imagesSources ImageAllPlatformStruct, onlyGetThumbnail bool) ImagesOutputObjectStruct {

	//Lấy hình V4 theo từng platform
	switch platform {
	case "web":
		formatOutput.Carousel_web_v4 = BuildImage(imagesSources.Web.Carousel_web_v4)
		formatOutput.Poster_v4 = BuildImage(imagesSources.Web.Poster_v4)
		formatOutput.Thumbnail_big_v4 = BuildImage(imagesSources.Web.Thumbnail_big_v4)
		formatOutput.Thumbnail_hot_v4 = BuildImage(imagesSources.Web.Thumbnail_hot_v4)
		formatOutput.Thumbnail_v4 = BuildImage(imagesSources.Web.Thumbnail_v4)
		//lấy lại hình cũ nếu chưa có hình v4
		if formatOutput.Carousel_web_v4 == "" {
			formatOutput.Carousel_web_v4 = formatOutput.Home_carousel_web
		}
		if formatOutput.Thumbnail_v4 == "" {
			formatOutput.Thumbnail_v4 = formatOutput.Vod_thumb
		}
		formatOutput.Promotion_banner = BuildImage(imagesSources.Web.Promotion_banner)
		formatOutput.Title_card_light = BuildImage(imagesSources.Web.Title_card_light)
		formatOutput.Title_card_dark = BuildImage(imagesSources.Web.Title_card_dark)
	case "smarttv":
		formatOutput.Carousel_tv_v4 = BuildImage(imagesSources.Smarttv.Carousel_tv_v4)
		formatOutput.Poster_v4 = BuildImage(imagesSources.Smarttv.Poster_v4)
		formatOutput.Thumbnail_big_v4 = BuildImage(imagesSources.Smarttv.Thumbnail_big_v4)
		formatOutput.Thumbnail_hot_v4 = BuildImage(imagesSources.Smarttv.Thumbnail_hot_v4)
		formatOutput.Thumbnail_v4 = BuildImage(imagesSources.Smarttv.Thumbnail_v4)
		//lấy lại hình cũ nếu chưa có hình v4
		if formatOutput.Carousel_tv_v4 == "" {
			formatOutput.Carousel_tv_v4 = formatOutput.Home_carousel_tv
		}
		formatOutput.Promotion_banner = BuildImage(imagesSources.Smarttv.Promotion_banner)
		formatOutput.Title_card_light = BuildImage(imagesSources.Smarttv.Title_card_light)
		formatOutput.Title_card_dark = BuildImage(imagesSources.Smarttv.Title_card_dark)
	case "app":
		formatOutput.Carousel_web_v4 = BuildImage(imagesSources.Web.Carousel_web_v4)
		formatOutput.Carousel_app_v4 = BuildImage(imagesSources.App.Carousel_app_v4)
		formatOutput.Poster_v4 = BuildImage(imagesSources.App.Poster_v4)
		formatOutput.Thumbnail_big_v4 = BuildImage(imagesSources.App.Thumbnail_big_v4)
		formatOutput.Thumbnail_hot_v4 = BuildImage(imagesSources.App.Thumbnail_hot_v4)
		formatOutput.Thumbnail_v4 = BuildImage(imagesSources.App.Thumbnail_v4)
		//lấy lại hình cũ nếu chưa có hình v4
		if formatOutput.Carousel_app_v4 == "" {
			formatOutput.Carousel_app_v4 = formatOutput.Banner
		}
		formatOutput.Promotion_banner = BuildImage(imagesSources.App.Promotion_banner)
		formatOutput.Title_card_light = BuildImage(imagesSources.App.Title_card_light)
		formatOutput.Title_card_dark = BuildImage(imagesSources.App.Title_card_dark)
	}

	if onlyGetThumbnail {
		formatOutput.Thumbnail_big_v4 = ""
		formatOutput.Thumbnail_hot_v4 = ""
		formatOutput.Carousel_app_v4 = ""
		formatOutput.Carousel_tv_v4 = ""
		formatOutput.Carousel_web_v4 = ""
	} else {
		//lấy lại hình cũ nếu chưa có hình v4
		if formatOutput.Thumbnail_big_v4 == "" {
			formatOutput.Thumbnail_big_v4 = formatOutput.Vod_thumb_big
		}
		if formatOutput.Thumbnail_hot_v4 == "" {
			formatOutput.Thumbnail_hot_v4 = formatOutput.Home_vod_hot
		}
	}

	if formatOutput.Poster_v4 == "" {
		formatOutput.Poster_v4 = formatOutput.Poster
	}
	if formatOutput.Thumbnail_big_v4 == "" {
		formatOutput.Thumbnail_big_v4 = formatOutput.Vod_thumb_big
	}
	if formatOutput.Thumbnail_hot_v4 == "" {
		formatOutput.Thumbnail_hot_v4 = formatOutput.Home_vod_hot
	}

	if formatOutput.Thumbnail_v4 == "" {
		formatOutput.Thumbnail_v4 = formatOutput.Thumbnail
	}

	return formatOutput
}

func SupportUrlImg(img_url string, size string) string {
	if img_url == "" {
		return ""
	}
	var url string
	arr_url := strings.Split(img_url, ".")
	for index, val := range arr_url {
		if index == len(arr_url)-1 {
			url = fmt.Sprintf(`%s_%s.%s`, url, size, val)
		} else {
			url = fmt.Sprintf(`%s%s`, url, val)
		}
	}
	return url
}

func HandleSeoTitleVieON(txt string) string {
	txt = strings.Replace(txt, "- ViePlay", "- VieON", -1)
	return txt
}

func HandleLinkPlayVieON(txt string) string {
	txt = strings.Replace(txt, ".vieplay.vn/", ".vieon.vn/", -1)
	return txt
}

func HandleLinkThumbVtt(link string) string {
	link = strings.Replace(link, "/playlist.m3u8", "/thumbs.vtt", -1)
	return link
}
func HandleLinkThumbImage(link string) string {
	link = strings.Replace(link, "/playlist.m3u8", "/thumbs.jpg", -1)
	return link
}

func HandleLinkPlayVodByPremium(linkPlayRoot string, is_premium bool) string {
	// Disable this func
	return linkPlayRoot

	if strings.Contains(linkPlayRoot, "?DVR") {
		// return linkPlayRoot
		if is_premium == true {
			return linkPlayRoot
		}
		linkPlayRoot = strings.Replace(linkPlayRoot, ".m3u8", "_free.m3u8", -1)
		linkPlayRoot = strings.Replace(linkPlayRoot, ".mpd", "_free.mpd", -1)
		return linkPlayRoot
	}

	return linkPlayRoot
	// if is_premium == true {
	// 	linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.m3u8", "playlist-vip.m3u8", -1)
	// 	linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.mpd", "playlist-vip.mpd", -1)
	// 	return linkPlayRoot
	// }
	// linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.m3u8", "playlist-std.m3u8", -1)
	// linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.mpd", "playlist-std.mpd", -1)
	// return linkPlayRoot
}

func GetLinkPlayEpg(linkPlayRoot string, timeStart, duration int) string {
	nameFile := fmt.Sprintf("playlist_dvr_range-%d-%d.m3u8", timeStart, duration)
	return strings.Replace(linkPlayRoot, "playlist.m3u8", nameFile, -1)
}

func HandleShareUrlVieON(txt string) string {
	txt = strings.Replace(txt, "https://web.vieon.vn/", "http://vieon.vn/", -1)
	return txt
}

func HandleSeoShareUrlVieON(seo_url string) string {
	return WEB_VIEON + seo_url
}

// func HandleLinkPlayVodByPremium(linkPlayRoot string, is_premium bool) string {
// 	return linkPlayRoot
// 	if is_premium == true {
// 		linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.m3u8", "playlist-vip.m3u8", -1)
// 		linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.mpd", "playlist-vip.mpd", -1)
// 		return linkPlayRoot
// 	}
// 	linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.m3u8", "playlist-std.m3u8", -1)
// 	linkPlayRoot = strings.Replace(linkPlayRoot, "playlist.mpd", "playlist-std.mpd", -1)
// 	return linkPlayRoot
// }

// func GetLinkPlayEpg(linkPlayRoot string, timeStart, duration int) string {
// 	nameFile := fmt.Sprintf("playlist_dvr_range-%d-%d.m3u8", timeStart, duration)
// 	return strings.Replace(linkPlayRoot, "playlist.m3u8", nameFile, -1)
// }

func FormatResultAPI(status int, errStr string, data interface{}) interface{} {
	var dataR gin.H
	// if reflect.ValueOf(data).IsNil() && status == 200 {
	// 	dataR = gin.H{"data": []string{}}
	// 	return dataR["data"]
	// }

	if status == 400 {
		dataR = gin.H{
			"message": errStr,
			"data":    data,
			"error":   status,
		}
		return dataR
	}
	dataR = gin.H{"data": data}
	return dataR["data"]
}

func GenerateKeyLog(action string, data []byte) string {
	time_current := time.Now().String()

	h := md5.New()
	io.WriteString(h, time_current+action+string(data)+strconv.Itoa(rand.Int())+strconv.Itoa(rand.Int())+strconv.Itoa(rand.Int())+strconv.Itoa(rand.Int())+strconv.Itoa(rand.Int()))
	var key = fmt.Sprintf("%x", h.Sum(nil))
	return key
}

func GetMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	var key = fmt.Sprintf("%x", h.Sum(nil))
	return key
}

func Expire_second_one_day() int {
	return 86400
}

func Expire_second_one_month() int {
	return 86400 * 30
}

func Expire_second_one_hour() int {
	return 3600
}

func error_catch(err error) {
	if err != nil {
		// fmt.Println(err)
	}
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func Sentry_log(err error) {
	if err != nil {
		raven.SetEnvironment(ENV)
		raven.SetTagsContext(map[string]string{"server_ip": GetLocalIP()})
		raven.CaptureError(err, nil)
	}
}

func Sentry_log_request_cache(c *gin.Context) {
	packet := &raven.Packet{
		Message:     "Update cache: " + c.Request.URL.Path,
		Level:       raven.INFO,
		Environment: ENV,
		Extra: map[string]interface{}{
			"URL": c.Request.URL,
		},
		Tags: []raven.Tag{
			raven.Tag{
				Key:   "server_ip",
				Value: GetLocalIP(),
			},
		},
	}
	raven.Capture(packet, nil)
}

func Sentry_log_with_msg(txt string) {
	packet := &raven.Packet{
		Message:     "Log: " + txt,
		Level:       raven.INFO,
		Environment: ENV,
		Tags: []raven.Tag{
			raven.Tag{
				Key:   "server_ip",
				Value: GetLocalIP(),
			},
		},
	}
	raven.Capture(packet, nil)
}

func Sentry_log_mysql_miss(txt string) {
	packet := &raven.Packet{
		Message:     "Log: " + txt + " - Host: " + hostMySQL + "-" + dbnameMySQL,
		Level:       raven.INFO,
		Environment: ENV,
		Tags: []raven.Tag{
			raven.Tag{
				Key:   "server_ip",
				Value: GetLocalIP(),
			},
		},
	}
	raven.Capture(packet, nil)
}

func Sentry_log_jobs(txt string) {
	raven.SetEnvironment(ENV)
	raven.SetTagsContext(map[string]string{"server_ip": GetLocalIP()})
	raven.CaptureMessageAndWait("Log job: "+txt, map[string]string{"category": "logging", "level": "info"})
}

func Sentry_log_slow_request(uri string) {
	packet := &raven.Packet{
		Message:     "Slow Request: " + uri,
		Level:       raven.INFO,
		Environment: ENV,
		Tags: []raven.Tag{
			raven.Tag{
				Key:   "server_ip",
				Value: GetLocalIP(),
			},
		},
	}
	raven.Capture(packet, nil)
}

func FormatResultString(status int, errStr string, data interface{}) string {
	var dataResult DataResultStruct
	dataResult.Status = status
	dataResult.Error = errStr
	dataResult.Data = data
	dataResultStr, _ := json.Marshal(dataResult)
	return string(dataResultStr)
}

func GetCurrentTimeStamp() (string, string) {
	t := time.Now()

	// rand.Seed(time.Now().UTC().UnixNano())
	// day := rand.Intn(10)
	// t = t.AddDate(0, 0, -day)

	return fmt.Sprint(t.Unix()), t.Format("20060102")
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func Platform(platform string) PlatformOutputStruct {

	var PlatformOutput PlatformOutputStruct
	switch platform {
	case "all":
		PlatformOutput.Id = 0
		PlatformOutput.Type = "all"
	case "app":
		PlatformOutput.Id = 1
		PlatformOutput.Type = "app"
	case "smarttv":
		PlatformOutput.Id = 2
		PlatformOutput.Type = "smarttv"
	case "web":
		PlatformOutput.Id = 3
		PlatformOutput.Type = "web"
	case "ios":
		PlatformOutput.Id = 4
		PlatformOutput.Type = "app"
	case "android":
		PlatformOutput.Id = 5
		PlatformOutput.Type = "app"
	case "samsung_tv":
		PlatformOutput.Id = 6
		PlatformOutput.Type = "smarttv"
	case "sony_androidtv":
		PlatformOutput.Id = 7
		PlatformOutput.Type = "smarttv"
	case "lg_tv":
		PlatformOutput.Id = 8
		PlatformOutput.Type = "smarttv"
	case "mobile_web":
		PlatformOutput.Id = 9
		PlatformOutput.Type = "web"
	case "androidtv":
		PlatformOutput.Id = 12
		PlatformOutput.Type = "smarttv"
	}

	// var mRedis RedisModelStruct
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// key_cache := "key_" + platform
	// dataRedis, err := mRedis.GetString(key_cache)
	// //Get data to cache
	// if err == nil {
	// 	json.Unmarshal([]byte(dataRedis), &PlatformOutput)
	// } else {
	// 	// Get data in DB
	// 	session, db, err := GetCollection()
	// 	if err != nil {
	// 		return PlatformOutput
	// 	}
	// 	defer session.Close()
	// 	var where = bson.M{
	// 		"slug": platform,
	// 	}
	// 	db.C(COLLECTION_PLATFORM).Find(where).One(&PlatformOutput)

	// 	// Write cache
	// 	dataByte, _ := json.Marshal(PlatformOutput)
	// 	mRedis.SetString(key_cache, string(dataByte), 0)

	// }

	return PlatformOutput
}

func StringToInt(str string) (int, error) {
	if str == "" {
		return 0, errors.New("Not Parseint")
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return i, errors.New("Not Parseint")
	}
	return i, err
}

func StringToBool(str string) bool {
	if str == "true" {
		return true
	} else {
		return false
	}
}

func In_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func Isset(arr []string, index int) bool {
	return (len(arr) > index)
}

//GetStartEndTimeFromTimestamp timeInput = 0 to get current date
func GetStartEndTimeFromTimestamp(timeInput int64) (int, int) {
	timeStamp := time.Now().Local().Unix()
	if timeInput > 0 {
		timeStamp = timeInput
	}
	start := time.Unix(timeStamp, 0)
	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	endTime := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, time.Local).Unix()
	return int(startTime), int(endTime)
}

//GetTimeStampStartEndFromString timeInputYmd type Y-m-d
func GetTimeStampStartEndFromString(timeInputYmd string) (int, int) {
	const layoutISO = "2006-01-02"
	t, err := time.Parse(layoutISO, timeInputYmd)
	if err != nil {
		return 0, 0
	}

	start := time.Unix(t.Local().Unix(), 0)
	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	endTime := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, time.Local).Unix()
	return int(startTime), int(endTime)
}

func GetKeyCacheCurrentMinuteLiveTV(mAdd int) (string, time.Time) {
	t := time.Now().Add(time.Duration(mAdd) * time.Minute)
	y := t.Year()
	mon := t.Month().String()[:3]
	d := t.Day()
	h := t.Hour()
	m := t.Minute()

	txt := fmt.Sprintf("%d/%s/%02d-%02d:%02d", y, mon, int(d), h, m)
	return txt, t
}

func RandString(length int, str string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = str[rand.Intn(len(str))]
	}
	return string(b)
}

func RandNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandStringAny(str string, min, max int) string {
	randomNum := RandNumber(min, max)
	return RandString(randomNum, str)
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func ParseStringToArray(text string) ([]string, error) {
	var data []string
	err := json.Unmarshal([]byte(text), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetTagType(tagType int) string {
	switch tagType {
	case GENRE:
		return "genre"
	case COUNTRY:
		return "country"
	case CATEGORY:
		return "category"
	case LIVETV:
		return "livetv"
	case AUDIO:
		return "audio"
	}
	return "tag"
}

func GetTagTypeInt(tagType string) int {
	switch tagType {
	case "genre":
		return GENRE
	case "country":
		return COUNTRY
	case "category":
		return CATEGORY
	case "livetv":
		return LIVETV
	case "audio":
		return AUDIO
	}
	return TAG
}

func HttpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

func StringToTimestamp(str string) int64 {
	layout := "2006-01-02 15:04:05.000000"
	// str := "2019-03-05 13:33:31.000000"
	t, _ := time.Parse(layout, str)

	return int64(t.Unix())
}

// GetClientIPHelper gets the client IP using a mixture of techniques.
// This is how it is with golang at the moment.
func GetClientIPHelper(req *http.Request, c *gin.Context) (ipResult string, errResult error) {
	//Check ip white list
	if req.Header.Get("vie_whitelist") != "" {
		return "127.0.0.1", nil
	}

	if req.Header.Get("vie_adminn") != "" {
		return "127.0.0.1", nil
	}

	ip := c.ClientIP()
	if ip != "" {
		return ip, nil
	}
	// Try Request Headers (X-Forwarder). Client could be behind a Proxy
	ip, err := GetClientIPByHeaders(req)
	if err == nil {
		return ip, nil
	}

	//  Try Request Header ("Origin")
	url, err := url.Parse(req.Header.Get("Origin"))
	if err == nil {
		host := url.Host
		ip, _, err := net.SplitHostPort(host)
		if err == nil {
			return ip, nil
		}
	}

	// Try by Request
	ip, err = GetClientIPByRequestRemoteAddr(req)
	if err == nil {
		return ip, nil
	}

	err = errors.New("error: Could not find clients IP address")
	return "", err
}

// getClientIPByRequest tries to get directly from the Request.
// https://blog.golang.org/context/userip/userip.go
func GetClientIPByRequestRemoteAddr(req *http.Request) (ip string, err error) {
	// Try via request
	ip, _, err = net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		message := fmt.Sprintf("debug: Parsing IP from Request.RemoteAddr got nothing.")
		return "", fmt.Errorf(message)

	}
	return userIP.String(), nil
}

// getClientIPByHeaders tries to get directly from the Request Headers.
// This is only way when the client is behind a Proxy.
func GetClientIPByHeaders(req *http.Request) (ip string, err error) {

	// Client could be behid a Proxy, so Try Request Headers (X-Forwarder)
	ipSlice := []string{}

	ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
	ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
	ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

	for _, v := range ipSlice {
		if v != "" {
			return v, nil
		}
	}
	err = errors.New("error: Could not find clients IP address from the Request Headers")
	return "", err
}

func BodyReaderToString(req *http.Request) (string, error) {
	if req.Body == nil {
		return "", errors.New("Error")
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBuffer(body)
	req.Body = ioutil.NopCloser(buff)
	return string(body), nil
}

func UUIDV4() []byte {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println("UUIDV4:", err)
		return out
	}
	return out
}

func UUIDV4Generate() string {

	b := make([]byte, 16)
	_, err := cryptoRand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func isValidURL(stringTest string) bool {
	_, err := url.ParseRequestURI(stringTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

func CheckPlatformUsingLinkplayH265(platform, modelPlatform string) bool {

	if platform == "lg_tv" {
		return true
	} else if platform == "samsung_tv" && len(modelPlatform) > 5 {
		// T – 2020
		// R – 2019
		// N – 2018
		// M – 2017
		// K – 2016
		// J – 2015
		// H – 2014
		// F – 2013
		//https://en.tab-tv.com/?page_id=7123

		modelYear := modelPlatform[4]
		model2018 := byte('N') //model 2018

		//chỉ áp dụng đối với các model 2018 trở lên
		if modelYear >= model2018 {
			return true
		}
	}

	return false
}

func GetDateFromTimestamp(timestamp int64) string {
	t := time.Unix(int64(timestamp), 0)
	y := t.Year()
	mon := t.Month().String()[:3]
	d := t.Day()
	// layout := "02/Jan/2006:15:04:05"
	str := fmt.Sprintf("%02d/%s/%d", y, mon, int(d))
	return str
}

func DecodeDataformToString(str string) string {
	decodedValue, err := url.QueryUnescape(str)
	if err != nil {
		return ""
	}
	decodedValue = "{" + decodedValue + "}"
	decodedValue = strings.Replace(decodedValue, "=", ":", -1)
	return decodedValue
}

func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range PRIVATE_IPBLOCKS {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func ReplaceDomainMBF(input_path string) string {
	if input_path == "" {
		return input_path
	}

	resutlUri, err := url.Parse(input_path)
	if err != nil {
		return ""
	}

	CDNInfo := CheckPartner(input_path)
	if CDNInfo.Hostname_mobifone != "" {
		output_path := strings.Replace(input_path, resutlUri.Host, CDNInfo.Hostname_mobifone, -1)
		return output_path
	}
	return input_path
}

func GetExpiredDateFromDuration(duration int, duration_type string, outPutFormat string) string {
	if outPutFormat == "" {
		outPutFormat = "2006-01-02 15:04:05"
	}

	//default in hours
	if duration_type == "" {
		duration_type = "hours"
	}
	//get current time
	currentTime := time.Now().Local()

	//Ngày dời lại mặc định là 1/ Trừ các gói có duration type = months_exact/years_exact
	var dayToBack = -1
	if duration_type == "months_end_day" {
		dayToBack = 0
	}

	var expiredTime time.Time

	switch duration_type {
	case "hours":
		expiredTime = currentTime.Add(time.Duration(duration) * time.Hour)
	case "days":
		expiredTime = currentTime.AddDate(0, 0, duration)
	case "months", "months_end_day":

		_, currentMonth, _ := currentTime.Date()

		var expireMonth = (int(currentMonth) + duration)
		if expireMonth > 12 {
			expireMonth = expireMonth % 12
		}

		//trường hợp bình thường + theo công thức d-1/m+1
		expiredTime = currentTime.AddDate(0, duration, dayToBack)
		expiredYear, expiredMonth, _ := expiredTime.Date()
		expiredHours, expiredMinutes, expiredSeconds := expiredTime.Clock()

		//trường hợp tháng hiện tại + duration rơi vào tháng 2 và tháng hết hạn sau khi tính toán lại nhảy tới tháng 3 thì ngày hết hạn cần back lại về cuối tháng 2
		//trường hợp ngày mua là ngày cuối tháng thì ngày hết hạn và gia hạn cũng là ngày cuối tháng (chỉ áp dụng đối với months_end_day)
		if (expireMonth == 2 && expireMonth < int(expiredMonth)) || (duration_type == "months_end_day" && IsEndOfMonth(currentTime)) {
			expiredTime = time.Date(expiredYear, time.Month(expireMonth)+1, 0, expiredHours, expiredMinutes, expiredSeconds, 0, currentTime.Location())
		}

	case "years":
		expiredTime = currentTime.AddDate(duration, 0, dayToBack)
	}

	year, month, day := expiredTime.Date()
	expiredTime = time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	// outPutFormat := "2006-01-02 23:59:59"
	return expiredTime.Format(outPutFormat)
}

//ReCalculatorRating (minRate float64, currentRate float64) float64 {
func ReCalculatorRating(minRate float64, currentRate float64) float64 {
	if minRate <= 0 {
		minRate = 3
	}
	if currentRate < minRate {
		currentRate = minRate
	}
	if currentRate > 5 {
		currentRate = 5
	}
	return currentRate
}

func IsEndOfMonth(t time.Time) bool {
	_, m, _ := t.Date()
	_, month, _ := t.AddDate(0, 0, 1).Date()
	return int(month) != int(m)
}
