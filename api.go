// api.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	// updated module path
	"pokedexcli/internal/pokecache"
)

const baseLocationAreaURL = "https://pokeapi.co/api/v2/location-area/"

// ---- LIST (used by map / mapb) ----
type locationAreaListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func fetchLocationAreasCached(cache *pokecache.Cache, url string) (locationAreaListResponse, error) {
	var out locationAreaListResponse
	if url == "" {
		url = baseLocationAreaURL
	}
	// try cache
	if b, ok := cache.Get(url); ok {
		if err := json.Unmarshal(b, &out); err != nil {
			return out, fmt.Errorf("cache unmarshal failed: %w", err)
		}
		return out, nil
	}
	// network
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return out, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	cache.Add(url, body)
	if err := json.Unmarshal(body, &out); err != nil {
		return out, err
	}
	return out, nil
}

// ---- DETAILS (used by explore <area>) ----
type locationAreaDetailResponse struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func fetchLocationAreaDetailCached(cache *pokecache.Cache, area string) (locationAreaDetailResponse, error) {
	var out locationAreaDetailResponse
	if area == "" {
		return out, fmt.Errorf("area name required")
	}
	url := baseLocationAreaURL + area // PokeAPI accepts name or id

	// try cache
	if b, ok := cache.Get(url); ok {
		if err := json.Unmarshal(b, &out); err != nil {
			return out, fmt.Errorf("cache unmarshal failed: %w", err)
		}
		return out, nil
	}
	// network
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return out, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	cache.Add(url, body)
	if err := json.Unmarshal(body, &out); err != nil {
		return out, err
	}
	return out, nil
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func fetchPokemonCached(cache *pokecache.Cache, name string) (Pokemon, error) {
	var out Pokemon
	if name == "" {
		return out, fmt.Errorf("pokemon name required")
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + name

	// cache hit?
	if b, ok := cache.Get(url); ok {
		if err := json.Unmarshal(b, &out); err != nil {
			return out, fmt.Errorf("cache unmarshal failed: %w", err)
		}
		return out, nil
	}

	// network
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return out, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	// store + decode
	cache.Add(url, body)
	if err := json.Unmarshal(body, &out); err != nil {
		return out, err
	}
	return out, nil
}

// shared helper
func strptr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
