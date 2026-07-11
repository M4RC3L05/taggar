package mediatags

type IProvider interface {
	FetchMediaTags(term string) (*MediaTags, error)
}
