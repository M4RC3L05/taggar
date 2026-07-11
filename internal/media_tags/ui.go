package mediatags

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	"github.com/blacktop/go-termimg"
)

func DisplayMediaTags(tags MediaTags) error {
	padW := 13

	printLine := func(label string, val *string) {
		valStr := "-"
		if val != nil {
			valStr = *val
		}

		dots := strings.Repeat(" ", max(1, padW-len(label)))
		fmt.Printf("%s%s: %s\n", dots, label, valStr)
	}

	if tags.Cover != nil {
		img, _, err := image.Decode(bytes.NewReader(tags.Cover.Data))
		if err != nil {
			return err
		}

		err = termimg.New(img).Size(50, 50).Scale(termimg.ScaleFit).Print()
		if err != nil {
			return err
		}
	} else {
		printLine("Cover", nil)
	}

	fmt.Println("")

	printLine("Title", tags.Title)
	printLine("Artist", tags.Artist)
	printLine("Album", tags.Album)
	printLine("Album artist", tags.AlbumArtist)
	printLine("Genre", tags.Genre)
	printLine("Year", tags.Year)
	printLine("Track", tags.Track)
	printLine("Track count", tags.TrackCount)
	printLine("Disc", tags.Disc)
	printLine("Disc count", tags.DiscCount)

	return nil
}
