package edit

import (
	"errors"
	"io"
	"os"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/zconst"
	"github.com/m4rc3l05/taggar/internal"
	mediatags "github.com/m4rc3l05/taggar/internal/media_tags"
	"github.com/spf13/cobra"
)

type EditFlags struct {
	Path        *string
	Provider    *string
	Term        *string
	Cover       *string
	Title       *string
	Artist      *string
	Album       *string
	AlbumArtist *string
	Genre       *string
	Year        *string
	Track       *string
	TrackCount  *string
	Disc        *string
	DiscCount   *string
}

var editFlagsSchema = zog.Struct(zog.Shape{
	"Path": zog.Ptr(
		zog.String().
			Min(1).
			TestFunc(func(val *string, ctx internals.Ctx) bool {
				if val == nil {
					return false
				}

				f, err := os.Stat(*val)
				if err != nil {
					return false
				}

				return f.Mode().IsRegular()
			}, zog.Message("Invalid filepath"),
			),
	).NotNil(),
	"Provider": zog.Ptr(zog.String().Optional().OneOf([]string{"itunes"})),
	"Term":     zog.Ptr(zog.String()),
	"Cover": zog.Ptr(
		zog.String().
			TestFunc(func(val *string, ctx internals.Ctx) bool {
				if val == nil {
					return false
				}

				f, err := os.Stat(*val)
				if err != nil {
					return false
				}

				return f.Mode().IsRegular()
			}, zog.Message("Invalid filepath")),
	),
	"Title":       zog.Ptr(zog.String().Optional()),
	"Artist":      zog.Ptr(zog.String().Optional()),
	"Album":       zog.Ptr(zog.String().Optional()),
	"AlbumArtist": zog.Ptr(zog.String().Optional()),
	"Genre":       zog.Ptr(zog.String().Optional()),
	"Year":        zog.Ptr(zog.String().Optional()),
	"Track":       zog.Ptr(zog.String().Optional()),
	"TrackCount":  zog.Ptr(zog.String().Optional()),
	"Disc":        zog.Ptr(zog.String().Optional()),
	"DiscCount":   zog.Ptr(zog.String().Optional()),
}).Test(zog.Test[any]{
	Func: func(val any, ctx internals.Ctx) {
		x := val.(*EditFlags)

		if x.Provider != nil && x.Term == nil {
			ctx.AddIssue(&internals.ZogIssue{
				Code:    zconst.IssueCodeRequired,
				Path:    []string{"Term"},
				Value:   x.Term,
				Message: "Terms must be specified if provider is specified",
			})

			return
		}

		if x.Provider == nil &&
			x.Album == nil &&
			x.AlbumArtist == nil &&
			x.Artist == nil &&
			x.Cover == nil &&
			x.Disc == nil &&
			x.DiscCount == nil &&
			x.Genre == nil &&
			x.Title == nil &&
			x.Track == nil &&
			x.TrackCount == nil &&
			x.Year == nil {
			ctx.AddIssue(&internals.ZogIssue{
				Code:    zconst.IssueCodeRequired,
				Message: "Tags must be specified when no provider is selected",
			})

			return
		}
	},
})

func (ef EditFlags) Validate() zog.ZogIssueList {
	return editFlagsSchema.Validate(&ef)
}

type cmd struct {
	data EditFlags
}

func (c *cmd) Pre(cmd *cobra.Command) error {
	x := EditFlags{}
	if err := Set(cmd, "path", &x.Path); err != nil {
		return err
	}

	if err := Set(cmd, "provider", &x.Provider); err != nil {
		return err
	}
	if err := Set(cmd, "term", &x.Term); err != nil {
		return err
	}

	if err := Set(cmd, "cover", &x.Cover); err != nil {
		return err
	}
	if err := Set(cmd, "title", &x.Title); err != nil {
		return err
	}
	if err := Set(cmd, "artist", &x.Artist); err != nil {
		return err
	}
	if err := Set(cmd, "album", &x.Album); err != nil {
		return err
	}
	if err := Set(cmd, "albumArtist", &x.AlbumArtist); err != nil {
		return err
	}
	if err := Set(cmd, "genre", &x.Genre); err != nil {
		return err
	}
	if err := Set(cmd, "year", &x.Year); err != nil {
		return err
	}
	if err := Set(cmd, "track", &x.Track); err != nil {
		return err
	}
	if err := Set(cmd, "trackCount", &x.TrackCount); err != nil {
		return err
	}
	if err := Set(cmd, "dist", &x.Disc); err != nil {
		return err
	}
	if err := Set(cmd, "distCount", &x.DiscCount); err != nil {
		return err
	}

	if errs := x.Validate(); errs != nil {
		return errors.New(zog.Issues.Prettify(errs))
	}

	c.data = x

	return nil
}

func (c cmd) Run(cmd *cobra.Command) error {
	m := mediatags.MediaTags{
		AlbumArtist: c.data.AlbumArtist,
		Album:       c.data.Album,
		Title:       c.data.Title,
		Year:        c.data.Year,
		Artist:      c.data.Artist,
		Genre:       c.data.Genre,
		Track:       c.data.Track,
		TrackCount:  c.data.TrackCount,
		Disc:        c.data.Disc,
		DiscCount:   c.data.DiscCount,
	}

	if c.data.Cover != nil {
		f, err := os.Open(*c.data.Cover)
		if err != nil {
			return err
		}

		defer func() { internal.Must(f.Close()) }()

		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		m.Cover = &mediatags.MediaTagsCover{
			Data: data,
		}
	}

	tags, err := mediatags.TaglibMediaTagsRepository{}.SetMediaTagsFromPath(*c.data.Path, m)
	if err != nil {
		return err
	}

	return mediatags.DisplayMediaTags(*tags)
}

func Set(cmd *cobra.Command, n string, dest *(*string)) error {
	if cmd.Flags().Changed(n) {
		path, err := cmd.Flags().GetString(n)
		if err != nil {
			return err
		}

		*dest = &path
	}

	return nil
}

func NewCommand() *cobra.Command {
	c := cmd{}
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit audio tags",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.Pre(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(cmd)
		},
	}

	editCmd.Flags().StringP("path", "p", "", "the path to the audio file to view metadata")

	editCmd.Flags().StringP("provider", "s", "", "the provider to use to prefill metadata")
	editCmd.Flags().
		StringP("term", "t", "", "the search term the provider will use to fetch metadata")

	editCmd.Flags().String("cover", "", "the cover")
	editCmd.Flags().String("title", "", "the title")
	editCmd.Flags().String("artist", "", "the artist")
	editCmd.Flags().String("album", "", "the album")
	editCmd.Flags().String("albumArtist", "", "the album artist")
	editCmd.Flags().String("genre", "", "the genre")
	editCmd.Flags().String("year", "", "the year")
	editCmd.Flags().String("track", "", "the track")
	editCmd.Flags().String("trackCount", "", "the track count")
	editCmd.Flags().String("disc", "", "the disc")
	editCmd.Flags().String("discCount", "", "the disc count")

	internal.Must(editCmd.MarkFlagRequired("path"))
	internal.Must(editCmd.MarkFlagFilename("path"))

	return editCmd
}
