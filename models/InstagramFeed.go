package models

type InstagramImageDimensions struct {
	Width int
	Height int
}

type InstagramPost struct {
	URL string
	ThumbnailURL string
	Date int
	Dimensions InstagramImageDimensions
	Caption string
}

type InstagramFeed struct {
	Id string
	Bio string
	ExternalURL string
	FullName string
	ProfilePicURL string
	Posts []InstagramPost
}
