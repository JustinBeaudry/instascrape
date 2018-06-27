package services

import "instascrape/models"

func Transform(raw *models.InstagramPageData) *models.InstagramFeed {
	// setup feed model
	feed := &models.InstagramFeed{}
	// we only get the first profile, not sure when there'd be multiple
	user := raw.EntryData.ProfilePage[0].Graphql.User

	feed.Id = user.Id
	feed.Bio = user.Biography
	feed.ExternalURL = user.ExternalURL
	feed.FullName = user.FullName
	feed.ProfilePicURL = user.ProfilePicURL
	feed.Posts = make([]models.InstagramPost, len(user.Media.Edges))
	// go through all media nodes
	for i := 0; i < len(user.Media.Edges); i++ {
		node := user.Media.Edges[i].Node
		post := models.InstagramPost{
			URL: node.ImageURL,
			ThumbnailURL: node.ThumbnailURL,
			Date: node.Date,
			Dimensions: models.InstagramImageDimensions{
				Width: node.Dimensions.Width,
				Height: node.Dimensions.Height,
			},
			Caption: node.Caption.Edges[0].Node.Text,
		}
		feed.Posts[i] = post
	}
	return feed
}
