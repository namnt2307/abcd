package giftcode

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedisUSC RedisUSCModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

func Route_UseCode(c *gin.Context) {
	userId := c.GetString("user_id")
	code := c.PostForm("code")

	data, err := Giftcode_UseCode(userId, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), data))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func GiftCodeLG(c *gin.Context) {
	userId := c.GetString("user_id")
	code := c.Param("code")

	data, err := CheckGiftCode(userId, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), data))
		return
	}

	var Resp struct {
		Data interface{} `json:"data" `
	}
	Resp.Data = data

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", Resp))
}


func PromotionCodeDetail(c *gin.Context) {
	userId := c.GetString("user_id")
	code := c.Param("code")
	package_id, _ := StringToInt(c.Param("package_id"))

	data, err := CheckPromotionCode(userId, code, package_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), data))
		return
	}

	var Resp struct {
		Data interface{} `json:"data" `
	}
	Resp.Data = data

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", Resp))
}
