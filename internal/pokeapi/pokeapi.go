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

	var locationsResp = LocationsResp{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return LocationsResp{}, err
	}

	c.cache.Add(url, dat)

	return locationsResp, nil
}

func (c *Client) GetLocation(name string) (LocationInfoResp, error) {
	urlEndpoint := "/location-area"
	url := baseUrl + urlEndpoint + "/" + name

	cachedValue, ok := c.cache.Get(url)
	if ok {
		var locationInfo = LocationInfoResp{}
		err := json.Unmarshal(cachedValue, &locationInfo)
		if err != nil {
			return LocationInfoResp{}, err
		}
		c.cache.Add(url, cachedValue)
		return locationInfo, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationInfoResp{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationInfoResp{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationInfoResp{}, err
	}

	var locationInfo = LocationInfoResp{}
	err = json.Unmarshal(data, &locationInfo)
	if err != nil {
		return LocationInfoResp{}, err
	}

	c.cache.Add(url, data)

	return locationInfo, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	urlEndpoint := "/pokemon"
	url := baseUrl + urlEndpoint + "/" + name

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon = Pokemon{}
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
