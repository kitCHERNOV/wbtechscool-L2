package handlers

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
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

	// stringSliceType for checking fields
	stringSliceType := reflect.TypeOf([]string{})
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type() != stringSliceType {
			continue
		}
		for _, url := range v.Field(i).Interface().([]string) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(Logger, "getting url: %s error: %v", url, err)
			}

			field := v.Field(i)
			fieldName := field.

		}
	}
}
