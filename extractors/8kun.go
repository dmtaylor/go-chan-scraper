package extractors

import (
	"io"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func Extract8Kun(body io.ReadCloser, fileChannel chan ImageFile) error {

	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return err
	}
	document.Find("p.fileinfo").Each(func(i int, s *goquery.Selection) {
		aNode := s.Find("a")
		if aNode.Length() == 0 {
			log.Printf("Node has no inner link: %s\n", s.Text())
		}
		link, exists := aNode.Attr("href")
		if !exists {
			log.Printf("Error: no href in element: %+v\n", s)
			return
		}
		fileLink, err := url.Parse("https:" + link)
		if err != nil {
			log.Printf("Failed to parse url %s: %v\n", link, err)
			return
		}
		var filename string
		// default to file name from link
		title, exists := aNode.Attr("title")
		if exists {
			filename = title
		} else {
			filename = s.Text()
		}

		// Use file name from inner span if it exists
		innerSpan := s.Find("span.postfilename")
		if innerSpan.Length() > 0 {
			title, exists := innerSpan.Attr("title")
			if exists {
				filename = title
			} else {
				filename = innerSpan.Text()
			}
		}
		imgFile := ImageFile{
			FileName: filename,
			FileURL:  *fileLink,
		}
		fileChannel <- imgFile
	})
	return nil
}
