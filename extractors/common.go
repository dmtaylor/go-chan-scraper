package extractors

import "net/url"

type ImageFile struct {
	FileName string
	FileUrl  url.URL
}

type ImageError struct {
	Err     error
	FileUrl url.URL
}
