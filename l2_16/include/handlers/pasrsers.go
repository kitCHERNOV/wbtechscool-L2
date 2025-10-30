package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"wget/include/logger"

	"golang.org/x/net/html"
)

type htmlIncludedLinks struct {
	logger   logger.Logger
	CssLinks []string
	JsLinks  []string
	Images   []string
}

func NewHtmlIncludedLinksStruct() *htmlIncludedLinks {
	return &htmlIncludedLinks{}
}

func (inclLinks *htmlIncludedLinks) HtmlParser(object io.Reader) {
	htmlNode, err := html.Parse(object)
	if err != nil {
		fmt.Fprintf(&inclLinks.logger, "parse html error: %v", err)
	}

	// firstly get all links
	findCertainTagAndItsData(htmlNode, inclLinks)
}

func (incLinks *htmlIncludedLinks) DownloadPages() {
	Logger := &incLinks.logger

	v := reflect.ValueOf(incLinks).Elem()
	t := v.Type()
	// stringSliceType for checking fields
	stringSliceType := reflect.TypeOf([]string{})
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type() != stringSliceType {
			continue
		}
		// get description of elem
		field := v.Field(i)
		fieldName := t.Field(i).Name

		var folderName string
		switch fieldName {
        case "CssLinks":
            folderName = "css"
        case "JsLinks":
            folderName = "js"
        case "Images":
            folderName = "images"
        default:
            continue
        }
		// all directories shoud be created in main func
		urls := field.Interface().([]string)
		for _, url := range urls {
			
			
		}
	}
}

func downloadFile(urlStr, folder string, logger *logger.Logger) {
	// get file name
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		fmt.Fprintf(logger, "parse url error: %w", err)
		return
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "/" || fileName == "." {
		// if file name is not determined, just generate its
		fileName = generateFileName(urlStr, folder)
	}

	filePath := filepath.Join(folder, fileName)

	// download file
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Fprintf(logger, "getting url: %s error: %v", urlStr, err)
	}

	defer resp.Body.Close()

	// Check download status
	if resp.StatusCode != http.StatusOK {
		
	}
}


func generateFileName(urlStr, folder string) string {
    // Используем хеш URL или timestamp
    hash := fmt.Sprintf("%d", len(urlStr)) // Простой пример
    
    var ext string
    switch folder {
    case "css":
        ext = ".css"
    case "js":
        ext = ".js"
    case "images":
        // Можно попробовать определить из Content-Type
        ext = ".png" // По умолчанию
    }
    
    return fmt.Sprintf("file_%s%s", hash, ext)
}