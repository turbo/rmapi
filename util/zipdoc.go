package util

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/juruen/rmapi/log"
)

type content struct {
	ExtraMetadata  map[string]string `json:"extraMetadata"`
	FileType       string            `json:"fileType"`
	LastOpenedPage int               `json:"lastOpenedPage"`
	LineHeight     int               `json:"lineHeight"`
	Margins        int               `json:"margins"`
	TextScale      int               `json:"textScale"`
	Transform      map[string]string `json:"transform"`
}

func CreateZipDocument(id, srcPath string) (string, error) {
	pdf, err := ioutil.ReadFile(srcPath)

	if err != nil {
		log.Error.Println("failed to open source PDF file to read", err)
		return "", err
	}

	tmp, err := ioutil.TempFile("", "rmapizip")
	log.Trace.Println("creating temp zip file:", tmp.Name())
	defer tmp.Close()

	if err != nil {
		log.Error.Println("failed to create tmpfile for zip doc", err)
		return "", err
	}

	w := zip.NewWriter(tmp)
	defer w.Close()

	// Create PDF file
	f, err := w.Create(fmt.Sprintf("%s.pdf", id))
	if err != nil {
		log.Error.Println("failed to create pdf entry in zip file", err)
		return "", err
	}

	f.Write(pdf)

	// Create pagedata file
	f, err = w.Create(fmt.Sprintf("%s.pagedata", id))
	if err != nil {
		log.Error.Println("failed to create content entry in zip file", err)
		return "", err
	}

	f.Write(make([]byte, 0))

	// Create content content
	f, err = w.Create(fmt.Sprintf("%s.content", id))
	if err != nil {
		log.Error.Println("failed to create content entry in zip file", err)
		return "", err
	}

	c, err := createContent()
	if err != nil {
		return "", err
	}

	f.Write([]byte(c))

	return tmp.Name(), nil
}

func createContent() (string, error) {
	c := content{
		make(map[string]string),
		"pdf",
		0,
		-1,
		180,
		1,
		make(map[string]string),
	}

	cstring, err := json.Marshal(c)

	log.Trace.Println("content: ", string(cstring))

	if err != nil {
		log.Error.Println("failed to serialize content file", err)
		return "", err
	}

	return string(cstring), nil
}
