package wgetmanager

import (
	_ "bufio"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"wget/include/handlers"
	"wget/include/logger"
)


func WgetManager() {
	// logger initializing
	logger := logger.NewLogger()


	var (
		url string = "https://www.example.com"
	)

	// scanner := bufio.NewScanner(os.Stdin)
	// // read url from terminal
	// scanner.Scan()
	// url = scanner.Text()

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

	// create dirs
	filePath := "./downloads"

	// directory isExist checking
	dir := filepath.Dir(filePath)
	if err := os.Mkdir(dir, 0755); err != nil {
		fmt.Fprintf(logger, "creating downloads directory error")
		os.Exit(1)
	}

	dirCss := filepath.Dir(filePath + "css")
	if err := os.Mkdir(dirCss, 0755); err != nil {
		fmt.Fprintf(logger, "creating css directory error")
		os.Exit(1)
	}

	includedLinks.DownloadPages()
	// TODO: to wrap next code in for loop
	// for {
	// }
}