# Go Text Extractor

This is meant as a simple way to extract raw text contents from different file formats to be used in search indexing. It is not meant to display contents true to their orignal layout.

It currently supports the following file formats:

  * `pdf` -- using [ledongthuc/pdf](https://github.com/ledongthuc/pdf) 
  * `docx` -- naive extraction from the xml
  * `odt` -- naive extraction from the xml 