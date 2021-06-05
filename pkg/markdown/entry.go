package markdown

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed template
var templatesFS embed.FS

type Entry struct {
	SavePath   string
	Title      string
	Body       string
	Category   string
	DateString string
}

func (e Entry) Save(saveDirectory string) error {
	content, err := e.Render()
	if err != nil {
		return err
	}
	fullPath := path.Join(saveDirectory, e.SavePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return err
	}

	return e.UpdateReadme(saveDirectory)
}

func (e Entry) Render() (string, error) {
	entryTemplate, err := template.ParseFS(templatesFS, "template/entry.gomd")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := entryTemplate.Execute(&buf, e); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e Entry) UpdateReadme(saveDirectory string) error {
	readme, err := RenderReadme(saveDirectory)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(saveDirectory, "README.md"), []byte(readme), 0644)
}

type ReadmeContents struct {
	Categories []string
	Entries    map[string][]Entry
	TotalCount int
}

func RenderReadme(saveDirectory string) (string, error) {
	contents, err := generateContents(os.DirFS(saveDirectory))
	if err != nil {
		return "", err
	}

	b, err := templatesFS.ReadFile("template/README.gomd")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = template.Must(template.New("").Funcs(template.FuncMap{"RenderLinks": renderLinks}).
		Parse(string(b))).
		Execute(&buf, contents)

	return buf.String(), err
}

func generateContents(fileSystem fs.FS) (ReadmeContents, error) {
	contents := ReadmeContents{
		Categories: []string{},
		Entries:    make(map[string][]Entry),
	}

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == "README.md" || strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if d.IsDir() {
			contents.Categories = append(contents.Categories, d.Name())
			return nil
		}

		bytes, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return err
		}

		entry := readEntry(string(bytes))
		entry.SavePath = path
		if arr, ok := contents.Entries[entry.Category]; ok {
			arr = append(arr, entry)
			contents.Entries[entry.Category] = arr
			contents.TotalCount++
		} else if !ok {
			contents.Entries[entry.Category] = []Entry{entry}
			contents.TotalCount++
		}

		return nil
	})

	return contents, err
}

func renderCategoryLinks(contents ReadmeContents) string {
	b := &strings.Builder{}
	for _, cat := range contents.Categories {
		link(b, cat, cat)
	}

	if b.Len() > 0 {
		return "---\n### Categories\n\n" + b.String()
	}

	return ""
}

func link(b *strings.Builder, text, url string) {
	b.WriteString(fmt.Sprintf("* [%s](%s)\n", text, url))
}

func renderEntryLinks(contents ReadmeContents) string {
	b := &strings.Builder{}
	for _, e := range contents.Entries[""] {
		link(b, e.Title, e.SavePath)
	}

	for _, cat := range contents.Categories {
		b.WriteString(fmt.Sprintf("### %s\n\n", cat))

		for _, e := range contents.Entries[cat] {
			link(b, e.Title, e.SavePath)
		}
		b.WriteRune('\n')
	}

	if b.Len() > 0 {
		return "---\n\n" + b.String()
	}

	return ""
}

func renderLinks(contents ReadmeContents) string {
	categories := renderCategoryLinks(contents)
	entries := renderEntryLinks(contents)
	return strings.Join([]string{categories, entries}, "\n")
}

var (
	entryRE = regexp.MustCompile(
		`(?sm)^# (?P<title>[^\n]+)\n` +
			`(?P<body>.*)\n---` +
			`(\n?)+Date: (?P<date>\d{4}-\d\d-\d\d)(\n?)+` +
			`(Category: (?P<category>[^\n]+))?`)
)

func readEntry(s string) (entry Entry) {
	matches := entryRE.FindStringSubmatch(s)
	entry.Title = matches[entryRE.SubexpIndex("title")]
	entry.Body = strings.TrimSpace(matches[entryRE.SubexpIndex("body")])
	entry.DateString = matches[entryRE.SubexpIndex("date")]
	entry.Category = matches[entryRE.SubexpIndex("category")]
	return entry
}
