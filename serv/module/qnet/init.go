package qnet

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	. "cm-v5/serv/module"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	json             = jsoniter.ConfigCompatibleWithStandardLibrary
	QNET_API         string
	QNET_SECRECT_KEY string
	QNET_OPERATOR_ID int
	QNET_JWT_TOKEN   string
)

func init() {

	QNET_API, _ = CommonConfig.GetString("QNET", "domain_api")
	QNET_SECRECT_KEY, _ = CommonConfig.GetString("QNET", "secret_key")
	QNET_OPERATOR_ID, _ = CommonConfig.GetInt("QNET", "operator_id")
	QNET_JWT_TOKEN, _ = CommonConfig.GetString("QNET", "jwt_token")

}

func GetTokenQNet(c *gin.Context) {
	UserIdBK := c.GetString("user_id")
	UserId := strings.Replace(UserIdBK, "-", "", -1)
	UserId = fmt.Sprint(QNET_OPERATOR_ID) + "-" + UserId

	t := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"timestamp": t,
		"sessionId": UserIdBK+fmt.Sprint(t),
		"userId": UserId,
	})

	ss, err := token.SignedString([]byte(QNET_JWT_TOKEN))

	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Fail"))
		return
	}

	var ResponeData struct {
		SessionId string
		UserId    string
	}
	ResponeData.SessionId = ss
	ResponeData.UserId = UserId

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", ResponeData))
}
