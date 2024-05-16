package flickr

import (
	"fmt"
	"net/url"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"sort"
	"strings"
)

type FlickrRequest struct {
	ApiKey      string
	TokenSecret string
	SecretKey   string
	Method      string
	Url         string
	RequestUrl  string
	Args        map[string]string
}

func (request *FlickrRequest) Sign() string {
	delete(request.Args, "oauth_signature")

	s := request.SortArgs()

	signKey := config.APP.FlickrSecret + "&" + request.TokenSecret
	signItem := request.Method + "&" + url.QueryEscape(request.Url) + "&" + url.QueryEscape(s)
	signature := helpers.GetSignature(signItem, signKey)
	fmt.Println("Sign() - Signing", signItem, " SIGN WITH: "+signKey)

	request.Args["oauth_signature"] = url.QueryEscape(signature)

	return signature
}

func (request *FlickrRequest) Do() {
	s := request.SortArgs()
	request.RequestUrl = request.Url + "?" + s
}

func (request *FlickrRequest) SortArgs() string {
	args := request.Args
	sorted_keys := make([]string, len(args))
	// Sort array keys
	i := 0
	for k := range args {
		sorted_keys[i] = k
		i++
	}
	sort.Strings(sorted_keys)

	s := ""
	hasFirst := false
	for _, key := range sorted_keys {
		delimiter := "&"
		if args[key] != "" {
			if !hasFirst {
				delimiter = ""
				hasFirst = true
			}

			keyVal := args[key]
			if key == "title" {
				keyVal = strings.Replace(url.QueryEscape(keyVal), "+", "%20", -1)
			}

			s += delimiter + fmt.Sprintf("%s=%s", key, keyVal)
		}
	}
	return s
}
