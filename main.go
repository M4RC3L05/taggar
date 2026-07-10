package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"syscall"

	"github.com/crgimenes/glaze"
	"github.com/m4rc3l05/taggar/internal"
	"go.senan.xyz/taglib"
)

type MediaTags struct {
	Cover       *string `json:"cover,omitempty"`
	AlbumArtist *string `json:"albumArtist,omitempty"`
	Album       *string `json:"album,omitempty"`
	Title       *string `json:"title,omitempty"`
	Year        *string `json:"year,omitempty"`
	Artist      *string `json:"artist,omitempty"`
	Genre       *string `json:"genre,omitempty"`
	Track       *string `json:"track,omitempty"`
	TrackCount  *string `json:"trackCount,omitempty"`
	Disc        *string `json:"disc,omitempty"`
	DiscCount   *string `json:"discCount,omitempty"`
}

//go:embed frontend/\.dist/index.html
var index string

var (
	log              = internal.NewLogger()
	numberAndCountRe = regexp.MustCompile(`^\s*(\d+)\s*(?:/\s*(\d+)\s*)?$`)
)

func parseNumberAndCount(s string) (num *string, count *string) {
	matches := numberAndCountRe.FindStringSubmatch(s)

	var nRes *string
	var cRes *string
	if len(matches) > 1 {
		n := matches[1]
		t := matches[2]

		if n != "" {
			nRes = &matches[1]
		}

		if t != "" {
			cRes = &matches[2]
		}
	}

	return nRes, cRes
}

func init() {
	runtime.LockOSThread()
}

func run(ctx context.Context) int {
	w, err := glaze.New(true)
	if err != nil {
		log.Error("Error creating webview window", "err", err)

		return 1
	}

	defer w.Destroy()

	err = w.Bind("chooseFile", func() (string, error) {
		return w.OpenFile(glaze.FileDialogOptions{
			Title: "Select file to continue",
		})
	})
	if err != nil {
		log.Error("Error binding method to webview", "err", err)

		return 1
	}

	err = w.Bind("chooseDirectory", func() (string, error) {
		return w.OpenDirectory(glaze.FileDialogOptions{
			Title: "Select directory to continue",
		})
	})
	if err != nil {
		log.Error("Error binding method to webview", "err", err)

		return 1
	}

	err = w.Bind("getMediaTags", func(path string) (*MediaTags, error) {
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

			base64Img := fmt.Sprintf(
				"data:%s;base64,%s",
				props.Images[coverImgIndex].MIMEType,
				base64.StdEncoding.EncodeToString(imgBytes),
			)

			mediaTags.Cover = &base64Img
		}

		if val, ok := tags[taglib.AlbumArtist]; ok && len(val) > 0 {
			mediaTags.AlbumArtist = &val[0]
		}

		if val, ok := tags[taglib.Album]; ok && len(val) > 0 {
			mediaTags.Album = &val[0]
		}

		if val, ok := tags[taglib.Title]; ok && len(val) > 0 {
			mediaTags.Title = &val[0]
		}

		if val, ok := tags[taglib.Date]; ok && len(val) > 0 {
			mediaTags.Year = &val[0]
		}

		if val, ok := tags[taglib.Artist]; ok && len(val) > 0 {
			mediaTags.Artist = &val[0]
		}

		if val, ok := tags[taglib.Genre]; ok && len(val) > 0 {
			mediaTags.Genre = &val[0]
		}

		if val, ok := tags[taglib.TrackNumber]; ok && len(val) > 0 {
			if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
				n, c := parseNumberAndCount(val[0])

				mediaTags.Track = n
				mediaTags.TrackCount = c
			} else {
				mediaTags.Track = &val[0]
			}
		}

		if val, ok := tags["TRACKTOTAL"]; ok && len(val) > 0 {
			mediaTags.TrackCount = &val[0]
		}

		if val, ok := tags[taglib.DiscNumber]; ok && len(val) > 0 {
			if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
				n, c := parseNumberAndCount(val[0])

				mediaTags.Disc = n
				mediaTags.DiscCount = c
			} else {
				mediaTags.Disc = &val[0]
			}
		}

		if val, ok := tags["DISCTOTAL"]; ok && len(val) > 0 {
			mediaTags.DiscCount = &val[0]
		}

		return &mediaTags, nil
	})
	if err != nil {
		log.Error("Error binding method to webview", "err", err)

		return 1
	}

	err = w.Bind("setMediaTags", func(path string, tags MediaTags) error {
		fmt.Printf("path: %v\n", path)
		fmt.Printf("tags: %v\n", tags)

		props, err := taglib.ReadProperties(path)
		if err != nil {
			return err
		}

		tagsToSet := map[string][]string{}

		if tags.AlbumArtist != nil {
			tagsToSet[taglib.AlbumArtist] = []string{*tags.AlbumArtist}
		}

		if tags.Album != nil {
			tagsToSet[taglib.Album] = []string{*tags.Album}
		}

		if tags.Title != nil {
			tagsToSet[taglib.Title] = []string{*tags.Title}
		}

		if tags.Year != nil {
			tagsToSet[taglib.Date] = []string{*tags.Year}
		}

		if tags.Artist != nil {
			tagsToSet[taglib.Artist] = []string{*tags.Artist}
		}

		if tags.Genre != nil {
			tagsToSet[taglib.Genre] = []string{*tags.Genre}
		}

		if tags.Track != nil {
			if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
				final := ""
				final += *tags.Track

				if tags.TrackCount != nil {
					final += " / " + *tags.TrackCount
				}
				tagsToSet[taglib.TrackNumber] = []string{final}
			} else {
				tagsToSet[taglib.TrackNumber] = []string{*tags.Track}
			}
		}

		if tags.TrackCount != nil && props.Format != "flac" && props.Format != "opus" &&
			props.Format != "ogg" {
			tagsToSet["TRACKTOTAL"] = []string{*tags.TrackCount}
		}

		if tags.Disc != nil {
			if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
				final := ""
				final += *tags.Disc

				if tags.DiscCount != nil {
					final += " / " + *tags.DiscCount
				}
				tagsToSet[taglib.DiscNumber] = []string{final}
			} else {
				tagsToSet[taglib.DiscNumber] = []string{*tags.Disc}
			}
		}

		if tags.DiscCount != nil && props.Format != "flac" && props.Format != "opus" &&
			props.Format != "ogg" {
			tagsToSet["DISCTOTAL"] = []string{*tags.Disc}
		}

		err = taglib.WriteTags(path, tagsToSet, 0)
		if err != nil {
			return err
		}

		if tags.Cover != nil {
			img, err := base64.StdEncoding.DecodeString(
				(*tags.Cover)[strings.Index(*tags.Cover, "base64,")+len("base64,"):],
			)
			if err != nil {
				return err
			}

			err = taglib.WriteImage(path, img)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Error("Error binding method to webview", "err", err)

		return 1
	}

	go func() {
		<-ctx.Done()

		log.Info("Terminating via os signals")

		w.Dispatch(func() {
			w.Terminate()
		})
	}()

	log.Info("Starting webview")

	w.SetTitle("Taggar")
	w.SetSize(1200, 800, glaze.HintNone)
	w.SetHtml(index)
	w.Run()

	log.Info("Webview closed")

	return 0
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT)

	os.Exit(run(ctx))
}
