package module

import (
	// "time"
	"errors"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	jwt.StandardClaims
	ExpiresAt int64  `json:"exp"`
	Id        string `json:"jti"`
	Audience  string `json:"aud"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	NotBefore int64  `json:"nbf"`
	Subject   string `json:"sub"`
	Scope     string `json:"scope"`
	DeviceId  string `json:"di"`
	UserAgent string `json:"ua"`
	Platform  string `json:"dt"`
	Mth       string `json:"mth"`
	Model     string `json:"md"`
	Ispremium int    `json:"ispre"`
}

var (
	LOCAL_AUTH_TOKEN_ISSUER = "VieOn"
	LOCAL_AUTH_TOKEN_EXPIRE = 2592000
	LOCAL_AUTH_SECRET_KEY   = "fw3g7)36w=_(2ace3s5t&m+t^hd%o*i*&1c9z7h0++7$5__hji"
)

func LocalAuthVerify(tokenString string) (Token, error) {
	var jot Token

	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(LOCAL_AUTH_SECRET_KEY), nil
	})

	if token.Valid {
		if claims, ok := token.Claims.(*Token); ok && token.Valid {
			if strings.HasPrefix(claims.Subject, "anonymous_") == true {
				return jot, errors.New("require_login")
			}
			dataResultStr, _ := json.Marshal(claims)
			json.Unmarshal(dataResultStr, &jot)
			return jot, nil
		} else {
			return jot, errors.New("require_login")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return jot, errors.New("token_invalid")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return jot, errors.New("token_expired")
		} else {
			return jot, errors.New("token_invalid")
		}
	} else {
		return jot, errors.New("token_invalid")
	}

	return jot, nil

	// // Timestamp the beginning.
	// now := time.Now()
	// // Define a signer.
	// hs256 := jwt.NewHS256(LOCAL_AUTH_SECRET_KEY)

	// // First, extract the payload and signature.
	// // This enables unmarshaling the JWT first and
	// // verifying it later or vice versa.
	// payload, sig, err := jwt.Parse(token)
	// if err != nil {
	//     return jot , err
	// }
	// if err = hs256.Verify(payload, sig); err != nil {
	//     return jot , err
	// }
	// if err = jwt.Unmarshal(payload, &jot); err != nil {
	//     return jot , err
	// }

	// // log.Println(jot.Subject)

	// // Validate fields.
	// iatValidator := jwt.IssuedAtValidator(now)
	// expValidator := jwt.ExpirationTimeValidator(now)
	// audValidator := jwt.AudienceValidator("admin")
	// if err = jot.Validate(iatValidator, expValidator, audValidator); err != nil {
	//     switch err {
	//     case jwt.ErrIatValidation:
	//         return jot , errors.New("token_invalid")
	//     case jwt.ErrExpValidation:
	//         return jot , errors.New("token_expired")
	//     // case jwt.ErrAudValidation:
	//     //     return jot , err
	//     }

	//     if strings.HasPrefix(jot.Subject, "anonymous_") == true {
	//         return jot , errors.New("require_login")
	//     }
	// }
	return jot, nil
}

func Run_test_auth(jot Token) {
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaSI6IjMxNTE4ODExNDgiLCJzY29wZSI6ImNtOnJlYWQgY2FzOnJlYWQgYmlsbGluZzpyZWFkIiwiZHQiOiJ3ZWIiLCJuYmYiOjE1NDkyODk2MDgsIm1kIjoiTWFjL2lPUyIsInN1YiI6ImFub255bW91c18zMTUxODgxMTQ4IiwiaXNzIjoiVmllT24iLCJtdGgiOiJhbm9ueW1vdXNfbG9naW4iLCJqdGkiOiJ0a180Yzg4YjE3MzQ2NTg0YzEyYjEyMTMwODM2OTg1NDEwMiIsImV4cCI6MTU1MTg4MTYwOSwiaWF0IjoxNTQ5Mjg5NjA5LCJ1YSI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE0XzApIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS83MS4wLjM1NzguOTggU2FmYXJpLzUzNy4zNiJ9.JSy-5UB33qtJXKRK5iKISUpcEQhfahFpsaOoK98OjgI"
	// jot , err := LocalAuthVerify(token)
	// if err != nil {
	//     log.Println("Err: " , err)
	//     return
	// }

	log.Println("Issuer: ", jot.Issuer)
	log.Println("Subject: ", jot.Subject)
	log.Println("Audience: ", jot.Audience)
	log.Println("ExpirationTime: ", jot.ExpiresAt)
	log.Println("NotBefore: ", jot.NotBefore)
	log.Println("IssuedAt: ", jot.IssuedAt)
	log.Println("ID: ", jot.Id)
	log.Println("Scope: ", jot.Scope)
	log.Println("DeviceId: ", jot.DeviceId)
	log.Println("UserAgent: ", jot.UserAgent)
	log.Println("Platform: ", jot.Platform)
	log.Println("Mth: ", jot.Mth)
	log.Println("Model: ", jot.Model)
}

type ProfileStruct struct {
	Is_premium int
	Status     int
	Avatar     string
	Given_name string
	Mobile     string
	Gender     int
	Email      string
}

var (
	STATUS_BANNED         = 2
	PREFIX_HASH_USER_INFO = "hash_user_info"
)

func CheckAccountBanned(user_id string) bool {
	var Profile ProfileStruct

	//key cache from backuser
	dataRedis, err := mRedis.HGet(PREFIX_HASH_USER_INFO, user_id)
	if err != nil {
		return true
	}
	str, ok := dataRedis.(string)

	if str == "" && !ok {
		return true
	}
	err = json.Unmarshal([]byte(str), &Profile)
	if err == nil && Profile.Status == STATUS_BANNED {
		return false
	}

	return true
}

func GetInfoUserById(user_id string) (ProfileStruct, error) {
	var Profile ProfileStruct

	//key cache from backuser
	dataRedis, err := mRedis.HGet(PREFIX_HASH_USER_INFO, user_id)
	if err != nil {
		return Profile, err
	}
	str, ok := dataRedis.(string)

	if str == "" && !ok {
		return Profile, err
	}
	err = json.Unmarshal([]byte(str), &Profile)
	if err == nil && Profile.Status == 1 {
		return Profile, nil
	}
	return Profile, err
}
