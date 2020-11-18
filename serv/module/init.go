package module

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/dlintw/goconf"
	jsoniter "github.com/json-iterator/go"
)

var (
	ErrLoadConfig                  = errors.New("Cannot load file config")
	ErrConnectDatabase             = errors.New("Cannot connect Database")
	ErrLogin                       = errors.New("Cannot login")
	ErrDecode                      = errors.New("Decode Error")
	DOMAIN_IMAGE_CDN               string
	USE_MENCACHE                   string
	MENU_ID                        string
	PAGE_HOME_ID                   string
	STREAMING_SECRET_TOKEN         string
	STREAMING_SECRET_TOKEN_LIVETV  string
	ENV                            string
	
	STATUS_LISTEPISODE_GETUSERINFO string
	TTL_LOCALCACHE                 int64
	TTL_KVCACHE                    time.Duration
	ENABLE_KAFKA                   string
	ENABLE_WATCHMORE               string
	ENABLE_WATCHLATER              string
	PRIVATE_IPBLOCKS               []*net.IPNet
	LIST_ID_PACKAGE_MBF            []string
	WEB_VIEON                      string

	// PLAY_LIST
	PLAY_LIST_CLIENT_ID     string
	PLAY_LIST_CLIENT_SECRET string
	PLAY_LIST_HOST_NAME     string

	PRIVATE_IPWHILELIST_BLOCKS []*net.IPNet
	PRIVATE_IPWHILELIST_LIST   map[string]bool = make(map[string]bool, 0)
	LIST_TOP_VIEW_CONTENT                      = make(map[string]int)

	//LIVETV
	LIVETV_LIMIT_REQUEST int
	LIVETV_LIST_IP_LOCK  []string
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var mRedis RedisModelStruct

func init() {
	var err error
	DOMAIN_IMAGE_CDN, _ = CommonConfig.GetString("IMAGE", "domain_cdn")
	MENU_ID, _ = CommonConfig.GetString("MENU", "menu_id")
	PAGE_HOME_ID, _ = CommonConfig.GetString("MENU", "page_home_id")
	if PAGE_HOME_ID == "" {
		PAGE_HOME_ID = "cf73673f-b553-45ab-ab2f-329e3698543e"
	}

	STREAMING_SECRET_TOKEN, _ = CommonConfig.GetString("STREAMING", "secret_token")
	STREAMING_SECRET_TOKEN_LIVETV, _ = CommonConfig.GetString("STREAMING", "secret_token_livetv")
	ENV, _ = CommonConfig.GetString("GLOBAL_CONFIG", "env")
	if ENV == "" {
		ENV = "develop"
	}


	// Init Sentry
	SENTRY_DSN, _ := CommonConfig.GetString("SENTRY", "dsn")
	raven.SetDSN(SENTRY_DSN)

	// Global config
	STATUS_LISTEPISODE_GETUSERINFO, _ = CommonConfig.GetString("GLOBAL_CONFIG", "status_listepisode_getuserinfo")

	if STATUS_LISTEPISODE_GETUSERINFO == "" {
		STATUS_LISTEPISODE_GETUSERINFO = "true"
	}
	ttl, err := CommonConfig.GetInt("GLOBAL_CONFIG", "ttl_localcache")
	if err != nil {
		TTL_LOCALCACHE = 30
	} else {
		TTL_LOCALCACHE = int64(ttl)
	}

	ttl, err = CommonConfig.GetInt("GLOBAL_CONFIG", "ttl_kvcache")
	if err != nil {
		TTL_KVCACHE = 30
	} else {
		TTL_KVCACHE = time.Duration(ttl)
	}

	ENABLE_KAFKA, _ = CommonConfig.GetString("GLOBAL_CONFIG", "enable_kafka")
	ENABLE_WATCHMORE, _ = CommonConfig.GetString("GLOBAL_CONFIG", "enable_watchmore")
	ENABLE_WATCHLATER, _ = CommonConfig.GetString("GLOBAL_CONFIG", "enable_watchlater")

	list_pack_mbf, _ := CommonConfig.GetString("CONFIG_MBF", "list_pack")
	LIST_ID_PACKAGE_MBF = strings.Split(list_pack_mbf, ",")

	//Add ip MBF
	list_ip_mbf, _ := CommonConfig.GetString("CONFIG_MBF", "list_ip")
	arr_ip_mbf := strings.Split(list_ip_mbf, ",")
	for _, cidr := range arr_ip_mbf {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		PRIVATE_IPBLOCKS = append(PRIVATE_IPBLOCKS, block)
	}

	WEB_VIEON, _ = CommonConfig.GetString("GLOBAL_CONFIG", "web_vieon")

	// Link play
	PLAY_LIST_CLIENT_ID, _ = CommonConfig.GetString("PLAY_LIST", "client_id")
	PLAY_LIST_CLIENT_SECRET, _ = CommonConfig.GetString("PLAY_LIST", "client_secret")
	PLAY_LIST_HOST_NAME, _ = CommonConfig.GetString("PLAY_LIST", "host_name")

	// Livetv
	LIVETV_LIMIT_REQUEST, err = CommonConfig.GetInt("LIVETV", "limit_request")
	if err != nil || LIVETV_LIMIT_REQUEST == 0 {
		LIVETV_LIMIT_REQUEST = 20
	}

	list_ip_lock, _ := CommonConfig.GetString("LIVETV", "list_ip_lock")
	LIVETV_LIST_IP_LOCK = strings.Split(list_ip_lock, ",")

	//DT-14489
	ReadConfigIpWhileList()
}

func ReadConfigIpWhileList() {
	IPConfig, err := goconf.ReadConfigFile("config/list_ip.config")
	if err != nil {
		log.Println("ReadConfigFile list_ip err", err)
		return
	}

	//Add ip by blocks IP
	listIPBlocks, _ := IPConfig.GetString("IP_WHILELIST", "list_ip_blocks")
	if listIPBlocks != "" {
		arrIP := strings.Split(listIPBlocks, ",")
		for _, cidr := range arrIP {
			_, block, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Println("parse error on ", cidr, err)
				continue
			}
			PRIVATE_IPWHILELIST_BLOCKS = append(PRIVATE_IPWHILELIST_BLOCKS, block)
		}
	}

	//Add ip by list IP defined
	listIP, _ := IPConfig.GetString("IP_WHILELIST", "list_ip")
	if listIP != "" {
		arrIP := strings.Split(listIP, ",")
		for _, IP := range arrIP {
			PRIVATE_IPWHILELIST_LIST[IP] = true
		}
	}

	//Add ip by list IP range from x to y
	listIPRange, _ := IPConfig.GetString("IP_WHILELIST", "list_range_ip")
	if listIPRange != "" {
		arrIP := strings.Split(listIPRange, ",")
		for _, ipRange := range arrIP {
			rangeIPParse := strings.Split(ipRange, "::")
			if len(rangeIPParse) < 2 {
				log.Println("Invalid ip range format IP::[number]: ", ipRange)
				continue
			}
			ipArrParse := strings.Split(rangeIPParse[0], ".")
			if len(ipArrParse) < 3 {
				log.Println("Invalid ip format: ", rangeIPParse[0])
				continue
			}
			firstIndex, _ := StringToInt(ipArrParse[len(ipArrParse)-1])
			lastIndex, _ := StringToInt(rangeIPParse[1])

			for i := firstIndex; i <= lastIndex; i++ {
				ipArrParse[len(ipArrParse)-1] = fmt.Sprint(i)
				PRIVATE_IPWHILELIST_LIST[strings.Join(ipArrParse, ".")] = true
			}
		}
	}
}
