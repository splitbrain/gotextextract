package main

import (
	"archive/zip"
	"encoding/xml"
	"flag"
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

	err := error(nil)

	// if type not given, use extension of filename
	if fileType == "" {
		fileType = strings.Trim(filepath.Ext(flag.Args()[0]), ".")
	}

	switch fileType {
	case "docx":
		err = dumpXMLFilesFromZip(flag.Args()[0], "word/document.xml")
		break
	case "odt":
		err = dumpXMLFilesFromZip(flag.Args()[0], "content.xml")
		break
	case "pdf":
		err = dumpPdf(flag.Args()[0])
		break
	case "pptx":
		err = dumpXMLFilesFromZip(flag.Args()[0], "ppt/slides/slide*.xml")
	case "odp":
		err = dumpXMLFilesFromZip(flag.Args()[0], "content.xml")
	default:
		log.Fatalf("Unsupported file type: '%s'", fileType)
	}

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return
}

// dumpPdf Dump the plain text contents from a PDF file to stdout.
func dumpPdf(path string) error {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return err
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
				print(word.S)
			}
		}
	}
	return nil
}

// dumpXMLFilesFromZip Dump the plain text contents of all matching XML files inside a zip document.
func dumpXMLFilesFromZip(path string, pattern string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		ok, _ := filepath.Match(pattern, f.Name)
		if !ok {
			continue
		}

		// Found it, print its content to terminal:
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		err = StripWordXMLTags(rc)
		if err != nil {
			return err
		}
	}
	return nil
}

// StripWordXMLTags Returns the raw text from a .docx Word document, adding newlines for paragraphs and table contents.
func StripWordXMLTags(r io.Reader) error {

	dec := xml.NewDecoder(r)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.CharData:
			print(string(tok))
		case xml.StartElement:
			// paragraph or table content
			if tok.Name.Local == "p" || tok.Name.Local == "tc" {
				println()
			}
		}
	}
	return nil
}
