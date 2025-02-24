package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var baseUrl = "https://pokeapi.co/api/v2/"

type pokeURL struct {
	hostname string
	endpoint string
	queries  []string
}

func locationURL(urlOld string, offset, limit int) string {
	urlFull := ""
	urlEndpoint := "location-area"
	if urlOld == "" {
		urlQueries := fmt.Sprintf("offset=-%d&limit=%d", offset, limit)
		urlFull = baseUrl + urlEndpoint + "?" + urlQueries
	} else {
		urlFull = urlOld
	}

	urlParsed, err := url.Parse(urlFull)
	if err != nil {
		return ""
	}

	urlNew := baseUrl + urlEndpoint + "?"
	for _, queryRaw := range strings.Split(urlParsed.RawQuery, "&") {
		if queryRaw == "" {
			continue
		}

		queryParsed := strings.Split(queryRaw, "=")
		key, value := queryParsed[0], queryParsed[1]
		if key == "offset" {
			valueNew, err := strconv.Atoi(value)
			if err != nil {
				continue
			}

			valueNew += offset
			if valueNew < 0 {
				valueNew = 0
			}
			urlNew += key + "=" + strconv.Itoa(valueNew)

		} else {
			urlNew += key + "=" + value
		}
		urlNew += "&"
	}

	return urlNew
}
