package module

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	edgeauth "github.com/mobilerider/EdgeAuth-Token-Golang"
)

type CDNStruct struct {
	Hostname          string
	Protocol          string
	Secret_token      string
	Partner           string
	Hostname_mobifone string
}

var (
	KEYCACHE_CDN = "CDN_KEY"
	VNNETWORK    = "VNNETWORK"
	VINADATA     = "VINADATA"
	VINADATA2    = "VINADATA2"
	AKAMAI       = "AKAMAI"
)

func ConnectMySQLStreaming() (db *sql.DB, err error) {
	var instanceMySQLCli = new(MySQLCli)
	instanceMySQLCli.db, err = sql.Open("mysql", usernameMySQL+":"+passwordMySQL+"@tcp("+hostMySQL+":"+portMySQL+")/"+dbnameMySQLStream)
	if err != nil {
		log.Fatalf("Connect MySQL: %s\n", err)
		Sentry_log(err)
		return instanceMySQLCli.db, err
	}

	return instanceMySQLCli.db, nil
}

func init() {
	//Connect mysql
	DB, err := ConnectMySQLStreaming()
	if err != nil {
		log.Fatalf("Connect MySQL Streaming: %s\n", err)
		return
	}
	defer DB.Close()

	Rows, err := DB.Query("SELECT hostname, protocol, secret_token, partner, hostname_mobifone FROM Playlist_server")
	if err != nil {
		return
	}

	var ListCDN []CDNStruct
	for Rows.Next() {
		var CDN CDNStruct
		err = Rows.Scan(&CDN.Hostname, &CDN.Protocol, &CDN.Secret_token, &CDN.Partner, &CDN.Hostname_mobifone)
		if err != nil {
			continue
		}

		ListCDN = append(ListCDN, CDN)
	}

	// Write Redis
	dataByte, _ := json.Marshal(ListCDN)
	mRedis.SetString(KEYCACHE_CDN, string(dataByte), -1)
}

func CheckPartner(input_path string) CDNStruct {
	var CDN CDNStruct
	valueCache, err := mRedis.GetString(KEYCACHE_CDN)
	if err != nil && valueCache == "" {
		return CDN
	}

	var ListCDN []CDNStruct
	err = json.Unmarshal([]byte(valueCache), &ListCDN)
	if err != nil {
		return CDN
	}

	for _, val := range ListCDN {
		i := strings.Index(input_path, val.Hostname)
		if i != -1 {
			CDN = val
			break
		}
	}

	return CDN
}

func BuildTokenUrl(input_path, type_path, client_ip string) string {
	if input_path == "" {
		return ""
	}
	var uri string = input_path

	//Xoa query method get link play
	resutlUri, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	//Get domain link play
	domainLinkPlay := resutlUri.Scheme + "://" + resutlUri.Host

	if strings.Contains(input_path, "?DVR") == false {
		uri = strings.Replace(uri, "?"+resutlUri.RawQuery, "", -1)
	} else {
		type_path = "LiveTV"
	}

	uri = strings.Replace(uri, domainLinkPlay, "", -1)

	//build token
	CDNInfo := CheckPartner(input_path)
	var Result string
	switch CDNInfo.Partner {
	case VNNETWORK:
		Result = VnnetWorkBuildTokenUrl(CDNInfo.Secret_token, uri, client_ip)
	case VINADATA:
		Result = VinadataBuildTokenUrl_V2(CDNInfo.Secret_token, uri, type_path)
	case VINADATA2:
		Result = Vinadata2BuildTokenUrl(CDNInfo.Secret_token, resutlUri.Path, type_path)
	case AKAMAI:
		Result = AkamaiBuildTokenUrl(CDNInfo.Secret_token, uri)
	default:
		Result = uri
	}

	if Result == "" {
		return ""
	}

	return domainLinkPlay + Result
}

func AkamaiBuildTokenUrl(secret_token, input_path string) string {
	config := &edgeauth.Config{
		Algo:           crypto.SHA256,
		Key:            secret_token,
		DurationWindow: time.Second * 3600 * 9,
	}

	client := edgeauth.NewClient(config)

	token, err := client.GenerateToken(input_path, true)
	var output_path string = ""
	if err != nil {
		// Handle error
		return output_path
	}

	if strings.Contains(input_path, "?") == false {
		// Chưa có param bắt đầu
		output_path = input_path + "?hdnts=" + token
	} else {
		// Đã có param bắt đầu (?abc=111)
		output_path = input_path + "&hdnts=" + token
	}

	return output_path
}

