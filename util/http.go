package util

import (
	"fmt"
	"mime"
	"net/http"
)

func ValidateHttpResponse(resp *http.Response, requiredType *string) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Failed to get %s with status %s: %+v", resp.Request.URL.String(), resp.Status, resp)
	}
	if requiredType != nil {
		mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			return err
		}
		if mediatype != *requiredType {
			return fmt.Errorf("got invalid media type %s for page %s", mediatype, resp.Request.URL.String())
		}
	}
	return nil
}
