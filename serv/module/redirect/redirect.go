package redirect

import (
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	"net/http"
	"net/http/httptest"
)

type OutputRedirectStruct struct {
	From_url    string `json:"from_url" `
	From_to     string `json:"from_to" `
	Http_status int    `json:"http_status" `
}

func CheckURLDirect(from_url string) (OutputRedirectStruct, error) {
	var OutputRedirect OutputRedirectStruct
	rules, err := GetRules()
	if err != nil {
		return OutputRedirect, err
	}

	handler := NewHandler(rules)
	req, err := http.NewRequest("GET", from_url, nil)
	if err != nil {
		return OutputRedirect, err
	}

	hReq := NewHttpRequest(req)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, hReq)

	OutputRedirect.From_url = from_url
	OutputRedirect.From_to = hReq.Req.URL.Path
	OutputRedirect.Http_status = hReq.StatusCode
	return OutputRedirect, nil
}

func GetRules() ([]ruleHandle, error) {
	var rules []ruleHandle
	var keyCache = "rule_redirect"
	valueCache, err := mRedis.GetString(keyCache)
	if err == nil && valueCache != "" {
		err = json.Unmarshal([]byte(valueCache), &rules)
		if err == nil {
			// return rules, nil
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return rules, err
	}
	defer db_mysql.Close()

	rawData, err := db_mysql.Query(`SELECT pattern_from, pattern_to, http_status FROM redirect_tool WHERE status = 1`)
	if err != nil {
		return rules, err
	}

	for rawData.Next() {
		var rule ruleHandle
		err = rawData.Scan(&rule.pattern, &rule.to, &rule.status)
		if err != nil {
			continue
		}
		rules = append(rules, rule)
	}

	// Write Redis
	dataByte, _ := json.Marshal(rules)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_2_HOURS)
	return rules, nil
}
