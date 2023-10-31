package anizip

import (
	"errors"
	"github.com/goccy/go-json"
	"github.com/seanime-app/seanime-server/internal/result"
	"io"
	"net/http"
	"strconv"
)

type Episode struct {
	TvdbEid               int               `json:"tvdbEid,omitempty"`
	Airdate               string            `json:"airdate,omitempty"`
	SeasonNumber          int               `json:"seasonNumber,omitempty"`
	EpisodeNumber         int               `json:"episodeNumber,omitempty"`
	AbsoluteEpisodeNumber int               `json:"absoluteEpisodeNumber,omitempty"`
	Title                 map[string]string `json:"title,omitempty"`
	Image                 string            `json:"image,omitempty"`
	Summary               string            `json:"summary,omitempty"`
	Overview              string            `json:"overview,omitempty"`
	Runtime               int               `json:"runtime,omitempty"`
	Length                int               `json:"length,omitempty"`
	Episode               string            `json:"episode,omitempty"`
	AnidbEid              int               `json:"anidbEid,omitempty"`
	Rating                string            `json:"rating,omitempty"`
}

type Mappings struct {
	AnimeplanetID string `json:"animeplanet_id,omitempty"`
	KitsuID       int    `json:"kitsu_id,omitempty"`
	MalID         int    `json:"mal_id,omitempty"`
	Type          string `json:"type,omitempty"`
	AnilistID     int    `json:"anilist_id,omitempty"`
	AnisearchID   int    `json:"anisearch_id,omitempty"`
	AnidbID       int    `json:"anidb_id,omitempty"`
	NotifymoeID   string `json:"notifymoe_id,omitempty"`
	LivechartID   int    `json:"livechart_id,omitempty"`
	ThetvdbID     int    `json:"thetvdb_id,omitempty"`
	ImdbID        string `json:"imdb_id,omitempty"`
	ThemoviedbID  string `json:"themoviedb_id,omitempty"`
}

type Media struct {
	Titles       map[string]string  `json:"titles"`
	Episodes     map[string]Episode `json:"episodes"`
	EpisodeCount int                `json:"episodeCount"`
	SpecialCount int                `json:"specialCount"`
	Mappings     *Mappings          `json:"mappings"`
}

func (m *Media) GetTitle() string {
	if m == nil {
		return ""
	}
	if len(m.Titles["en"]) > 0 {
		return m.Titles["en"]
	}
	return m.Titles["ro"]
}

func (m *Media) GetMappings() *Mappings {
	if m != nil {
		return m.Mappings
	}
	return &Mappings{}
}

type Cache struct {
	*result.Cache[string, *Media]
}

func NewCache() *Cache {
	return &Cache{result.NewCache[string, *Media]()}
}

//----------------------------------------------------------------------------------------------------------------------

// FetchAniZipMedia fetches anizip.Media from the AniZip API.
func FetchAniZipMedia(from string, id int) (*Media, error) {
	// Construct the API URL
	apiUrl := "https://api.ani.zip/mappings?" + from + "_id=" + strconv.Itoa(id)

	// Send an HTTP GET request
	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("not found")
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into AniZipData
	var media Media
	if err := json.Unmarshal(responseBody, &media); err != nil {
		return nil, err
	}

	return &media, nil
}

// FetchAniZipMediaC is the same as FetchAniZipMedia but uses a cache.
func FetchAniZipMediaC(from string, id int, cache *Cache) (*Media, error) {

	cacheV, ok := cache.Get(GetCacheKey(from, id))
	if ok {
		return cacheV, nil
	}

	media, err := FetchAniZipMedia(from, id)
	if err != nil {
		return nil, err
	}

	cache.Set(GetCacheKey(from, id), media)

	return media, nil
}

//----------------------------------------------------------------------------------------------------------------------

func GetCacheKey(from string, id int) string {
	return from + strconv.Itoa(id)
}
