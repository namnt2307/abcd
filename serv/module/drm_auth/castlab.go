package drm_auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	. "cm-v5/serv/module"
	livetv "cm-v5/serv/module/livetv_v3"
	packages "cm-v5/serv/module/packages"
	vod "cm-v5/serv/module/vod"
)

var CONFIG_CASTLAB CastlabConfig

func init() {

	configStr, _ := CommonConfig.GetString("DRM", "castlab_config")
	err := json.Unmarshal([]byte(configStr), &CONFIG_CASTLAB)
	if err != nil {
		//nếu có lỗi thì lấy default value
		CONFIG_CASTLAB.Livetv_output_protect = true
		CONFIG_CASTLAB.Livetv_output_protect_analogue = true
		CONFIG_CASTLAB.Livetv_output_protect_digital = true
		CONFIG_CASTLAB.Livetv_output_protect_enforce = false
		CONFIG_CASTLAB.Livetv_store_licence = true
		CONFIG_CASTLAB.Save_request = true
		CONFIG_CASTLAB.Vod_output_protect = true
		CONFIG_CASTLAB.Vod_output_protect_analogue = true
		CONFIG_CASTLAB.Vod_output_protect_digital = true
		CONFIG_CASTLAB.Vod_output_protect_enforce = false
		CONFIG_CASTLAB.Vod_store_licence = true
	}
}

func CheckValidLicenseCastLab(licenseRequest CastlabLicenseRequest, isLiveTV bool) bool {
	if licenseRequest.Asset == "" || licenseRequest.User == "" || licenseRequest.Session == "" {
		log.Println("Empty data")
		return false
	}
	// decode token
	jwt, err := LocalAuthVerify(licenseRequest.Session)
	if err != nil {
		log.Println("decode jwt error", err)
		return false
	}
	statusUserIsPremium := jwt.Ispremium

	if licenseRequest.User != jwt.Subject {
		log.Println("user id not match")
		return false
	}

	assetID := licenseRequest.Asset

	//check if asset_id is livetv id
	if isLiveTV {
		// assetID = strings.Replace(assetID, "livetv_", "", 1)

		//get id livetv
		livetvID := livetv.GetIdLiveTVByAssetID(assetID)

		if livetvID == "" {
			log.Println("asset id not valid", assetID)
			Sentry_log(errors.New("Livetv asset id not valid: " + assetID))
			return false
		}

		//check exists livetv
		var listIds []string = []string{livetvID}
		infoLivetv, err := livetv.GetLiveTVByListID(listIds, 0, true)
		if err != nil || len(infoLivetv) == 0 {
			return false
		}
		assetID = livetvID
	} else {
		//check exists vod
		infoVod, err := vod.GetVodDetail(assetID, 0, true)
		if err != nil || infoVod.Id == "" {
			return false
		}
	}

	//check permission user. Only SVOD/SLivetv is allow to play castlab
	isValid, numPack := packages.CheckPermissionUser(assetID, jwt.Subject, statusUserIsPremium)
	if isValid && numPack > 0 {
		log.Println("CheckPermissionUser "+jwt.Subject, isValid)
		return isValid
	}

	log.Println("CheckPermissionUser Fail "+jwt.Subject+" | numPack :  ", numPack)
	return false
}

func GenerateInfoResponseCastlab(infoRequest CastlabLicenseRequest, isLiveTV bool) CastlabLicenseResponseSuccess {
	var infoResponse CastlabLicenseResponseSuccess

	// var uuid = UUIDV4Generate()
	// infoResponse.Ref = []string{UUIDV4Generate()}
	infoResponse.AssetID = infoRequest.Asset
	infoResponse.VariantID = infoRequest.Variant
	infoResponse.Profile = CreateProfileCastlab(isLiveTV)
	infoResponse.AccountingID = UUIDV4Generate()

	if isLiveTV {
		infoResponse.StoreLicense = CONFIG_CASTLAB.Livetv_store_licence
		if CONFIG_CASTLAB.Livetv_output_protect {
			var outputProtection OutputProtection
			outputProtection.Analogue = CONFIG_CASTLAB.Livetv_output_protect_analogue
			outputProtection.Digital = CONFIG_CASTLAB.Livetv_output_protect_digital
			outputProtection.Enforce = CONFIG_CASTLAB.Livetv_output_protect_enforce

			dataByte, err := json.Marshal(outputProtection)
			if err == nil {
				json.Unmarshal(dataByte, &infoResponse.OutputProtection)
			}
		}

	} else {
		infoResponse.StoreLicense = CONFIG_CASTLAB.Vod_store_licence
		if CONFIG_CASTLAB.Vod_output_protect {

			var outputProtection OutputProtection

			outputProtection.Analogue = CONFIG_CASTLAB.Vod_output_protect_analogue
			outputProtection.Digital = CONFIG_CASTLAB.Vod_output_protect_digital
			outputProtection.Enforce = CONFIG_CASTLAB.Vod_output_protect_enforce

			dataByte, err := json.Marshal(outputProtection)
			if err == nil {
				json.Unmarshal(dataByte, &infoResponse.OutputProtection)
			}
		}
	}

	return infoResponse
}

func CreateProfileCastlab(isLiveTV bool) CastlabProfile {
	// t := time.Now().UTC().Add(time.Minute * 60 * 24).Format("2006-01-02T15:04:05Z") //1 day
	// t := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	var profile CastlabProfile
	// profile.Rental.AbsoluteExpiration = t
	profile.Rental.RelativeExpiration = "P2D" //2 day from now

	if isLiveTV {
		profile.Rental.PlayDuration = 24 * 3600000 // 24hour in miliseconds
	} else {
		profile.Rental.PlayDuration = 3 * 3600000 // 3hour in miliseconds
	}

	return profile
}

func InsertRespLog(dataRequest interface{}, dataResponse interface{}, drmServiceName string) error {
	if CONFIG_CASTLAB.Save_request == false {
		return nil
	}
	session, db, err := GetCollection()
	if err != nil {
		fmt.Println("error connect db")
		return err
	}

	defer session.Close()

	var dataInsert DrmLogInfo
	dataInsert.Drm_service_name = drmServiceName
	dataInsert.Log_data.Request = dataRequest
	dataInsert.Log_data.Response = dataResponse
	dataInsert.Created_at = time.Now().Unix()

	err = db.C(COLLECTION_DRM_RESP_LOG).Insert(dataInsert)
	if err != nil {
		fmt.Println("InsertRespLog err ", err)
		return err
	}

	return nil
}
