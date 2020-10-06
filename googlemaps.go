package googlemaps

import (
	"net/http"

	"github.com/dghubble/sling"
)

const (
	apiURL string = "https://maps.googleapis.com/maps/api/"
)

// GoogleMaps stores GoogleMaps configuration
//
type GoogleMaps struct {
	sling     *sling.Sling
	GeoCoding *GeoCodingService
}

// NewGoogleMaps returns a pointer to a new GoogleMaps instance
//
func NewGoogleMaps(httpClient *http.Client, isLive bool) *GoogleMaps {
	base := sling.New().Client(httpClient).Base(apiURL)
	return &GoogleMaps{
		sling:     base,
		GeoCoding: newGeoCodingService(base.New()),
	}
}
