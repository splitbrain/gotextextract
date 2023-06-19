package main

import (
	"archive/zip"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/ledongthuc/pdf"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	fileTypePtr := flag.String("type", "", "File type to extract text from")
	flag.Parse()
	fileType := *fileTypePtr

	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: gotextextract <filename>")
	}

	content := ""
	err := error(nil)

	// if type not given, use extension of filename
	if fileType == "" {
		fileType = strings.Trim(filepath.Ext(flag.Args()[0]), ".")
	}

	switch fileType {
	case "docx":
		content, err = readXMLFileFromZip(flag.Args()[0], "word/document.xml")
		break
	case "odt":
		content, err = readXMLFileFromZip(flag.Args()[0], "content.xml")
		break
	case "pdf":
		content, err = readPdf(flag.Args()[0])
		break
	default:
		log.Fatalf("Unsupported file type: '%s'", fileType)
	}

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Println(content)
	return
}

// readPdf Reads the plain text contents from a PDF file.
func readPdf(path string) (string, error) {
	var sb strings.Builder

	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				sb.WriteString(word.S)
			}
		}
	}
	return sb.String(), nil
}

// readXMLFileFromZip Reads the plain text contents from an XML file inside a zip document.
//
// This is used to extract the text from a .docx Word document.
func readXMLFileFromZip(path string, file string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != file {
			continue
		}

		// Found it, print its content to terminal:
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		return StripWordXMLTags(rc)
	}
	return "", nil
}

// StripWordXMLTags Returns the raw text from a .docx Word document, adding newlines for paragraphs and table contents.
func StripWordXMLTags(r io.Reader) (string, error) {
	var sb strings.Builder
	dec := xml.NewDecoder(r)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch tok := tok.(type) {
		case xml.CharData:
			sb.WriteString(string(tok))
		case xml.StartElement:
			// paragraph or table content
			if tok.Name.Local == "p" || tok.Name.Local == "tc" {
				sb.WriteString("\n")
			}
		}
	}
	return sb.String(), nil
}
