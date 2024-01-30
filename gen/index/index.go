package index

import (
	"encoding/xml"
	"github.com/gorilla/feeds"
	"time"
)

type Index struct {
	Title         string    `xml:"-"`
	Description   string    `xml:"-"`
	PublishedDate time.Time `xml:"lastmod"`
	Href          string    `xml:"loc"`
	EntryName     string    `xml:"-"`
	Priority      float32   `xml:"priority,omitempty"`
}

const PageSize = 1

var Pages = [][]Index{
	{
		{
			EntryName:     "2024-01-27-webauthn-guide",
			Title:         "A guide to WebAuthn.",
			Description:   "Developing a simple WebAuthn authentication service in Go, as there are few functional implementations of WebAuthn with Go, and only a few existing guides.",
			PublishedDate: time.Unix(1706313600, 0),
			Href:          "/blog/2024-01-27-webauthn-guide",
			Priority:      0.5,
		},
		{
			EntryName:     "2024-01-11-cgo-guide",
			Title:         "Using C libraries in Go with CGO",
			Description:   "Simple guide and recommendations about CGO. For documentation purposes.",
			PublishedDate: time.Unix(1704931200, 0),
			Href:          "/blog/2024-01-11-cgo-guide",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-12-28-architecture-paradigms",
			Title:         "Learn software architecture, paradigms and patterns... even the wrong ones.",
			Description:   "Have you ever wondered whether learning the wrong software architecture is really \"wrong\"? Personally, I've always asked myself this question, and more often than not I've found my answer on the job.",
			PublishedDate: time.Unix(1703721600, 0),
			Href:          "/blog/2023-12-28-architecture-paradigms",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-12-14-about-gentoo-linux",
			Title:         "Gentoo Linux is the best OS for gaming and software development on desktop.",
			Description:   "The review about Gentoo Linux after 1 year of intensive usage in gaming and development: it's the best OS in the world.",
			PublishedDate: time.Unix(1702512000, 0),
			Href:          "/blog/2023-12-14-about-gentoo-linux",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-11-08-go-with-portage-and-crossdev",
			Title:         "Go with Portage and Crossdev, for easy static multi-platform compilation of CGO_ENABLED software.",
			Description:   "Want to statically compile for multi-platform in Go super-easily? Let me introduce Portage, Gentoo's package manager, and Crossdev, Gentoo's solution for cross-compilation.",
			PublishedDate: time.Unix(1699401600, 0),
			Href:          "/blog/2023-11-08-go-with-portage-and-crossdev",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-10-09-understanding-authentication",
			Title:         "Just use OAuth2/OIDC.",
			Description:   "A rant about people implementing their own user database. Also, a guide with detailed implementations on OAuth2/OIDC.",
			PublishedDate: time.Unix(1696809600, 0),
			Href:          "/blog/2023-10-09-understanding-authentication",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-09-22-learn-programming-language",
			Title:         "Learning your first programming language",
			Description:   "About learning your first programming language in 2023. Yes, it's a filler post.",
			PublishedDate: time.Unix(1695340800, 0),
			Href:          "/blog/2023-09-22-learn-programming-language",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-09-16-road-to-replicable-infrastructure",
			Title:         "Road to replicable infrastructure with OverlayFS and dracut live image",
			Description:   "About replicable infrastructure when containerization and virtualization are not allowed.",
			PublishedDate: time.Unix(1694822400, 0),
			Href:          "/blog/2023-09-16-road-to-replicable-infrastructure",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-09-10-developing-blog",
			Title:         "Developing this blog in Go and HTMX",
			Description:   "This article documents about how this blog came to be. From technical choices to deploying this blog.",
			PublishedDate: time.Unix(1694304000, 0),
			Href:          "/blog/2023-09-10-developing-blog",
			Priority:      0.5,
		},
		{
			EntryName:     "2023-09-09-hello-world",
			Title:         "Hello world!",
			Description:   "The very first article. About the motivations of developing this blog from scratch with Go and HTMX, and why I want to write articles on this blog.",
			PublishedDate: time.Unix(1694217600, 0),
			Href:          "/blog/2023-09-09-hello-world",
			Priority:      0.5,
		},
	},
}

func ToSiteMap(ii [][]Index) ([]byte, error) {
	var sitemap = struct {
		XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
		Urls    []Index  `xml:"url"`
	}{
		Urls: []Index{
			{
				Title:         "Marc Nguyen's Blog",
				Description:   "Marc Nguyen's blog is a personal and technical blog about documenting some processes, implementations, etc.",
				PublishedDate: time.Now(),
				Href:          "https://blog.mnguyen.fr",
				Priority:      0.8,
			},
		},
	}

	for _, i := range ii {
		sitemap.Urls = append(sitemap.Urls, i...)
	}

	return xml.MarshalIndent(sitemap, "", "  ")
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
	Updated: time.Unix(1706656793, 0),
	Items: []*feeds.Item{
		{
			Title:       "A guide to WebAuthn.",
			Description: "Developing a simple WebAuthn authentication service in Go, as there are few functional implementations of WebAuthn with Go, and only a few existing guides.",
			Created:     time.Unix(1706313600, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-01-27-webauthn-guide",
			},
		},
		{
			Title:       "Using C libraries in Go with CGO",
			Description: "Simple guide and recommendations about CGO. For documentation purposes.",
			Created:     time.Unix(1704931200, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-01-11-cgo-guide",
			},
		},
		{
			Title:       "Learn software architecture, paradigms and patterns... even the wrong ones.",
			Description: "Have you ever wondered whether learning the wrong software architecture is really \"wrong\"? Personally, I've always asked myself this question, and more often than not I've found my answer on the job.",
			Created:     time.Unix(1703721600, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-12-28-architecture-paradigms",
			},
		},
		{
			Title:       "Gentoo Linux is the best OS for gaming and software development on desktop.",
			Description: "The review about Gentoo Linux after 1 year of intensive usage in gaming and development: it's the best OS in the world.",
			Created:     time.Unix(1702512000, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-12-14-about-gentoo-linux",
			},
		},
		{
			Title:       "Go with Portage and Crossdev, for easy static multi-platform compilation of CGO_ENABLED software.",
			Description: "Want to statically compile for multi-platform in Go super-easily? Let me introduce Portage, Gentoo's package manager, and Crossdev, Gentoo's solution for cross-compilation.",
			Created:     time.Unix(1699401600, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2023-11-08-go-with-portage-and-crossdev",
			},
		},
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
