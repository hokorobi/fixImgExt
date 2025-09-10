package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// カレントディレクトリ配下の画像ファイルの拡張子修正
	pwd, err := os.Getwd()
	if err != nil {
		logf(err)
	}

	err = filepath.Walk(pwd, walkFn())
	if err != nil {
		logf(err)
	}
}

func walkFn() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		ext := getImgExt(path)
		if ext == "" {
			return nil
		}
		if filepath.Ext(path) == ext {
			return nil
		}

		newPath := getFullpathWithoutExt(path) + ext
		err2 := os.Rename(path, newPath)
		logg("Rename: " + path + " to " + newPath)
		if err2 != nil {
			logf(err2)
		}

		return nil
	}
}
func getFullpathWithoutExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}
func getImgExt(path string) string {
	f, err := os.Open(path)
	if err != nil {
		logf(err)
	}
	defer f.Close()

	buffer := make([]byte, 512)
	f.Read(buffer)
	contentType := http.DetectContentType(buffer)

	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/bmp":
		return ".bmp"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		if strings.HasPrefix(contentType, "image/") {
			logg(path + ": " + contentType)
		}
		return ""
	}
}

func logf(m interface{}) {
	logg(m)
	os.Exit(1)
}
func logg(m interface{}) {
	f, err := os.OpenFile(getFilename(".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Cannot open log file: " + err.Error())
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stderr))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(m)
}
func getFilename(ext string) string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), getFileNameWithoutExt(exec)+ext)
}
func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
