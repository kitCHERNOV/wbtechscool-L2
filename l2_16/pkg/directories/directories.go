package directories

import (
	"log"
	"os"
	"path/filepath"
)

func ToCreateDirectories() {
	// create dirs
	filePath := "./downloads"

	// directory isExist checking
	dir := filepath.Dir(filePath)
	if err := os.Mkdir(dir, 0755); err != nil {
		log.Fatalln("creating downloads directory error")
	}

	dirCss := filepath.Dir(filePath + "css")
	if err := os.Mkdir(dirCss, 0755); err != nil {
		log.Fatalln("creating css directory error")
	}

	dirJs := filepath.Dir(filePath + "js")
	if err := os.Mkdir(dirJs, 0755); err != nil {
		log.Fatalln("creating js directory error")
	}

	dirHtml := filepath.Dir(filePath + "html")
	if err := os.Mkdir(dirHtml, 0755); err != nil {
		log.Fatalln("creating html directory error")
	}

}
