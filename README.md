# Go Text Extractor

This is meant as a simple way to extract raw text contents from different file formats to be used in search indexing. It is not meant to display contents true to their orignal layout.

It currently supports the following file formats:

  * `pdf` -- PDF using [ledongthuc/pdf](https://github.com/ledongthuc/pdf) 
  * `docx` -- Microsoft Word, naive extraction from the xml
  * `odt` -- Open/Libreoffice Document, naive extraction from the xml
  * `pptx` -- Microsoft Powerpoint, naive extraction from the xml
  * `odp` -- Open/Libreoffice Presentation, naive extraction from the xml

## Usage

    gotextextract [--type <type>] <file>

Simply give the file to extract as argument. If no file type (see above) is given, it will try to guess it from the file extension. The extracted text will be printed to stdout.
