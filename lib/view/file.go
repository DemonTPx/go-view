package view

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var supportedFileExtensions = map[string]bool{
	".bmp":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

type FileCursor struct {
	directory string
	files     []string

	current int
}

func NewFileCursor(directory string) (*FileCursor, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var images []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if FileExtensionSupported(file.Name()) {
			images = append(images, file.Name())
		}
	}

	return &FileCursor{
		directory: directory,
		files:     images,
	}, nil
}

func NewFileCursorFromFilename(filename string) (*FileCursor, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return NewFileCursor(filename)
	}

	cursor, err := NewFileCursor(filepath.Dir(filename))
	if err != nil {
		return cursor, err
	}

	if !FileExtensionSupported(filename) {
		log.Println("file extension not supported, starting at first file in directory")
		return cursor, err
	}

	basename := filepath.Base(filename)

	for index, file := range cursor.files {
		if file == basename {
			cursor.current = index
			break
		}
	}

	return cursor, err
}

func NewFileCursorFromWorkingDirectory() (*FileCursor, error) {
	directory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return NewFileCursor(directory)
}

func FileExtensionSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if _, ok := supportedFileExtensions[ext]; ok {
		return true
	}
	return false
}

func (c *FileCursor) GetFilename() string {
	if len(c.files) == 0 {
		return ""
	}
	return filepath.Join(c.directory, c.files[c.current])
}

func (c *FileCursor) First() {
	c.current = 0
}

func (c *FileCursor) Last() {
	c.current = len(c.files) - 1
}

func (c *FileCursor) Next() {
	c.current = c.current + 1
	if c.current >= len(c.files) {
		c.current = 0
	}
}

func (c *FileCursor) Previous() {
	c.current = c.current - 1
	if c.current < 0 {
		c.current = len(c.files) - 1
	}
}
