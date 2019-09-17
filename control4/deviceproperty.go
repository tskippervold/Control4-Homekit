package control4

import (
	"net/http"
	"strings"
)

type Property string

const (
	Power      Property = "1000"
	Brightness Property = "1001"
)

func PropertyFrom(r *http.Request) (Property, string) {
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	return Property(parts[1]), parts[2]
}
