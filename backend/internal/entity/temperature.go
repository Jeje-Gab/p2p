package entity

import (
	"errors"
	"fmt"
)

type GetList struct {
	Offset  int     `json:"offset"`
	Limit   int     `json:"limit"`
	Filter  *string `json:"filter"`
	OrderBy string  `json:"order_by"`
}

type BeerFlat struct {
	Id          *int    `json:"id"`
	Style       string  `json:"style"`
	MaxTemp     int     `json:"max_temp"`
	MinTemp     int     `json:"min_temp"`
	AverageTemp float32 `json:"average_temp"`
}

type BeerResp struct {
	Id          *int    `json:"id,omitempty" db:"id"`
	Style       string  `json:"style" db:"style"`
	MaxTemp     int     `json:"max_temp" db:"max_temp"`
	MinTemp     int     `json:"min_temp" db:"min_temp"`
	AverageTemp float32 `json:"average_temp" db:"average_temp"`
}

type BeerRespArray []*BeerResp

// action

type AnalyzeReq struct {
	Temperature int `json:"temperature"`
}

type TrackResp struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Link   string `json:"link"` // link do artista (igual ao exemplo)
}

type PlaylistResp struct {
	Name   string      `json:"name"`
	Tracks []TrackResp `json:"tracks"`
}

type AnalyzeResp struct {
	BeerStyle string        `json:"beerStyle"`
	Playlist  *PlaylistResp `json:"playlist,omitempty"`
}

func (f *BeerFlat) AverageTemperature() error {
	if f == nil {
		return errors.New("BeerFlat is nil")
	}
	if f.MinTemp > f.MaxTemp {
		return fmt.Errorf("min_temp (%d) > max_temp (%d)", f.MinTemp, f.MaxTemp)
	}

	f.AverageTemp = (float32(f.MaxTemp) + float32(f.MinTemp)) / 2.0
	return nil
}
