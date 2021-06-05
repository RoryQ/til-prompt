package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readEntry(t *testing.T) {
	md := `# Title

This is what I learned.

---
Date: 2021-06-05
`
	t.Run(" category", func(t *testing.T) {
		entry := readEntry(md)

		assert.Equal(t, "Title", entry.Title)
		assert.Equal(t, "This is what I learned.", entry.Body)
		assert.Equal(t, "2021-06-05", entry.DateString)

		// category is optional
		assert.Equal(t, "", entry.Category)
		entry = readEntry(md + "\nCategory: gcp")
		assert.Equal(t, "gcp", entry.Category)
	})
}

func Test_renderCategoryLinks(t *testing.T) {
	t.Run("NoCategories", func(t *testing.T) {
		s := renderCategoryLinks(ReadmeContents{})
		assert.Equal(t, "", s)
	})

	t.Run("WithCategories", func(t *testing.T) {
		contents := ReadmeContents{
			Categories: []string{"gcp", "golang"},
		}
		s := renderCategoryLinks(contents)
		expected := `### Categories
* [gcp](gcp)
* [golang](golang)
`
		assert.Equal(t, expected, s)
	})
}

func Test_renderEntryLinks(t *testing.T) {
	entry := Entry{
		SavePath: "2021-06-05-title.md",
		Title:    "Title",
	}
	t.Run("NoCategories", func(t *testing.T) {
		contents := ReadmeContents{Entries: map[string][]Entry{
			"": {entry, entry},
		}}
		s := renderEntryLinks(contents)
		expected := `---

* [Title](2021-06-05-title.md)
* [Title](2021-06-05-title.md)

`
		assert.Equal(t, expected, s)
	})

	t.Run("WithCategories", func(t *testing.T) {
		contents := ReadmeContents{
			Categories: []string{"gcp"},
			Entries: map[string][]Entry{
				"gcp": {entry, entry},
			},
			TotalCount: 0,
		}
		s := renderEntryLinks(contents)
		expected := `---

### gcp

* [Title](2021-06-05-title.md)
* [Title](2021-06-05-title.md)

`
		assert.Equal(t, expected, s)
	})
}
