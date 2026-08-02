package mediatags

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AppleMusicResponse struct {
	ResultCount int                `json:"resultCount"`
	Results     []AppleMusicResult `json:"results"`
}

type AppleMusicResult struct {
	ArtistName           *string    `json:"artistName,omitempty"`
	CollectionName       *string    `json:"collectionName,omitempty"`
	TrackName            *string    `json:"trackName,omitempty"`
	ArtworkURL100        *string    `json:"artworkUrl100,omitempty"`
	ReleaseDate          *time.Time `json:"releaseDate,omitempty"`
	PrimaryGenreName     *string    `json:"primaryGenreName,omitempty"`
	CollectionArtistName *string    `json:"collectionArtistName,omitempty"`
	DiscCount            *int       `json:"discCount,omitempty"`
	DiscNumber           *int       `json:"discNumber,omitempty"`
	TrackCount           *int       `json:"trackCount,omitempty"`
	TrackNumber          *int       `json:"trackNumber,omitempty"`
}

type ITunesMediaTagsProvider struct {
	Id string
}

var _ IProvider = ITunesMediaTagsProvider{}

func IntPtrToStrPtr(pi *int) *string {
	if pi == nil {
		return nil
	}

	return new(strconv.Itoa(*pi))
}

func StrPtrToStrSlicePtr(pi *string) *[]string {
	if pi == nil {
		return nil
	}

	return new([]string{*pi})
}

func (i ITunesMediaTagsProvider) FetchMediaTags() (*MediaTags, error) {
	res, err := http.Get(fmt.Sprintf("https://itunes.apple.com/lookup?id=%s", i.Id))
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var x AppleMusicResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}

	if len(x.Results) <= 0 {
		return nil, fmt.Errorf("could not find tags for \"%s\"", i.Id)
	}

	match := x.Results[0]
	mediaTags := MediaTags{
		AlbumArtist: StrPtrToStrSlicePtr(match.ArtistName),
		Artist:      StrPtrToStrSlicePtr(match.ArtistName),
		Album:       StrPtrToStrSlicePtr(match.CollectionName),
		Title:       StrPtrToStrSlicePtr(match.TrackName),
		Genre:       StrPtrToStrSlicePtr(match.PrimaryGenreName),
		Track:       StrPtrToStrSlicePtr(IntPtrToStrPtr(match.TrackNumber)),
		TrackCount:  StrPtrToStrSlicePtr(IntPtrToStrPtr(match.TrackCount)),
		Disc:        StrPtrToStrSlicePtr(IntPtrToStrPtr(match.DiscNumber)),
		DiscCount:   StrPtrToStrSlicePtr(IntPtrToStrPtr(match.DiscCount)),
	}

	if match.CollectionArtistName != nil {
		mediaTags.AlbumArtist = StrPtrToStrSlicePtr(match.CollectionArtistName)
	}

	if match.ReleaseDate != nil {
		mediaTags.Year = StrPtrToStrSlicePtr(
			new(strconv.Itoa(match.ReleaseDate.In(time.Local).Year())),
		)
	}

	if match.ArtworkURL100 != nil {
		finalUrl := strings.Replace(*match.ArtworkURL100, "100x100bb", "1200x1200bb", 1)
		res, err := http.Get(finalUrl)
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = res.Body.Close()
		}()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		mediaTags.Cover = &MediaTagsCover{
			Data: data,
		}
	}

	return &mediaTags, nil
}
