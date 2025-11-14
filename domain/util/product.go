package util

import "fmt"

func GetProductImageURL(filename string, baseURL string) string {
	return fmt.Sprintf("%s/images/%s", baseURL, filename)
}
