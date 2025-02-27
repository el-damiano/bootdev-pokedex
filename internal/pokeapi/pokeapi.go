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

	var locationsResp = LocationsResp{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return LocationsResp{}, err
	}

	c.cache.Add(url, dat)

	return locationsResp, nil
}

type LocationInfoResp struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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
