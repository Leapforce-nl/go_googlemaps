package googlemaps

import (
	"net/http"

	"github.com/dghubble/sling"
	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	apiURL string = "https://maps.googleapis.com/maps/api/"
)

type ServiceConfig struct {
	GeoCodingAPIKey *string
}

// Service stores Service configuration
//
type Service struct {
	sling            *sling.Sling
	GeoCodingService *GeoCodingService
}

// NewGoogleMaps returns a pointer to a new Service instance
//
func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	client := http.Client{}
	base := sling.New().Client(&client).Base(apiURL)

	service := Service{
		sling: base,
	}

	if serviceConfig.GeoCodingAPIKey != nil {
		geoCodingServiceConfig := geoCodingServiceConfig{
			APIKey: *serviceConfig.GeoCodingAPIKey,
			Sling:  base.New(),
		}

		geoCodingService, e := newGeoCodingService(&geoCodingServiceConfig)
		if e != nil {
			return nil, e
		}

		service.GeoCodingService = geoCodingService
	}

	return &service, nil
}
