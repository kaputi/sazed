package utils

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"
)

type fsEntry struct {
	path  string
	name  string
	ext   string
	isDir bool
}

func (e fsEntry) Name() string {
	return e.name
}

func (e fsEntry) Path() string {
	return e.path
}

func (e fsEntry) IsDir() bool {
	return e.isDir
}

func (e fsEntry) Ext() string {
	return e.ext
}

func ReadDir(dirPath string) []fsEntry {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	entries := make([]fsEntry, len(files))

	dirPathParts := strings.Split(dirPath, string(os.PathSeparator))
	for i, file := range files {
		entries[i] = fsEntry{
			path:  strings.Join(append(dirPathParts, file.Name()), string(os.PathSeparator)),
			name:  file.Name(),
			isDir: file.IsDir(),
			ext:   path.Ext(file.Name()),
		}
	}

	return entries
}

func CreateDirIfNotExist(dirPath string) {
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
