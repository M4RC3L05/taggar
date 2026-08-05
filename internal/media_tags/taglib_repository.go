package mediatags

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"go.senan.xyz/taglib"
)

var numberAndCountRe = regexp.MustCompile(`^\s*(\d+)\s*(?:/\s*(\d+)\s*)?$`)

type TaglibMediaTagsRepository struct{}

var _ MediaTagsRepository = TaglibMediaTagsRepository{}

func parseNumberAndCount(s string) (num string, count string) {
	matches := numberAndCountRe.FindStringSubmatch(s)

	var nRes string
	var cRes string
	if len(matches) > 1 {
		n := matches[1]
		t := matches[2]

		if n != "" {
			nRes = matches[1]
		}

		if t != "" {
			cRes = matches[2]
		}
	}

	return nRes, cRes
}

func (t TaglibMediaTagsRepository) GetMediaTagsFromPath(path string) (*MediaTags, error) {
	tags, err := taglib.ReadTags(path)
	if err != nil {
		return nil, err
	}

	props, err := taglib.ReadProperties(path)
	if err != nil {
		return nil, err
	}

	mediaTags := MediaTags{}

	coverImgIndex := slices.IndexFunc(props.Images, func(img taglib.ImageDesc) bool {
		return strings.ToLower(img.Type) == "front cover"
	})

	if coverImgIndex == -1 && len(props.Images) > 0 {
		coverImgIndex = 0
	}

	if coverImgIndex != -1 {
		imgBytes, err := taglib.ReadImageOptions(path, coverImgIndex)
		if err != nil {
			return nil, err
		}

		mediaTags.Cover = &MediaTagsCover{
			Data:     imgBytes,
			Mimetype: props.Images[coverImgIndex].MIMEType,
		}
	}

	if val, ok := tags[taglib.AlbumArtist]; ok {
		mediaTags.AlbumArtist = &val
	}

	if val, ok := tags[taglib.Album]; ok {
		mediaTags.Album = &val
	}

	if val, ok := tags[taglib.Title]; ok {
		mediaTags.Title = &val
	}

	if val, ok := tags[taglib.Date]; ok {
		mediaTags.Year = &val
	}

	if val, ok := tags[taglib.Artist]; ok {
		mediaTags.Artist = &val
	}

	if val, ok := tags[taglib.Genre]; ok {
		mediaTags.Genre = &val
	}

	if val, ok := tags[taglib.TrackNumber]; ok {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			n, c := parseNumberAndCount(val[0])

			mediaTags.Track = new([]string{n})
			mediaTags.TrackCount = new([]string{c})
		} else {
			mediaTags.Track = &val
		}
	}

	if val, ok := tags["TRACKTOTAL"]; ok {
		mediaTags.TrackCount = &val
	}

	if val, ok := tags[taglib.DiscNumber]; ok {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			n, c := parseNumberAndCount(val[0])

			mediaTags.Disc = new([]string{n})
			mediaTags.DiscCount = new([]string{c})
		} else {
			mediaTags.Disc = &val
		}
	}

	if val, ok := tags["DISCTOTAL"]; ok {
		mediaTags.DiscCount = &val
	}

	return &mediaTags, nil
}

func (t TaglibMediaTagsRepository) SetMediaTagsFromPath(
	path string,
	tags *MediaTags,
) (*MediaTags, error) {
	if tags == nil {
		fmt.Println("No tags to set")
		return t.GetMediaTagsFromPath(path)
	}

	props, err := taglib.ReadProperties(path)
	if err != nil {
		return nil, err
	}

	tagsToSet := map[string][]string{}

	if tags.AlbumArtist != nil {
		tagsToSet[taglib.AlbumArtist] = *tags.AlbumArtist
	}

	if tags.Album != nil {
		tagsToSet[taglib.Album] = *tags.Album
	}

	if tags.Title != nil {
		tagsToSet[taglib.Title] = *tags.Title
	}

	if tags.Year != nil {
		tagsToSet[taglib.Date] = *tags.Year
	}

	if tags.Artist != nil {
		tagsToSet[taglib.Artist] = *tags.Artist
	}

	if tags.Genre != nil {
		tagsToSet[taglib.Genre] = *tags.Genre
	}

	if tags.Track != nil {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			final := ""
			final += (*tags.Track)[0]

			if tags.TrackCount != nil {
				final += " / " + (*tags.TrackCount)[0]
			}
			tagsToSet[taglib.TrackNumber] = []string{final}
		} else {
			tagsToSet[taglib.TrackNumber] = *tags.Track
		}
	}

	if tags.TrackCount != nil && (props.Format == "flac" || props.Format == "opus" ||
		props.Format == "ogg") {
		tagsToSet["TRACKTOTAL"] = *tags.TrackCount
	}

	if tags.Disc != nil {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			final := ""
			final += (*tags.Disc)[0]

			if tags.DiscCount != nil {
				final += " / " + (*tags.DiscCount)[0]
			}
			tagsToSet[taglib.DiscNumber] = []string{final}
		} else {
			tagsToSet[taglib.DiscNumber] = *tags.Disc
		}
	}

	if tags.DiscCount != nil && (props.Format == "flac" || props.Format == "opus" ||
		props.Format == "ogg") {
		tagsToSet["DISCTOTAL"] = *tags.DiscCount
	}

	err = taglib.WriteTags(path, tagsToSet, 0)
	if err != nil {
		return nil, err
	}

	if tags.Cover != nil {
		err = taglib.WriteImage(path, tags.Cover.Data)
		if err != nil {
			return nil, err
		}
	}

	return t.GetMediaTagsFromPath(path)
}
