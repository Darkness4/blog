package index

import (
	"time"

	"github.com/gorilla/feeds"
)

type Index struct {
	Title         string
	Description   string
	PublishedDate time.Time
	Href          string
	EntryName     string
}

const PageSize = 1

var Pages = [][]Index{
	{
		{
			EntryName:     "2023-10-09-understanding-authentication",
			Title:         "Just use OAuth2/OIDC.",
			Description:   "A rant about people implementing their own user database. Also, a guide with detailed implementations on OAuth2/OIDC.",
			PublishedDate: time.Unix(1696809600, 0),
			Href:          "/blog/2023-10-09-understanding-authentication",
		},
		{
			EntryName:     "2023-09-22-learn-programming-language",
			Title:         "Learning your first programming language",
			Description:   "About learning your first programming language in 2023. Yes, it's a filler post.",
			PublishedDate: time.Unix(1695340800, 0),
			Href:          "/blog/2023-09-22-learn-programming-language",
		},
		{
			EntryName:     "2023-09-16-road-to-replicable-infrastructure",
			Title:         "Road to replicable infrastructure with OverlayFS and dracut live image",
			Description:   "About replicable infrastructure when containerization and virtualization are not allowed.",
			PublishedDate: time.Unix(1694822400, 0),
			Href:          "/blog/2023-09-16-road-to-replicable-infrastructure",
		},
		{
			EntryName:     "2023-09-10-developing-blog",
			Title:         "Developing this blog in Go and HTMX",
			Description:   "This article documents about how this blog came to be. From technical choices to deploying this blog.",
			PublishedDate: time.Unix(1694304000, 0),
			Href:          "/blog/2023-09-10-developing-blog",
		},
		{
			EntryName:     "2023-09-09-hello-world",
			Title:         "Hello world!",
			Description:   "The very first article. About the motivations of developing this blog from scratch with Go and HTMX, and why I want to write articles on this blog.",
			PublishedDate: time.Unix(1694217600, 0),
			Href:          "/blog/2023-09-09-hello-world",
		},
	},
}

var Feed = &feeds.Feed{
	Title: "Marc Nguyen's Blog",
	Link: &feeds.Link{
		Href: "https://blog.mnguyen.fr",
	},
	Description: "Marc Nguyen's blog is a personal and technical blog about documenting some processes, implementations, etc.",
	Author: &feeds.Author{
		Name:  "Marc Nguyen",
		Email: "nguyen_marc@live.fr",
	},
	Created: time.Unix(1694131200, 0),
	Updated: time.Unix(1698584941, 0),
	Items: []*feeds.Item{
		{
			Title:       "Just use OAuth2/OIDC.",
			Description: "A rant about people implementing their own user database. Also, a guide with detailed implementations on OAuth2/OIDC.",
			Created:     time.Unix(1696809600, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-10-09-understanding-authentication",
			},
		},
		{
			Title:       "Learning your first programming language",
			Description: "About learning your first programming language in 2023. Yes, it's a filler post.",
			Created:     time.Unix(1695340800, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-09-22-learn-programming-language",
			},
		},
		{
			Title:       "Road to replicable infrastructure with OverlayFS and dracut live image",
			Description: "About replicable infrastructure when containerization and virtualization are not allowed.",
			Created:     time.Unix(1694822400, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-09-16-road-to-replicable-infrastructure",
			},
		},
		{
			Title:       "Developing this blog in Go and HTMX",
			Description: "This article documents about how this blog came to be. From technical choices to deploying this blog.",
			Created:     time.Unix(1694304000, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-09-10-developing-blog",
			},
		},
		{
			Title:       "Hello world!",
			Description: "The very first article. About the motivations of developing this blog from scratch with Go and HTMX, and why I want to write articles on this blog.",
			Created:     time.Unix(1694217600, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-09-09-hello-world",
			},
		},
	},
}
