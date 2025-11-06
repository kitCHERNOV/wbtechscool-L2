package wgetmanager

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"wget/include/handlers"
	"wget/include/logger"
)

func WgetManager() {
	// logger initializing
	logger := logger.NewLogger()

	var (
		url string = "https://parsemachine.com/sandbox/"//"https://www.example.com"
	)

	scanner := bufio.NewScanner(os.Stdin)
	// read url from terminal
	scanner.Scan()
	url = scanner.Text()

	// execute a get request
	includedLinks := handlers.NewHtmlIncludedLinksStruct()
	response, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(logger, "get request error: %v\n", err)
	}
	fmt.Fprintf(logger, "url is gotten\n")
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(logger, "response status is not OK\n")
	}

	includedLinks.HtmlParser(response.Body)

	includedLinks.DownloadPages()
}
