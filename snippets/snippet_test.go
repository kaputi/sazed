package snippets

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var path = "snippets.test"
var ft = "go"

var code = `package main

  import "fmt"

  func main() {
    fmt.Pringln("Hello, World!")
  }`

var notes = `This is a test snippet`

var tests = `package main
  import "testing"

  func TestMain(t *testing.T) {

  }`

func removeTest() {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func TestSnippet(t *testing.T) {
	removeTest()
	metadata := NewMetadata(ft, path)
	assert.Equal(t, ft, metadata.FileType)

	metadata.FileType = "js"
	metadata.Save()
	metadata = ReadMetadata(strings.Join([]string{path, "metadata.json"}, string(os.PathSeparator)))
	metadata.path = path
	assert.Equal(t, "js", metadata.FileType)

	metadata.FileType = ft

	metadata.AddSnippet("test", "this is just a test")

	TestSnippet := metadata.GetSnippet("test")
	assert.Equal(t, "test", TestSnippet.Name)

	metadata.RemoveSnippet("test")
	TestSnippet = metadata.GetSnippet("test")
	assert.Nil(t, TestSnippet)

	metadata.AddSnippet("test", "this is just a test")
	metadata.AddTag("test", "tag1")
	metadata.AddTag("test", "tag2")
	metadata.RemoveTag("test", "tag1")

	metadata.SetCode("test", code)

	metadata.SetTest("test", tests)
	metadata.SetNotes("test", notes)

  removeTest()
}
