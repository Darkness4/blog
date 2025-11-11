package meilisearch

import (
	"iter"

	"github.com/Darkness4/blog/web/gen/index"
	"github.com/rs/zerolog/log"
)

type RecordWithFormat struct {
	Record    `json:",inline"`
	Formatted Record `json:"_formatted"`
}

type Record struct {
	ObjectID      string `json:"objectID"`
	HierarchyLvl0 string `json:"hierarchy_lvl0,omitempty"`
	HierarchyLvl1 string `json:"hierarchy_lvl1,omitempty"`
	HierarchyLvl2 string `json:"hierarchy_lvl2,omitempty"`
	HierarchyLvl3 string `json:"hierarchy_lvl3,omitempty"`
	HierarchyLvl4 string `json:"hierarchy_lvl4,omitempty"`
	HierarchyLvl5 string `json:"hierarchy_lvl5,omitempty"`
	HierarchyLvl6 string `json:"hierarchy_lvl6,omitempty"`
	Content       string `json:"content"`
	URL           string `json:"url"`
	Anchor        string `json:"anchor"`
}

// IndexToRecords converts an Index to a slice of search Records
func IndexToRecords(index [][]index.Index) iter.Seq[Record] {

	return func(yield func(Record) bool) {
		for _, i := range index {
			for _, j := range i {
				// Base hierarchy levels
				lvl0 := j.PublishedDate.Format("January 2006")
				lvl1 := j.Title

				// 1. **Directly yield the Level 2 Record**
				lvl2Record := Record{
					ObjectID:      j.EntryName,
					HierarchyLvl0: lvl0,
					HierarchyLvl1: lvl1,
					HierarchyLvl2: "",
					HierarchyLvl3: "",
					HierarchyLvl4: "",
					HierarchyLvl5: "",
					HierarchyLvl6: "",
					Content:       "",
					URL:           j.Href,
					Anchor:        "",
				}

				if !yield(lvl2Record) {
					return
				}

				// 2. **Pass the yield function to processHeader**
				for _, header := range j.Hierarchy {
					// Check if we need to stop based on the return value of processHeader
					if !processHeader(
						header,
						j,
						lvl0,
						lvl1,
						"",
						"",
						"",
						"",
						"",
						yield, // Pass the yield function directly
					) {
						return // Stop the entire sequence
					}
				}
			}
		}
	}
}

// processHeader recursively processes headers and creates records, yielding them directly.
func processHeader(
	h index.Header,
	idx index.Index, // Renamed 'index' to 'idx' to avoid shadowing the package name
	lvl0, lvl1, lvl2, lvl3, lvl4, lvl5, lvl6 string,
	yield func(Record) bool,
) bool { // Now returns a boolean to indicate continuance (true = continue, false = stop)

	switch h.Level {
	case 1:
		// Skip H1 headers (assuming the top-level article title is lvl2)
		log.Warn().Msg("Skipping H1 header: conflicts with article title")
		return true // Continue processing other headers
	case 2:
		lvl2 = h.Text
	case 3:
		lvl3 = h.Text
	case 4:
		lvl4 = h.Text
	case 5:
		lvl5 = h.Text
	case 6:
		lvl6 = h.Text
	default:
		// Ignore headers beyond H5
		return true // Continue processing other headers
	}

	// 1. Create record for this header
	record := Record{
		ObjectID:      idx.EntryName + "-" + h.Anchor,
		HierarchyLvl0: lvl0,
		HierarchyLvl1: lvl1,
		HierarchyLvl2: lvl2,
		HierarchyLvl3: lvl3,
		HierarchyLvl4: lvl4,
		HierarchyLvl5: lvl5,
		HierarchyLvl6: lvl6,
		Content:       h.Content,
		URL:           idx.Href + "#" + h.Anchor,
		Anchor:        h.Anchor,
	}

	// 2. **Directly yield the record**
	if !yield(record) {
		return false // Stop if the consumer doesn't want more records
	}

	// 3. Process children, passing down the **new** hierarchy levels
	for _, child := range h.Children {
		// If a child call returns false, stop everything and propagate the 'false'
		if !processHeader(
			child,
			idx,
			lvl0,
			lvl1,
			lvl2,
			lvl3,
			lvl4,
			lvl5,
			lvl6,
			yield,
		) {
			return false
		}
	}

	return true // Continue to the next sibling header
}
