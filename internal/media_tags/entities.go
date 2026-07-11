package mediatags

type MediaTagsCover struct {
	Data     []byte
	Mimetype string
}

type MediaTags struct {
	Cover       *MediaTagsCover
	AlbumArtist *string
	Album       *string
	Title       *string
	Year        *string
	Artist      *string
	Genre       *string
	Track       *string
	TrackCount  *string
	Disc        *string
	DiscCount   *string
}
