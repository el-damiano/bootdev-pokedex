package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/el-damiano/bootdev-pokedex/internal/pokecache"
)

var baseUrl = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: *pokecache.NewCache(cacheInterval),
	}
}

type LocationsResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocations(urlPage *string) (LocationsResp, error) {
	urlEndpoint := "/location-area?offset=0&limit=20"
	url := baseUrl + urlEndpoint
	if urlPage != nil {
		url = *urlPage
	}

	cachedVal, ok := c.cache.Get(url)
	if ok {
		var locationsResp = LocationsResp{}
		err := json.Unmarshal(cachedVal, &locationsResp)
		if err != nil {
			return LocationsResp{}, err
		}
		c.cache.Add(url, cachedVal)
		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationsResp{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationsResp{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationsResp{}, err
	}

	c.cache.Add(url, dat)

	var locationsResp = LocationsResp{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return LocationsResp{}, err
	}

	return locationsResp, nil
}
