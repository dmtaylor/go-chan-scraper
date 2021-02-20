package extractors

import (
	"io"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func Extract4Chan(body io.ReadCloser, fileChannel chan ImageFile) error {
	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return err
	}
	document.Find(".fileText a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			log.Printf("Error: no href in element: %+v\n", s)
			return
		}
		fileLink, err := url.Parse("https:" + link)
		if err != nil {
			log.Printf("Failed to parse url %s: %v\n", link, err)
		}
		var filename string
		title, exists := s.Attr("title")
		if exists {
			filename = title
		} else {
			filename = s.Text()
		}
		imgFile := ImageFile{
			FileName: filename,
			FileURL:  *fileLink,
		}
		fileChannel <- imgFile
	})

	return nil
}
