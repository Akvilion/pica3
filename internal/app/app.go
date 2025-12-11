package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App структура для зберігання стану застосунку
type App struct {
	ctx          context.Context
	currentImage string
	imageFiles   []string
	currentIndex int
}

// NewApp створює новий екземпляр застосунку
func NewApp() *App {
	return &App{
		imageFiles:   []string{},
		currentIndex: -1,
	}
}

// Startup викликається при запуску застосунку
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Shutdown викликається при закритті застосунку
func (a *App) Shutdown(ctx context.Context) {
	// Тут можна виконати очищення ресурсів
}

// DomReady викликається після завантаження DOM
func (a *App) DomReady(ctx context.Context) {
	// Викликається після завантаження фронтенду
}

// OpenImage відкриває діалог вибору файлу та завантажує зображення
func (a *App) OpenImage() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Відкрити зображення",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Зображення",
				Pattern:     "*.jpg;*.jpeg;*.png;*.gif;*.bmp;*.webp",
			},
			{
				DisplayName: "Всі файли",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return "", err
	}

	if file == "" {
		return "", fmt.Errorf("файл не вибрано")
	}

	return a.LoadImage(file)
}

// LoadImage завантажує зображення та повертає його у форматі base64
func (a *App) LoadImage(filePath string) (string, error) {
	// Читаємо файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Визначаємо MIME тип за розширенням
	ext := strings.ToLower(filepath.Ext(filePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".bmp":
		mimeType = "image/bmp"
	case ".webp":
		mimeType = "image/webp"
	default:
		mimeType = "image/jpeg"
	}

	// Конвертуємо у base64
	encoded := base64.StdEncoding.EncodeToString(data)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	a.currentImage = filePath

	// Завантажуємо список файлів у директорії
	a.loadImagesFromDirectory(filepath.Dir(filePath))

	return dataURL, nil
}

// loadImagesFromDirectory завантажує список всіх зображень з директорії
func (a *App) loadImagesFromDirectory(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	a.imageFiles = []string{}
	supportedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if supportedExts[ext] {
				fullPath := filepath.Join(dir, file.Name())
				a.imageFiles = append(a.imageFiles, fullPath)
				if fullPath == a.currentImage {
					a.currentIndex = len(a.imageFiles) - 1
				}
			}
		}
	}

	return nil
}

// NextImage переходить до наступного зображення
func (a *App) NextImage() (string, error) {
	if len(a.imageFiles) == 0 {
		return "", fmt.Errorf("немає зображень")
	}

	a.currentIndex = (a.currentIndex + 1) % len(a.imageFiles)
	return a.LoadImage(a.imageFiles[a.currentIndex])
}

// PreviousImage переходить до попереднього зображення
func (a *App) PreviousImage() (string, error) {
	if len(a.imageFiles) == 0 {
		return "", fmt.Errorf("немає зображень")
	}

	a.currentIndex--
	if a.currentIndex < 0 {
		a.currentIndex = len(a.imageFiles) - 1
	}
	return a.LoadImage(a.imageFiles[a.currentIndex])
}

// GetImageInfo повертає інформацію про поточне зображення
func (a *App) GetImageInfo() map[string]interface{} {
	return map[string]interface{}{
		"current": a.currentIndex + 1,
		"total":   len(a.imageFiles),
		"path":    a.currentImage,
	}
}
