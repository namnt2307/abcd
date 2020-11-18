package drm_auth

import (
	"encoding/hex"
	"fmt"
	"log"

	// . "cm-v5/serv/module"
	livetv "cm-v5/serv/module/livetv_v3"
	packages "cm-v5/serv/module/packages"
	vod "cm-v5/serv/module/vod"
)

const (
	DRM_VIEON_TYPE_VOD    = "vod"
	DRM_VIEON_TYPE_LIVETV = "livetv"
	DRM_VIEON_NAME        = "vieon"
)

func CheckValidLicenseVieON(licenseRequest VieONLicenseRequest, statusUserIsPremium int) bool {

	isValid, _ := packages.CheckPermissionUser(licenseRequest.Id, licenseRequest.User_id, statusUserIsPremium)
	return isValid
}

func GenerateInfoResponseVieON(licenseRequest VieONLicenseRequest) string {
	log.Println("GenerateInfoResponseVieON start")
	log.Println("licenseRequest.Type", licenseRequest.Type)

	if licenseRequest.Type == DRM_VIEON_TYPE_LIVETV {
		var listIds []string = []string{licenseRequest.Id}

		infoLivetv, err := livetv.GetLiveTVByListID(listIds, 0, true)
		fmt.Println("infoLivetv", infoLivetv)
		if err != nil || len(infoLivetv) == 0 {
			log.Println("GenerateInfoResponseVieON err", err)
			return ""
		}

		if infoLivetv[0].Drm_service_name == DRM_VIEON_NAME && infoLivetv[0].Drm_key != "" {
			decoded, err := hex.DecodeString(infoLivetv[0].Drm_key)
			if err != nil {
				log.Println("DecodeString err", err)
				return ""
			}
			return fmt.Sprintf("%s", decoded)
		}
		log.Println("Livetv get drm err", licenseRequest.Id)
	} else if licenseRequest.Type == DRM_VIEON_TYPE_VOD {
		infoVod, err := vod.GetVodDetail(licenseRequest.Id, 0, true)
		if err != nil {
			log.Println("GenerateInfoResponseVieON err", err)
			return ""
		}

		if infoVod.Drm_service_name == DRM_VIEON_NAME && infoVod.Drm_vieon_key != "" {
			decoded, err := hex.DecodeString(infoVod.Drm_vieon_key)
			if err != nil {
				log.Println("DecodeString err", err)
				return ""
			}

			return fmt.Sprintf("%s", decoded)
		}
		log.Println("Vod get drm err", licenseRequest.Id)
	}
	return ""
}
