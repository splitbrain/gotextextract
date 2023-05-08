package main

import (
    "fmt"
    "os"

    "github.com/ledongthuc/pdf"
    "github.com/nguyenthenguyen/docx"
)

func main() {
    content, err := readPdf(os.Args[1]) // Read local pdf file
    if err != nil {
        panic(err)
    }
    fmt.Println(content)
    return
}

func readDocx(path string) (string, error) {
    r, err := docx.ReadDocxFile("./template.docx")

    if err != nil {
        panic(err)
    }

    docx.readText(zipfile)

}

func readPdf(path string) (string, error) {
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
                fmt.Println(word.S)
            }
        }
    }
    return "", nil
}
