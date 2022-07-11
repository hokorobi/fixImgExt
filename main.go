package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

func main() {
	// カレントディレクトリ配下のディレクトリ取得
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
		err = os.Rename(path, newPath)
		logg("Rename: " + path + " to " + newPath)
		if err != nil {
			logf(err)
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

	_, format, err := image.DecodeConfig(f)
	// 画像以外は無視
	if format == "" {
		return ""
	}
	if err != nil {
		logf(err)
	}

	if format == "jpeg" {
		format = "jpg"
	}
	return "." + format
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
