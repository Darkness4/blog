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
	Href          string    `xml:"-"`
	EntryName     string    `xml:"-"`
	Loc           string    `xml:"loc"`
	Priority      float32   `xml:"priority,omitempty"`
	Tags          []string  `xml:"-"`
}

const PageSize = 1

var Pages = [][]Index{
	{
		{
			EntryName:     "2025-01-25-home-raspi-part-2",
			Title:         "Pushing my Home Raspberry Pi cluster into a production state",
			Description:   "A new year, an overhaul of my home Raspberry Pi cluster.",
			PublishedDate: time.Unix(1737763200, 0),
			Href:          "/blog/2025-01-25-home-raspi-part-2",
			Loc:           "https://blog.mnguyen.fr/blog/2025-01-25-home-raspi-part-2",
			Priority:      0.5,
			Tags: []string{
				"raspberry pi",
				"hpc",
				"kubernetes",
				"cluster",
				"home",
				"monitoring",
				"storage",
				"devops",
			},
		},
		{
			EntryName:     "2024-12-18-k3s-crash-postmortem",
			Title:         "Migration from K3OS to K3s and post-mortem of an incident caused by a corrupted SQLite database.",
			Description:   "My cluster finally crashed! Let's goooooo! A little of context: I'm running a small k3s cluster with 3 Raspberry Pi 4 with a network storage, and I'm using SQLite as a database for my applications.",
			PublishedDate: time.Unix(1734480000, 0),
			Href:          "/blog/2024-12-18-k3s-crash-postmortem",
			Loc:           "https://blog.mnguyen.fr/blog/2024-12-18-k3s-crash-postmortem",
			Priority:      0.5,
			Tags: []string{
				"k3s",
				"k3os",
				"sqlite",
				"raspberry-pi",
				"migration",
				"post-mortem",
				"kubernetes",
				"devops",
			},
		},
		{
			EntryName:     "2024-09-11-fluxcd-argocd-gitops",
			Title:         "A comparison between FluxCD and ArgoCD",
			Description:   "My experience with FluxCD and ArgoCD.",
			PublishedDate: time.Unix(1726012800, 0),
			Href:          "/blog/2024-09-11-fluxcd-argocd-gitops",
			Loc:           "https://blog.mnguyen.fr/blog/2024-09-11-fluxcd-argocd-gitops",
			Priority:      0.5,
			Tags: []string{
				"gitops",
				"fluxcd",
				"argocd",
				"kubernetes",
				"devops",
			},
		},
		{
			EntryName:     "2024-06-23-migrating-cockroachdb",
			Title:         "Migrating from SQLite to CockroachDB",
			Description:   "Small article that review the migration from SQLite to CockroachDB.",
			PublishedDate: time.Unix(1719100800, 0),
			Href:          "/blog/2024-06-23-migrating-cockroachdb",
			Loc:           "https://blog.mnguyen.fr/blog/2024-06-23-migrating-cockroachdb",
			Priority:      0.5,
			Tags: []string{
				"database",
				"sqlite",
				"cockroachdb",
				"devops",
			},
		},
		{
			EntryName:     "2024-06-19-home-raspi",
			Title:         "Presenting my home Raspberry Pi cluster",
			Description:   "Presenting my home Raspberry Pi Kubernetes cluster which is hosting this blog.",
			PublishedDate: time.Unix(1718755200, 0),
			Href:          "/blog/2024-06-19-home-raspi",
			Loc:           "https://blog.mnguyen.fr/blog/2024-06-19-home-raspi",
			Priority:      0.5,
			Tags: []string{
				"raspberry pi",
				"hpc",
				"kubernetes",
				"cluster",
				"home",
				"monitoring",
				"storage",
				"devops",
			},
		},
		{
			EntryName:     "2024-06-18-a-take-zig-c-translate",
			Title:         "A first try on Zig and C interop",
			Description:   "Trying Zig with C libraries for the first time.",
			PublishedDate: time.Unix(1718668800, 0),
			Href:          "/blog/2024-06-18-a-take-zig-c-translate",
			Loc:           "https://blog.mnguyen.fr/blog/2024-06-18-a-take-zig-c-translate",
			Priority:      0.5,
			Tags: []string{
				"zig",
				"c",
				"ffmpeg",
				"av1",
				"ffi",
			},
		},
		{
			EntryName:     "2024-03-17-distributed-systems-in-go",
			Title:         "Fault-Tolerent Distributed Systems with Replicated State Machines in Go",
			Description:   "A simple example of a fault-tolerent distributed system in Go with the Raft consensus algorithm.",
			PublishedDate: time.Unix(1710633600, 0),
			Href:          "/blog/2024-03-17-distributed-systems-in-go",
			Loc:           "https://blog.mnguyen.fr/blog/2024-03-17-distributed-systems-in-go",
			Priority:      0.5,
			Tags: []string{
				"go",
				"distributed systems",
				"fault-tolerent",
				"raft",
				"etcd",
				"bitcoin",
				"ipfs",
			},
		},
		{
			EntryName:     "2024-02-24-gitops-systemd",
			Title:         "GitOps using SystemD",
			Description:   "Pull-based GitOps using SystemD and Git. An alternative to Ansible, Puppet, Chef, and SaltStack.",
			PublishedDate: time.Unix(1708732800, 0),
			Href:          "/blog/2024-02-24-gitops-systemd",
			Loc:           "https://blog.mnguyen.fr/blog/2024-02-24-gitops-systemd",
			Priority:      0.5,
			Tags: []string{
				"devops",
				"gitops",
				"systemd",
				"ansible",
				"puppet",
				"chef",
				"saltstack",
			},
		},
		{
			EntryName:     "2024-01-27-webauthn-guide",
			Title:         "A guide to WebAuthn.",
			Description:   "Developing a simple WebAuthn authentication service in Go, as there are few functional implementations of WebAuthn with Go, and only a few existing guides.",
			PublishedDate: time.Unix(1706313600, 0),
			Href:          "/blog/2024-01-27-webauthn-guide",
			Loc:           "https://blog.mnguyen.fr/blog/2024-01-27-webauthn-guide",
			Priority:      0.5,
			Tags: []string{
				"go",
				"webauthn",
				"authentication",
				"security",
			},
		},
		{
			EntryName:     "2024-01-11-cgo-guide",
			Title:         "Using C libraries in Go with CGO",
			Description:   "Simple guide and recommendations about CGO. For documentation purposes.",
			PublishedDate: time.Unix(1704931200, 0),
			Href:          "/blog/2024-01-11-cgo-guide",
			Loc:           "https://blog.mnguyen.fr/blog/2024-01-11-cgo-guide",
			Priority:      0.5,
			Tags: []string{
				"go",
				"cgo",
				"ffi",
				"c",
			},
		},
		{
			EntryName:     "2023-12-28-architecture-paradigms",
			Title:         "Learn software architecture, paradigms and patterns... even the wrong ones.",
			Description:   "Have you ever wondered whether learning the wrong software architecture is really \"wrong\"? Personally, I've always asked myself this question, and more often than not I've found my answer on the job.",
			PublishedDate: time.Unix(1703721600, 0),
			Href:          "/blog/2023-12-28-architecture-paradigms",
			Loc:           "https://blog.mnguyen.fr/blog/2023-12-28-architecture-paradigms",
			Priority:      0.5,
			Tags: []string{
				"software architecture",
				"paradigms",
				"patterns",
				"programming",
			},
		},
		{
			EntryName:     "2023-12-14-about-gentoo-linux",
			Title:         "Gentoo Linux is the best OS for gaming and software development on desktop.",
			Description:   "The review about Gentoo Linux after 1 year of intensive usage in gaming and development: it's the best OS in the world.",
			PublishedDate: time.Unix(1702512000, 0),
			Href:          "/blog/2023-12-14-about-gentoo-linux",
			Loc:           "https://blog.mnguyen.fr/blog/2023-12-14-about-gentoo-linux",
			Priority:      0.5,
			Tags: []string{
				"gentoo",
				"linux",
				"review",
			},
		},
		{
			EntryName:     "2023-11-08-go-with-portage-and-crossdev",
			Title:         "Go with Portage and Crossdev, for easy static multi-platform compilation of CGO_ENABLED software.",
			Description:   "Want to statically compile for multi-platform in Go super-easily? Let me introduce Portage, Gentoo's package manager, and Crossdev, Gentoo's solution for cross-compilation.",
			PublishedDate: time.Unix(1699401600, 0),
			Href:          "/blog/2023-11-08-go-with-portage-and-crossdev",
			Loc:           "https://blog.mnguyen.fr/blog/2023-11-08-go-with-portage-and-crossdev",
			Priority:      0.5,
			Tags: []string{
				"go",
				"cross-compilation",
				"portage",
				"crossdev",
				"gentoo",
				"static-compilation",
				"cgo",
				"docker",
				"multi-arch",
			},
		},
		{
			EntryName:     "2023-10-09-understanding-authentication",
			Title:         "Just use OAuth2/OIDC.",
			Description:   "A rant about people implementing their own user database. Also, a guide with detailed implementations on OAuth2/OIDC.",
			PublishedDate: time.Unix(1696809600, 0),
			Href:          "/blog/2023-10-09-understanding-authentication",
			Loc:           "https://blog.mnguyen.fr/blog/2023-10-09-understanding-authentication",
			Priority:      0.5,
			Tags: []string{
				"security",
				"authentication",
				"oauth2",
				"oidc",
				"dex",
				"389ds",
				"ldap",
			},
		},
		{
			EntryName:     "2023-09-22-learn-programming-language",
			Title:         "Learning your first programming language",
			Description:   "About learning your first programming language in 2023. Yes, it's a filler post.",
			PublishedDate: time.Unix(1695340800, 0),
			Href:          "/blog/2023-09-22-learn-programming-language",
			Loc:           "https://blog.mnguyen.fr/blog/2023-09-22-learn-programming-language",
			Priority:      0.5,
			Tags: []string{
				"programming",
				"go",
				"rust",
				"zig",
				"c",
				"c++",
				"python",
				"javascript",
				"typescript",
				"kotlin",
				"objective-c",
				"lua",
				"ruby",
			},
		},
		{
			EntryName:     "2023-09-16-road-to-replicable-infrastructure",
			Title:         "Road to replicable infrastructure with OverlayFS and dracut live image",
			Description:   "About replicable infrastructure when containerization and virtualization are not allowed.",
			PublishedDate: time.Unix(1694822400, 0),
			Href:          "/blog/2023-09-16-road-to-replicable-infrastructure",
			Loc:           "https://blog.mnguyen.fr/blog/2023-09-16-road-to-replicable-infrastructure",
			Priority:      0.5,
			Tags: []string{
				"devops",
				"linux",
				"infrastructure",
				"dracut",
				"squashfs",
				"overlayfs",
				"pxe",
				"gitops",
			},
		},
		{
			EntryName:     "2023-09-10-developing-blog",
			Title:         "Developing this blog in Go and HTMX",
			Description:   "This article documents about how this blog came to be. From technical choices to deploying this blog.",
			PublishedDate: time.Unix(1694304000, 0),
			Href:          "/blog/2023-09-10-developing-blog",
			Loc:           "https://blog.mnguyen.fr/blog/2023-09-10-developing-blog",
			Priority:      0.5,
			Tags: []string{
				"blog",
				"go",
				"htmx",
				"raspberry-pi",
				"kubernetes",
			},
		},
		{
			EntryName:     "2023-09-09-hello-world",
			Title:         "Hello world!",
			Description:   "The very first article. About the motivations of developing this blog from scratch with Go and HTMX, and why I want to write articles on this blog.",
			PublishedDate: time.Unix(1694217600, 0),
			Href:          "/blog/2023-09-09-hello-world",
			Loc:           "https://blog.mnguyen.fr/blog/2023-09-09-hello-world",
			Priority:      0.5,
			Tags: []string{
				"go",
				"htmx",
			},
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
				Loc:           "https://blog.mnguyen.fr",
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
	Updated: time.Unix(1734484547, 0),
	Items: []*feeds.Item{
		{
			Title:       "Migration from K3OS to K3s and post-mortem of an incident caused by a corrupted SQLite database.",
			Description: "My cluster finally crashed! Let's goooooo! A little of context: I'm running a small k3s cluster with 3 Raspberry Pi 4 with a network storage, and I'm using SQLite as a database for my applications.",
			Created:     time.Unix(1734480000, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-12-18-k3s-crash-postmortem",
			},
		},
		{
			Title:       "A comparison between FluxCD and ArgoCD",
			Description: "My experience with FluxCD and ArgoCD.",
			Created:     time.Unix(1726012800, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-09-11-fluxcd-argocd-gitops",
			},
		},
		{
			Title:       "Migrating from SQLite to CockroachDB",
			Description: "Small article that review the migration from SQLite to CockroachDB.",
			Created:     time.Unix(1719100800, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-06-23-migrating-cockroachdb",
			},
		},
		{
			Title:       "Presenting my home Raspberry Pi cluster",
			Description: "Presenting my home Raspberry Pi Kubernetes cluster which is hosting this blog.",
			Created:     time.Unix(1718755200, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-06-19-home-raspi",
			},
		},
		{
			Title:       "A first try on Zig and C interop",
			Description: "Trying Zig with C libraries for the first time.",
			Created:     time.Unix(1718668800, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-06-18-a-take-zig-c-translate",
			},
		},
		{
			Title:       "Fault-Tolerent Distributed Systems with Replicated State Machines in Go",
			Description: "A simple example of a fault-tolerent distributed system in Go with the Raft consensus algorithm.",
			Created:     time.Unix(1710633600, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-03-17-distributed-systems-in-go",
			},
		},
		{
			Title:       "GitOps using SystemD",
			Description: "Pull-based GitOps using SystemD and Git. An alternative to Ansible, Puppet, Chef, and SaltStack.",
			Created:     time.Unix(1708732800, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2024-02-24-gitops-systemd",
			},
		},
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
