package snippets

import (
	"encoding/json"
	"fmt"
	"os"
	"sazed/utils"
	"strings"
	"time"
)

type snippetMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
	Uses        int      `json:"uses"`
	code        string
	notes       string
	tests       string
}

func NewSnippetMetadata(name, description string) *snippetMetadata {
	return &snippetMetadata{
		Name:        name,
		Description: description,
		Date:        time.Now().Format("03-01-2006"),
		Tags:        []string{},
	}
}

type Metadata struct {
	FileType string             `json:"filetype"`
	Snippets []*snippetMetadata `json:"snippets"`
	path     string
}

func NewMetadata(filetype, path string) Metadata {
	filePath := fmt.Sprintf("%s%smetadata.json", path, string(os.PathSeparator))
	if _, err := os.Stat(filePath); err == nil {
		// if metadata Exists, read it and return it
		metadata := ReadMetadata(filePath)
		metadata.path = path
		return metadata
	}

	metadata := Metadata{
		FileType: filetype,
		Snippets: []*snippetMetadata{},
		path:     path,
	}

	utils.CreateDirIfNotExist(path)
	utils.CreateFileIfNotExist(filePath)

	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	utils.CheckErr(err)

	err = os.WriteFile(filePath, jsonData, 0644)
	utils.CheckErr(err)

	return metadata
}

func ReadMetadata(filePath string) Metadata {
	file, err := os.ReadFile(filePath)
	utils.CheckErr(err)
	var metadata Metadata

	err = json.Unmarshal(file, &metadata)
	utils.CheckErr(err)
	return metadata
}

func (m *Metadata) Save() {
	// metadata
	jsonData, err := json.MarshalIndent(m, "", "  ")
	utils.CheckErr(err)

	filePath := strings.Join([]string{m.path, "metadata.json"}, string(os.PathSeparator))
	err = os.WriteFile(filePath, jsonData, 0644)
	utils.CheckErr(err)
}

func (m *Metadata) AddSnippet(name, description string) {
	snip := NewSnippetMetadata(name, description)
	m.Snippets = append(m.Snippets, snip)
}

func (m *Metadata) RemoveSnippet(name string) {
	for i, snippet := range m.Snippets {
		if snippet.Name == name {
			m.Snippets = append(m.Snippets[:i], m.Snippets[i+1:]...)
			return
		}
	}
}

func (m *Metadata) GetSnippet(name string) *snippetMetadata {
	for _, snippet := range m.Snippets {
		if snippet.Name == name {
			return snippet
		}
	}
	return nil
}

func (m *Metadata) AddTag(snippetName, tag string) {
	snippet := m.GetSnippet(snippetName)
	if snippet == nil {
		return
	}
	snippet.Tags = append(snippet.Tags, tag)
	m.Save()
}

func (m *Metadata) RemoveTag(snippetName, tag string) {
	snippet := m.GetSnippet(snippetName)
	if snippet == nil {
		return
	}
	for i, t := range snippet.Tags {
		if t == tag {
			snippet.Tags = append(snippet.Tags[:i], snippet.Tags[i+1:]...)
			return
		}
	}

	m.Save()
}

func (m *Metadata) SetTags(snippetName string, tags []string) {
	snippet := m.GetSnippet(snippetName)
	if snippet == nil {
		return
	}

	snippet.Tags = tags

	m.Save()
}

func (m *Metadata) SetCode(snippetName, code string) {
	snippet := m.GetSnippet(snippetName)
	if snippet == nil {
		return
	}
	snippet.code = code

	dirPath := strings.Join([]string{m.path, snippet.Name}, string(os.PathSeparator))
	utils.CreateDirIfNotExist(dirPath)

	codeFilePath := strings.Join([]string{
		dirPath,
		fmt.Sprintf("code.%s", m.FileType),
	}, string(os.PathSeparator))
	err := os.WriteFile(codeFilePath, []byte(snippet.code), 0644)
	utils.CheckErr(err)

	m.Save()
}

func (m *Metadata) SetTest(snippetName, tests string) {
	snippet := m.GetSnippet(snippetName)
	if snippet == nil {
		return
	}
	snippet.tests = tests

	dirPath := strings.Join([]string{m.path, snippet.Name}, string(os.PathSeparator))
	utils.CreateDirIfNotExist(dirPath)

	testsFilePath := strings.Join([]string{
		dirPath,
		fmt.Sprintf("tests.%s", m.FileType),
	}, string(os.PathSeparator))

	err := os.WriteFile(testsFilePath, []byte(snippet.code), 0644)
	utils.CheckErr(err)

	m.Save()
}

func (m *Metadata) SetNotes(snippetName, notes string) {
	snippet := m.GetSnippet(snippetName)
	if snippet == nil {
		return
	}
	snippet.notes = notes

	dirPath := strings.Join([]string{m.path, snippet.Name}, string(os.PathSeparator))
	utils.CreateDirIfNotExist(dirPath)

	notesFilePath := strings.Join([]string{dirPath, "notes.md"}, string(os.PathSeparator))
	err := os.WriteFile(notesFilePath, []byte(snippet.notes), 0644)
	utils.CheckErr(err)

	m.Save()
}
