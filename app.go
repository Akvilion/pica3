package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OpenImageDialog opens a file dialog to select an image
func (a *App) OpenImageDialog() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Image",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images",
				Pattern:     "*.png;*.jpg;*.jpeg;*.gif;*.webp",
			},
		},
	})

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}

	return selection
}

// GetImagesInFolder returns a list of image files in the directory
func (a *App) GetImagesInFolder(dirPath string) []string {
	if dirPath == "" {
		return []string{}
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}
	}

	var paths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".webp" {
				paths = append(paths, filepath.Join(dirPath, entry.Name()))
			}
		}
	}
	return paths
}

// LoadImage loads an image from disk and returns it as a base64 string
func (a *App) LoadImage(path string) string {
	if path == "" {
		return ""
	}

	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	ext := strings.ToLower(filepath.Ext(path))
	mimeType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	}

	return "data:" + mimeType + ";base64," + encoded
}
