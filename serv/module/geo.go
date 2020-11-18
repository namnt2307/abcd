package module

import (
	// "fmt"

	"log"
	"net"
	"strings"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var db *geoip2.Reader

func init() {
	db, _ = geoip2.Open("db/GeoIP2-Country.mmdb")
}

func CheckIpIsVN(ipStr string) bool {
	// return true
	if ipStr == "::1" {
		return true
	}

	ip := net.ParseIP(ipStr)
	//check ip whilelist
	if _, ok := PRIVATE_IPWHILELIST_LIST[ipStr]; ok {
		return true
	}

	for _, block := range PRIVATE_IPWHILELIST_BLOCKS {
		if block.Contains(ip) {
			return true
		}
	}

	//check ip is IP VN

	record, err := db.Country(ip)
	if err != nil {
		return false
	}
	if record.Country.IsoCode == "VN" {
		return true
	}
	return false
}

func CheckIpWhiteList(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	//check ip whilelist
	if _, ok := PRIVATE_IPWHILELIST_LIST[ipStr]; ok {
		return true
	}

	for _, block := range PRIVATE_IPWHILELIST_BLOCKS {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func CheckIpRule(ipStr string, arrIp []string) bool {
	ipCheck := net.ParseIP(ipStr)
	for _, cidr := range arrIp {
		if cidr == "0.0.0.0" {
			return true
		}

		if strings.Index(cidr, "/") == -1 {
			if cidr == ipStr {
				return true
			}
			continue
		}

		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Println("parse error on ", cidr, err)
			continue
		}

		if block.Contains(ipCheck) {
			return true
		}
	}
	return false
}
