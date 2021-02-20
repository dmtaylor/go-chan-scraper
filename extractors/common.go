package extractors

import "net/url"

type ImageFile struct {
	FileName string
	FileURL  url.URL
}

type ImageError struct {
	Err     error
	FileURL url.URL
}
