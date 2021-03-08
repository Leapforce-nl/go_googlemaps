package googlemaps

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/dghubble/sling"
	errortools "github.com/leapforce-libraries/go_errortools"
)

// GeoCodes represents a GoogleMaps GeoCode set
type GeoCodes struct {
	Results []GeoCode `json:"results"`
	Status  string    `json:"status"`
}

type GeoCode struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	PlusCode          PlusCode           `json:"plus_code"`
	Types             []string           `json:"types"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Viewport `json:"viewport"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Viewport struct {
	NorthEast Location `json:"northeast"`
	SouthWest Location `json:"southwest"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

type geoCodingServiceConfig struct {
	APIKey string
	Sling  *sling.Sling
}

// GeoCodingService provides methods for accessing GoogleMaps geocode endpoints.
type GeoCodingService struct {
	apiKey string
	sling  *sling.Sling
}

// newGeoCodingService returns a new GeoCodingService.
func newGeoCodingService(geoCodingServiceConfig *geoCodingServiceConfig) (*GeoCodingService, *errortools.Error) {
	if geoCodingServiceConfig.APIKey == "" {
		return nil, errortools.ErrorMessage("APIKey not provided")
	}

	return &GeoCodingService{
		apiKey: geoCodingServiceConfig.APIKey,
		sling:  geoCodingServiceConfig.Sling.Path("geocode/"),
	}, nil
}

// GeoCodeParams are the parameters for GeoCodingService.GeoCode
type GeoCodeParams struct {
	Key     string `url:"key,omitempty"`
	Address string `url:"address,omitempty"`
}

func (s *GeoCodingService) GeoCode(params *GeoCodeParams) (*[]GeoCode, *errortools.Error) {

	p := struct {
		Key     string `url:"key,omitempty"`
		Address string `url:"address,omitempty"`
	}{
		s.apiKey,
		params.Address,
	}

	geoCodes := new(GeoCodes)
	errorResponse := new(ErrorResponse)

	retry := 0
	maxRetries := 5

	for retry <= maxRetries {
		if retry > 0 {
			fmt.Printf("Starting retry %v for GeoCode `%s`\n", retry, params.Address)
			waitSeconds := math.Pow(2, float64(retry-1))
			waitMilliseconds := int(rand.Float64() * 1000)
			time.Sleep(time.Duration(waitSeconds)*time.Second + time.Duration(waitMilliseconds)*time.Millisecond)
		}

		response, err := s.sling.New().Get("json").QueryStruct(p).Receive(geoCodes, errorResponse)

		statusCode := 0
		if response != nil {
			statusCode = response.StatusCode
		}

		if (statusCode == 500 || statusCode == 503) && retry < maxRetries { // retry in case of status 500/503 (server error)
			retry++
		} else {
			if err == nil && (statusCode/100 == 4 || statusCode/100 == 5) {
				err = fmt.Errorf("Server returned statuscode %v", statusCode)
			}

			if err != nil {
				e := new(errortools.Error)
				e.SetResponse(response)
				e.SetMessage(err.Error())

				return nil, e
			}

			return &geoCodes.Results, nil
		}
	}

	// should never reach this
	return nil, nil
}
