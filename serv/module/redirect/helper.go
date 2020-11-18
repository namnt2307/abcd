package redirect

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

const headerField = "X-Rewrite-Original-URI"

type Rule struct {
	Pattern string
	To      string
	Status  int
	*regexp.Regexp
}

type ruleHandle struct {
	pattern string
	to      string
	status  int
}

type httpRequest struct {
	Req        *http.Request
	StatusCode int
}

var regfmt = regexp.MustCompile(`:[^/#?()\.\\]+`)

func NewHttpRequest(req *http.Request) *httpRequest {
	return &httpRequest{
		Req:        req,
		StatusCode: 200,
	}
}

func NewRule(pattern, to string, status int) (*Rule, error) {
	pattern = regfmt.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Rule{
		pattern,
		to,
		status,
		reg,
	}, nil
}

func (r *Rule) Rewrite(hReq *httpRequest) bool {
	req := hReq.Req
	oriPath := req.URL.Path
	if !r.MatchString(oriPath) {
		return false
	}

	to := path.Clean(r.Replace(req.URL))

	u, e := url.ParseRequestURI(to)
	if e != nil {
		return false
	}

	req.Header.Set(headerField, req.URL.RequestURI())

	req.URL.Path = u.Path
	req.URL.RawPath = u.RawPath
	hReq.StatusCode = r.Status
	if u.RawQuery != "" {
		req.URL.RawQuery = u.RawQuery
	}

	return true
}

func (r *Rule) Replace(u *url.URL) string {
	if !hit("\\$|\\:", r.To) {
		return r.To
	}

	uri := u.RequestURI()

	regFrom := regexp.MustCompile(r.Pattern)
	match := regFrom.FindStringSubmatchIndex(uri)

	result := regFrom.ExpandString([]byte(""), r.To, uri, match)

	str := string(result[:])

	if hit("\\:", str) {
		return r.replaceNamedParams(uri, str)
	}

	return str
}

var urlreg = regexp.MustCompile(`:[^/#?()\.\\]+|\(\?P<[a-zA-Z0-9]+>.*\)`)

func (r *Rule) replaceNamedParams(from, to string) string {
	fromMatches := r.FindStringSubmatch(from)

	if len(fromMatches) > 0 {
		for i, name := range r.SubexpNames() {
			if len(name) > 0 {
				to = strings.Replace(to, ":"+name, fromMatches[i], -1)
			}
		}
	}

	return to
}

func NewHandler(rules []ruleHandle) RewriteHandler {
	var h RewriteHandler

	for _, val := range rules {
		r, e := NewRule(val.pattern, val.to, val.status)
		if e != nil {
			panic(e)
		}

		h.rules = append(h.rules, r)
	}

	return h
}

type RewriteHandler struct {
	rules []*Rule
}

func (h *RewriteHandler) ServeHTTP(res http.ResponseWriter, req *httpRequest) {
	for _, r := range h.rules {
		ok := r.Rewrite(req)
		if ok {
			break
		}
	}
	return
}

func hit(pattern, str string) bool {
	r, e := regexp.MatchString(pattern, str)
	if e != nil {
		return false
	}

	return r
}
