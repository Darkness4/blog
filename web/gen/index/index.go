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
	Hierarchy     []Header  `xml:"-"`
}

// Header represents a single header in the hierarchy
type Header struct {
	Level    int
	Text     string
	Anchor   string
	Content  string
	Children []Header
}

const PageSize = 1

var Pages = [][]Index{
	{
		{
			EntryName:     "2025-07-24-fpv-drone",
			Title:         "I'm back! And I'm now flying FPV drones!",
			Description:   "As an engineer, how I got started with FPV drones.",
			PublishedDate: time.Unix(1753315200, 0),
			Href:          "/blog/2025-07-24-fpv-drone",
			Loc:           "https://blog.mnguyen.fr/blog/2025-07-24-fpv-drone",
			Priority:      0.5,
			Tags: []string{
				"drone",
				"fpv",
				"programming",
				"electronics",
				"soldering",
			},
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "A Bit of My Story",
					Anchor:  "a-bit-of-my-story",
					Content: "",
				},

				{
					Level:   2,
					Text:    "From repairing mouses to building drones",
					Anchor:  "from-repairing-mouses-to-building-drones",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Many failures",
							Anchor:  "many-failures",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Repaired a Nintendo DS lite",
							Anchor:  "repaired-a-nintendo-ds-lite",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Repairing a Nintendo Switch Pro Controller",
							Anchor:  "repairing-a-nintendo-switch-pro-controller",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Making profitable my soldering hardware",
							Anchor:  "making-profitable-my-soldering-hardware",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Time to build an FPV drone",
					Anchor:  "time-to-build-an-fpv-drone",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Which drone?",
							Anchor:  "which-drone",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Shopping list",
							Anchor:  "shopping-list",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Building the drone",
							Anchor:  "building-the-drone",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Configuring the firmware",
							Anchor:  "configuring-the-firmware",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Test it",
							Anchor:  "test-it",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "How to pilot a drone",
					Anchor:  "how-to-pilot-a-drone",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Credits",
					Anchor:  "credits",
					Content: "",
				},
			},
		},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Hardware overhaul",
					Anchor:  "hardware-overhaul",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Replacing the router",
							Anchor:  "replacing-the-router",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Adding the storage node",
							Anchor:  "adding-the-storage-node",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Replacing k3os with simple RaspiOS with k3s",
							Anchor:  "replacing-k3os-with-simple-raspios-with-k3s",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Software overhaul",
					Anchor:  "software-overhaul",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Ansible",
							Anchor:  "ansible",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Migrating CockroachDB to PostgreSQL",
							Anchor:  "migrating-cockroachdb-to-postgresql",
							Content: "",
						},

						{
							Level:   3,
							Text:    "New services, and death to some",
							Anchor:  "new-services-and-death-to-some",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "FluxCD",
									Anchor:  "fluxcd",
									Content: "",
								},

								{
									Level:   4,
									Text:    "LLDAP",
									Anchor:  "lldap",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Authelia",
									Anchor:  "authelia",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Crowdsec",
									Anchor:  "crowdsec",
									Content: "",
								},

								{
									Level:   4,
									Text:    "VictoriaLogs and Vectors",
									Anchor:  "victorialogs-and-vectors",
									Content: "",
								},

								{
									Level:   4,
									Text:    "ArchiSteamFarm",
									Anchor:  "archisteamfarm",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Tried self-hosting a mail server with Maddy, using Scaleway Transactional mail instead",
									Anchor:  "tried-self-hosting-a-mail-server-with-maddy-using-scaleway-transactional-mail-instead",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Backups and AWS mountpoint S3 CSI Driver",
									Anchor:  "backups-and-aws-mountpoint-s3-csi-driver",
									Content: "",
								},
							},
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Migrating from K3os to K3s",
					Anchor:  "migrating-from-k3os-to-k3s",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Move PVCs to network storage",
							Anchor:  "move-pvcs-to-network-storage",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Migrating workers",
							Anchor:  "migrating-workers",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Migrating the master",
							Anchor:  "migrating-the-master",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Post-mortem of the crash",
					Anchor:  "post-mortem-of-the-crash",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Summary",
							Anchor:  "summary",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Leadup",
							Anchor:  "leadup",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Detection",
							Anchor:  "detection",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Recovery",
							Anchor:  "recovery",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Permanent fix",
					Anchor:  "permanent-fix",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Lesson learned and corrective actions",
					Anchor:  "lesson-learned-and-corrective-actions",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "What is GitOps?",
					Anchor:  "what-is-gitops",
					Content: "",
				},

				{
					Level:   2,
					Text:    "The comparison",
					Anchor:  "the-comparison",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "The dashboard",
							Anchor:  "the-dashboard",
							Content: "",
						},

						{
							Level:   3,
							Text:    "The setup",
							Anchor:  "the-setup",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Resource consumption",
							Anchor:  "resource-consumption",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Multi-cluster",
							Anchor:  "multi-cluster",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Dynamic applications",
							Anchor:  "dynamic-applications",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Multi-tenancy",
							Anchor:  "multi-tenancy",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Hooks",
							Anchor:  "hooks",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Post-rendering",
							Anchor:  "post-rendering",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Issues",
							Anchor:  "issues",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "What is CockroachDB?",
					Anchor:  "what-is-cockroachdb",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Missing features and critical differences compared to PostgreSQL",
					Anchor:  "missing-features-and-critical-differences-compared-to-postgresql",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Deployment",
					Anchor:  "deployment",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Migration",
					Anchor:  "migration",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Adding monitoring",
					Anchor:  "adding-monitoring",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Adding backup",
					Anchor:  "adding-backup",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Last improvements",
					Anchor:  "last-improvements",
					Content: "",
				},

				{
					Level:   2,
					Text:    "What has been migrated? What couldn't be migrated?",
					Anchor:  "what-has-been-migrated-what-couldnt-be-migrated",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Building the cluster",
					Anchor:  "building-the-cluster",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Storage",
							Anchor:  "storage",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Power supply",
							Anchor:  "power-supply",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Router",
							Anchor:  "router",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Nodes",
							Anchor:  "nodes",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Setting up Kubernetes and the services",
					Anchor:  "setting-up-kubernetes-and-the-services",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Core services",
							Anchor:  "core-services",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Monitoring services",
							Anchor:  "monitoring-services",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Application Services",
							Anchor:  "application-services",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Storage services",
							Anchor:  "storage-services",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "The concept",
					Anchor:  "the-concept",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Zig and tricks",
					Anchor:  "zig-and-tricks",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Memory allocators",
							Anchor:  "memory-allocators",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Easy memory management with ",
							Anchor:  "easy-memory-management-with-defer",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Error handling",
							Anchor:  "error-handling",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Zig's quick error handling",
							Anchor:  "zigs-quick-error-handling",
							Content: "",
						},

						{
							Level:   3,
							Text:    "C interop",
							Anchor:  "c-interop",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Limitations of the C interop",
							Anchor:  "limitations-of-the-c-interop",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Visibility",
							Anchor:  "visibility",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Syntax",
							Anchor:  "syntax",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Developing the AV1 transcoder",
					Anchor:  "developing-the-av1-transcoder",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Developing the Remuxer",
							Anchor:  "developing-the-remuxer",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Developing the Transcoder",
							Anchor:  "developing-the-transcoder",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Last part: the build system",
					Anchor:  "last-part-the-build-system",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "What is a distributed system?",
					Anchor:  "what-is-a-distributed-system",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Examples of stateful distributed systems",
					Anchor:  "examples-of-stateful-distributed-systems",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "ETCD",
							Anchor:  "etcd",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Bitcoin",
							Anchor:  "bitcoin",
							Content: "",
						},

						{
							Level:   3,
							Text:    "IPFS",
							Anchor:  "ipfs",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Summary",
							Anchor:  "summary",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "The example: A distributed key-value store with Raft",
					Anchor:  "the-example-a-distributed-key-value-store-with-raft",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "The objective",
							Anchor:  "the-objective",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Understanding Raft at its core",
							Anchor:  "understanding-raft-at-its-core",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Logs and state machines",
							Anchor:  "logs-and-state-machines",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Bootstrapping the project",
							Anchor:  "bootstrapping-the-project",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Implementing the key-value store",
							Anchor:  "implementing-the-key-value-store",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "The state machine",
									Anchor:  "the-state-machine",
									Content: "",
								},

								{
									Level:   4,
									Text:    "The implementation",
									Anchor:  "the-implementation",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Using Raft to distribute commands",
							Anchor:  "using-raft-to-distribute-commands",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Understanding Raft's lifecycle",
									Anchor:  "understanding-rafts-lifecycle",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Understanding Raft's RPCs and Term",
									Anchor:  "understanding-rafts-rpcs-and-term",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Implementing the Finite State Machine for Raft",
									Anchor:  "implementing-the-finite-state-machine-for-raft",
									Content: "",
									Children: []Header{

										{
											Level:   5,
											Text:    "Define commands for Raft",
											Anchor:  "define-commands-for-raft",
											Content: "",
										},

										{
											Level:   5,
											Text:    "Implementing the Finite State Machine",
											Anchor:  "implementing-the-finite-state-machine",
											Content: "",
										},
									},
								},

								{
									Level:   4,
									Text:    "The \"crash\" recovery: snapshots and restoring logs",
									Anchor:  "the-crash-recovery-snapshots-and-restoring-logs",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Preparing the storage for Raft",
									Anchor:  "preparing-the-storage-for-raft",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Replacing the network layer with a mutual TLS transport",
									Anchor:  "replacing-the-network-layer-with-a-mutual-tls-transport",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Adding the \"Join\" and \"Leave\" methods",
									Anchor:  "adding-the-join-and-leave-methods",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Sending commands",
									Anchor:  "sending-commands",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Adding an API for interactivity",
							Anchor:  "adding-an-api-for-interactivity",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Implementing the \"main\" function of the server",
							Anchor:  "implementing-the-main-function-of-the-server",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Usage",
									Anchor:  "usage",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Implementing the bootstrap function",
									Anchor:  "implementing-the-bootstrap-function",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Implementing helper functions for TLS",
									Anchor:  "implementing-helper-functions-for-tls",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Assembling everything",
									Anchor:  "assembling-everything",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Smoke tests with Vagrant and K0s",
									Anchor:  "smoke-tests-with-vagrant-and-k0s",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Implementing the client",
							Anchor:  "implementing-the-client",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Usage",
									Anchor:  "usage-1",
									Content: "",
								},

								{
									Level:   4,
									Text:    "About client-side load balancing",
									Anchor:  "about-client-side-load-balancing",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Implementing server-side request forwarding",
									Anchor:  "implementing-server-side-request-forwarding",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Implementing the membership service",
									Anchor:  "implementing-the-membership-service",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Implementing the \"main\" function of the client",
									Anchor:  "implementing-the-main-function-of-the-client",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Exposing the nodes to the external network",
									Anchor:  "exposing-the-nodes-to-the-external-network",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Setting mutual TLS on the client-side",
									Anchor:  "setting-mutual-tls-on-the-client-side",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Testing failures",
							Anchor:  "testing-failures",
							Content: "",
						},

						{
							Level:   3,
							Text:    "About Serf, the service discovery system",
							Anchor:  "about-serf-the-service-discovery-system",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "The recipe for pull-based GitOps",
					Anchor:  "the-recipe-for-pull-based-gitops",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Understanding the init system and service management",
					Anchor:  "understanding-the-init-system-and-service-management",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Implementing the script",
					Anchor:  "implementing-the-script",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Configuring the SystemD service",
					Anchor:  "configuring-the-systemd-service",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Building the OS image",
					Anchor:  "building-the-os-image",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Testing",
					Anchor:  "testing",
					Content: "",
				},

				{
					Level:   2,
					Text:    "(Optional) Make the node reboot",
					Anchor:  "optional-make-the-node-reboot",
					Content: "",
				},

				{
					Level:   2,
					Text:    "(Optional) Fetch the status and logs",
					Anchor:  "optional-fetch-the-status-and-logs",
					Content: "",
				},

				{
					Level:   2,
					Text:    "(Optional) Using a webhook to trigger the service",
					Anchor:  "optional-using-a-webhook-to-trigger-the-service",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion and Discussion",
					Anchor:  "conclusion-and-discussion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Quick remainder about authentication, identity and session",
					Anchor:  "quick-remainder-about-authentication-identity-and-session",
					Content: "",
				},

				{
					Level:   2,
					Text:    "What is WebAuthn?",
					Anchor:  "what-is-webauthn",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Quick definition",
							Anchor:  "quick-definition",
							Content: "",
						},

						{
							Level:   3,
							Text:    "The authentication flow",
							Anchor:  "the-authentication-flow",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Implementation",
					Anchor:  "implementation",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "The backend",
							Anchor:  "the-backend",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "The user",
									Anchor:  "the-user",
									Content: "",
								},

								{
									Level:   4,
									Text:    "The login/register session store",
									Anchor:  "the-loginregister-session-store",
									Content: "",
								},

								{
									Level:   4,
									Text:    "The JWT session generator and validator",
									Anchor:  "the-jwt-session-generator-and-validator",
									Content: "",
								},

								{
									Level:   4,
									Text:    "The HTTP handlers for WebAuthn",
									Anchor:  "the-http-handlers-for-webauthn",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Set up the HTTP server",
									Anchor:  "set-up-the-http-server",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "The frontend",
							Anchor:  "the-frontend",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Preface",
									Anchor:  "preface",
									Content: "",
								},

								{
									Level:   4,
									Text:    "The login/register page",
									Anchor:  "the-loginregister-page",
									Content: "",
								},

								{
									Level:   4,
									Text:    "The protected page",
									Anchor:  "the-protected-page",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Testing",
							Anchor:  "testing",
							Content: "",
						},

						{
							Level:   3,
							Text:    "CSRF and Templates",
							Anchor:  "csrf-and-templates",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Initial setup for templates",
									Anchor:  "initial-setup-for-templates",
									Content: "",
								},

								{
									Level:   4,
									Text:    "CSRF protection",
									Anchor:  "csrf-protection",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Adding/deleting devices for logged users",
							Anchor:  "addingdeleting-devices-for-logged-users",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion and thoughts",
					Anchor:  "conclusion-and-thoughts",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Quick reminder about C compilation",
					Anchor:  "quick-reminder-about-c-compilation",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Why use CGO instead of pure Go ?",
					Anchor:  "why-use-cgo-instead-of-pure-go-",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Foreign Function Interface (FFI)",
					Anchor:  "foreign-function-interface-ffi",
					Content: "",
				},

				{
					Level:   2,
					Text:    "About Go libraries offering CGO bindings",
					Anchor:  "about-go-libraries-offering-cgo-bindings",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Software architecture: \"frameworks\" for beginners",
					Anchor:  "software-architecture-frameworks-for-beginners",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Paradigms: it's simply a point of view",
					Anchor:  "paradigms-its-simply-a-point-of-view",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Patterns: weapons of war",
					Anchor:  "patterns-weapons-of-war",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Best practices and standards: the final word",
					Anchor:  "best-practices-and-standards-the-final-word",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "A small review of Void Linux, a SystemD-less binary OS",
					Anchor:  "a-small-review-of-void-linux-a-systemd-less-binary-os",
					Content: "",
				},

				{
					Level:   2,
					Text:    "The review of Gentoo Linux",
					Anchor:  "the-review-of-gentoo-linux",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Installation Review",
							Anchor:  "installation-review",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Package Management Review",
							Anchor:  "package-management-review",
							Content: "",
						},

						{
							Level:   3,
							Text:    "About Kernel Configuration",
							Anchor:  "about-kernel-configuration",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Stability/Maintenance Review",
							Anchor:  "stabilitymaintenance-review",
							Content: "",
						},

						{
							Level:   3,
							Text:    "For other stuff review",
							Anchor:  "for-other-stuff-review",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Learning curve... what's about it?",
							Anchor:  "learning-curve-whats-about-it",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Tackling the myths about Gentoo Linux",
					Anchor:  "tackling-the-myths-about-gentoo-linux",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "\"It takes time to compile, and you compile every day.\"",
							Anchor:  "it-takes-time-to-compile-and-you-compile-every-day",
							Content: "",
						},

						{
							Level:   3,
							Text:    "\"It's unmaintainable if you forget about it for one month. It's unstable.\"",
							Anchor:  "its-unmaintainable-if-you-forget-about-it-for-one-month-its-unstable",
							Content: "",
						},

						{
							Level:   3,
							Text:    "\"It's bloated with dev dependencies.\"",
							Anchor:  "its-bloated-with-dev-dependencies",
							Content: "",
						},

						{
							Level:   3,
							Text:    "\"It's for expert.\"",
							Anchor:  "its-for-expert",
							Content: "",
						},

						{
							Level:   3,
							Text:    "\"It's optimized as hell.\"",
							Anchor:  "its-optimized-as-hell",
							Content: "",
						},

						{
							Level:   3,
							Text:    "\"Compiling my browser takes 100 hours.\"",
							Anchor:  "compiling-my-browser-takes-100-hours",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Understanding static-linked and dynamically-linked.",
					Anchor:  "understanding-static-linked-and-dynamically-linked",
					Content: "",
				},

				{
					Level:   2,
					Text:    "About Go compilation and CGO",
					Anchor:  "about-go-compilation-and-cgo",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Portage: Gentoo's package manager",
					Anchor:  "portage-gentoos-package-manager",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Crossdev: Gentoo's cross-compilation environment",
					Anchor:  "crossdev-gentoos-cross-compilation-environment",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Containerization and multi-arch manifests",
					Anchor:  "containerization-and-multi-arch-manifests",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Splitting the base image",
					Anchor:  "splitting-the-base-image",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "A critique about \"users in database\" and example of monolithic",
					Anchor:  "a-critique-about-users-in-database-and-example-of-monolithic",
					Content: "",
				},

				{
					Level:   2,
					Text:    "The existing standards for authentication",
					Anchor:  "the-existing-standards-for-authentication",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Hands-on example: One OIDC provider with Google Auth+OIDC with Go",
					Anchor:  "hands-on-example-one-oidc-provider-with-google-authoidc-with-go",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Second example: Multiple providers in one with Dex+OIDC with Go",
					Anchor:  "second-example-multiple-providers-in-one-with-dexoidc-with-go",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Third example: Fully self-hosted with 389ds+Dex+OIDC with Go",
					Anchor:  "third-example-fully-self-hosted-with-389dsdexoidc-with-go",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Criteria",
					Anchor:  "criteria",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Your first programming language for embedded systems: Arduino (C++)",
					Anchor:  "your-first-programming-language-for-embedded-systems-arduino-c",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Your first programming language for software",
					Anchor:  "your-first-programming-language-for-software",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Python, JavaScript and Typescript: Simple, but full of the bad practices, like the scripting languages they are.",
							Anchor:  "python-javascript-and-typescript-simple-but-full-of-the-bad-practices-like-the-scripting-languages-they-are",
							Content: "",
						},

						{
							Level:   3,
							Text:    "C: simple to write, hard to bootstrap, compile and debug",
							Anchor:  "c-simple-to-write-hard-to-bootstrap-compile-and-debug",
							Content: "",
						},

						{
							Level:   3,
							Text:    "C++: same as C, but with a better standard library and easier to use for OOP",
							Anchor:  "c-same-as-c-but-with-a-better-standard-library-and-easier-to-use-for-oop",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Rust: an alternative to C++, but a lot more \"safe\"",
							Anchor:  "rust-an-alternative-to-c-but-a-lot-more-safe",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Go: THE language that I recommend",
							Anchor:  "go-the-language-that-i-recommend",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Bonus: Zig: THE next language that I recommend",
							Anchor:  "bonus-zig-the-next-language-that-i-recommend",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "One last point: To OOP or not",
					Anchor:  "one-last-point-to-oop-or-not",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Introduction",
					Anchor:  "introduction",
					Content: "",
				},

				{
					Level:   2,
					Text:    "About the Linux boot process",
					Anchor:  "about-the-linux-boot-process",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Stateless images",
					Anchor:  "stateless-images",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "OverlayFS",
							Anchor:  "overlayfs",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Building the SquashFS image",
							Anchor:  "building-the-squashfs-image",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Install packages",
									Anchor:  "install-packages",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Add root password",
									Anchor:  "add-root-password",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Post-boot script",
									Anchor:  "post-boot-script",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Building the kernel, install kernel modules, packing the kernel",
									Anchor:  "building-the-kernel-install-kernel-modules-packing-the-kernel",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Packing the image with mksquashfs",
									Anchor:  "packing-the-image-with-mksquashfs",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Building the initramfs with Dracut",
							Anchor:  "building-the-initramfs-with-dracut",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Testing Locally",
							Anchor:  "testing-locally",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Booting the OS with PXE",
					Anchor:  "booting-the-os-with-pxe",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Limitations",
					Anchor:  "limitations",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "OverlayFS limitations",
							Anchor:  "overlayfs-limitations",
							Content: "",
						},

						{
							Level:   3,
							Text:    "NVIDIA drivers",
							Anchor:  "nvidia-drivers",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Privileges and host kernel",
							Anchor:  "privileges-and-host-kernel",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Motivation",
					Anchor:  "motivation",
					Content: "",
				},

				{
					Level:   2,
					Text:    "State of the art and Inspiration",
					Anchor:  "state-of-the-art-and-inspiration",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Docusaurus",
							Anchor:  "docusaurus",
							Content: "",
						},

						{
							Level:   3,
							Text:    "SvelteKit or SveltePress",
							Anchor:  "sveltekit-or-sveltepress",
							Content: "",
						},

						{
							Level:   3,
							Text:    "HTML-only",
							Anchor:  "html-only",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Taking inspiration of my favorite blog structure: The Go dev blog.",
							Anchor:  "taking-inspiration-of-my-favorite-blog-structure-the-go-dev-blog",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Hugo",
							Anchor:  "hugo",
							Content: "",
						},

						{
							Level:   3,
							Text:    "HTMX and Go",
							Anchor:  "htmx-and-go",
							Content: "",
						},
					},
				},

				{
					Level:   2,
					Text:    "Development",
					Anchor:  "development",
					Content: "",
					Children: []Header{

						{
							Level:   3,
							Text:    "Proof of Concept",
							Anchor:  "proof-of-concept",
							Content: "",
						},

						{
							Level:   3,
							Text:    "Architecture",
							Anchor:  "architecture",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Content directory",
									Anchor:  "content-directory",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Templates and Components",
									Anchor:  "templates-and-components",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Static directory",
									Anchor:  "static-directory",
									Content: "",
								},
							},
						},

						{
							Level:   3,
							Text:    "Implementation",
							Anchor:  "implementation",
							Content: "",
							Children: []Header{

								{
									Level:   4,
									Text:    "Initial Request and Server-Side rendering",
									Anchor:  "initial-request-and-server-side-rendering",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Markdown rendering",
									Anchor:  "markdown-rendering",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Compile-time rendering",
									Anchor:  "compile-time-rendering",
									Content: "",
								},

								{
									Level:   4,
									Text:    "Index page and pagination",
									Anchor:  "index-page-and-pagination",
									Content: "",
								},
							},
						},
					},
				},

				{
					Level:   2,
					Text:    "Conclusion",
					Anchor:  "conclusion",
					Content: "",
				},

				{
					Level:   2,
					Text:    "References",
					Anchor:  "references",
					Content: "",
				},
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
			Hierarchy: []Header{

				{
					Level:   2,
					Text:    "Table of contents",
					Anchor:  "table-of-contents",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Why ?",
					Anchor:  "why-",
					Content: "",
				},

				{
					Level:   2,
					Text:    "HTMX and Go",
					Anchor:  "htmx-and-go",
					Content: "",
				},

				{
					Level:   2,
					Text:    "Okay, so what is this blog ?",
					Anchor:  "okay-so-what-is-this-blog-",
					Content: "",
				},

				{
					Level:   2,
					Text:    "So, what's next ?",
					Anchor:  "so-whats-next-",
					Content: "",
				},
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
	Updated: time.Unix(1762734328, 0),
	Items: []*feeds.Item{
		{
			Title:       "I'm back! And I'm now flying FPV drones!",
			Description: "As an engineer, how I got started with FPV drones.",
			Created:     time.Unix(1753315200, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2025-07-24-fpv-drone",
			},
		},
		{
			Title:       "Pushing my Home Raspberry Pi cluster into a production state",
			Description: "A new year, an overhaul of my home Raspberry Pi cluster.",
			Created:     time.Unix(1737763200, 0),
			Link: &feeds.Link{
				Href: "https://blog.mnguyen.fr/blog/2025-01-25-home-raspi-part-2",
			},
		},
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
