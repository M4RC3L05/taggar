package mediatags

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	TrackViewUrl         *string    `json:"trackViewUrl,omitempty"`
}

type AppleMusicSchemaSong struct {
	Audio *struct {
		ByArtist *[]struct {
			Name *string `json:"name,omitempty"`
		} `json:"byArtist,omitempty"`
		InAlbum *struct {
			ByArtist *[]struct {
				Name *string `json:"name,omitempty"`
			} `json:"byArtist,omitempty"`
		} `json:"inAlbum,omitempty"`
	} `json:"audio,omitempty"`
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

var songJsonLDRe = regexp.MustCompile(
	`(?is)<script\s+(?:[^>]*\s+)?id\s*=\s*["\']?schema:song["\']?(?:\s+[^>]*?)?\s+type\s*=\s*["\']?application/ld\+json["\']?(?:\s+[^>]*)?>(.+?)</script>`,
)

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

	// properly fetch all artists of this track
	if match.TrackViewUrl != nil {
		res, err := http.Get(*match.TrackViewUrl)
		defer func() {
			_ = res.Body.Close()
		}()

		if err == nil {
			html, err := io.ReadAll(res.Body)

			if err == nil {
				matches := songJsonLDRe.FindStringSubmatch(string(html))

				if len(matches) > 1 {
					jsonContent := matches[1]

					var x AppleMusicSchemaSong
					err = json.Unmarshal([]byte(jsonContent), &x)

					if err == nil {
						if x.Audio != nil && x.Audio.ByArtist != nil && len(*x.Audio.ByArtist) > 0 {
							finalArtist := []string{}

							for _, artist := range *x.Audio.ByArtist {
								if artist.Name != nil {
									finalArtist = append(finalArtist, *artist.Name)
								}
							}

							if len(finalArtist) > 0 {
								mediaTags.Artist = &finalArtist
							}
						}

						if x.Audio != nil && x.Audio.InAlbum != nil &&
							x.Audio.InAlbum.ByArtist != nil &&
							len(*x.Audio.InAlbum.ByArtist) > 0 {
							finalArtist := []string{}

							for _, artist := range *x.Audio.InAlbum.ByArtist {
								if artist.Name != nil {
									finalArtist = append(finalArtist, *artist.Name)
								}
							}

							if len(finalArtist) > 0 {
								mediaTags.AlbumArtist = &finalArtist
							}
						}
					}
				}
			}
		}
	}

	return &mediaTags, nil
}