func VinadataBuildTokenUrl_V2(secret_token, input_path, type_path string) string {
	if input_path == "" {
		return ""
	}

	now := time.Now() //.In(loc)
	nano := now.UnixNano()
	timestamp := fmt.Sprint(nano / 1000000)

	token := GetMd5(secret_token + "/" + timestamp)

	if type_path == "LiveTV" {
		return "/hls/" + token + "/" + timestamp + input_path
	}

	return "/" + token + "/" + timestamp + input_path
}

func Vinadata2BuildTokenUrl(secret_token, input_path, type_path string) string {
	if input_path == "" {
		return ""
	}

	path := input_path
	indexPath := strings.Index(input_path, "/playlist")
	if indexPath >= 0 {
		path = input_path[:indexPath]
	}

	now := time.Now()
	now = now.Add(time.Second * 3600 * 8)
	timeUnix := now.Unix()
	timestamp := fmt.Sprint(timeUnix * 1000)

	hash := secret_token + path + timestamp
	token := GetMd5(hash)
	if type_path == "LiveTV" {
		return fmt.Sprintf("/hls/%s/%s%s", token, timestamp, input_path)
	}

	return fmt.Sprintf("/%s/%s%s", token, timestamp, input_path)
}

func VnnetWorkBuildTokenUrl(secret_token, input_path, client_ip string) string {
	if input_path == "" {
		return ""
	}

	var g int = 0
	t := time.Now().UTC()
	var stime string = t.Format("20060102150405")
	t = t.Add(time.Second * 3600 * 9)
	var etime string = t.Format("20060102150405")

	if strings.Contains(input_path, "?") == false {
		// Chưa có param bắt đầu
		input_path = input_path + "?stime=" + stime + "&etime=" + etime
	} else {
		// Đã có param bắt đầu (?abc=111)
		input_path = input_path + "&stime=" + stime + "&etime=" + etime
	}

	GenerateToken := VnnetWorkGenerateToken(secret_token, input_path, g)
	return input_path + "&encoded=" + GenerateToken
}

func VnnetWorkGenerateToken(secret_token, resource string, generator int) string {
	hash := hmac.New(sha1.New, []byte(secret_token))
	io.WriteString(hash, resource)
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprint(generator) + sha1_hash[:20]
}

func VinadataGenerateToken(resource string, generator int) string {
	hash := hmac.New(sha1.New, []byte(STREAMING_SECRET_TOKEN))
	io.WriteString(hash, resource)
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprint(generator) + sha1_hash[:20]
}

func VinadataBuildTokenUrlLiveTV(input_path, client_ip string) string {
	if input_path == "" {
		return ""
	}

	var g int = 0
	t := time.Now().UTC()
	var stime string = t.Format("20060102150405")
	t = t.Add(time.Second * 3600 * 9)
	var etime string = t.Format("20060102150405")
	var uri string = input_path

	//Xoa query method get link play
	resutlUri, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	//Get domain link play
	domainLinkPlay := resutlUri.Scheme + "://" + resutlUri.Host

	uri = strings.Replace(uri, "?"+resutlUri.RawQuery, "", -1)
	uri = strings.Replace(uri, domainLinkPlay, "", -1)
	uri = uri + "?stime=" + stime + "&etime=" + etime

	GenerateToken := VinadataGenerateTokenLiveTV(uri, g)
	return domainLinkPlay + uri + "&encoded=" + GenerateToken
}

func VinadataGenerateTokenLiveTV(resource string, generator int) string {
	hash := hmac.New(sha1.New, []byte(STREAMING_SECRET_TOKEN_LIVETV))
	io.WriteString(hash, resource)
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprint(generator) + sha1_hash[:20]
}

func VieONBuildDrmTokenUrl(linkPlay, contentID, typeContent, ipUser, tokenUser string) string {
	//typeContent is livetv or vod
	//tokenUser is accesss token of user

	if linkPlay == "" {
		return ""
	}

	//Xoa query method get link play
	resutlURI, err := url.Parse(linkPlay)
	if err != nil {
		return ""
	}

	//remove all query string
	linkPlay = strings.Replace(linkPlay, "?"+resutlURI.RawQuery, "", -1)

	//add token type and id
	if typeContent == "livetv" {
		linkPlay = linkPlay + "?token=" + tokenUser + "&type=" + typeContent + "&id=" + contentID
	}

	return linkPlay
}
