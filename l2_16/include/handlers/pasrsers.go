package handlers

import (
	"fmt"
	"io"
	"wget/include/logger"

	"golang.org/x/net/html"
)

type htmlIncludedLinks struct {
	logger logger.Logger
	CssLinks []string
	JsLinks []string
	Images []string
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
	
}