package view

import (
	"errors"
	"os"

	"github.com/m4rc3l05/taggar/internal"
	mediatags "github.com/m4rc3l05/taggar/internal/media_tags"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View audio tags",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("path") {
				return nil
			}

			path, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			f, err := os.Stat(path)
			if err != nil {
				return err
			}

			if !f.Mode().IsRegular() {
				return errors.New("path is not a file")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			tags, err := mediatags.TaglibMediaTagsRepository{}.GetMediaTagsFromPath(path)
			if err != nil {
				return err
			}

			return mediatags.DisplayMediaTags(*tags)
		},
	}

	viewCmd.Flags().StringP("path", "p", "", "the path to the audio file to view metadata")

	internal.Must(viewCmd.MarkFlagRequired("path"))
	internal.Must(viewCmd.MarkFlagFilename("path"))

	return viewCmd
}
