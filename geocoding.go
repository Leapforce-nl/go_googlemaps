package googlemaps

import (
	"net/http"

	"github.com/dghubble/sling"
)

// GeoCodes represents a GoogleMaps GeoCode set
type GeoCodes struct {
	Results []GeoCode `json:"results"`
	Status  string    `json:"status"`
}

type GeoCode struct {
	AddressComponents string             `json:"address_components"`
	FormattedAddress  []AddressComponent `json:"formatted_address"`
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
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Viewport struct {
	NorthEast Location `json:"northeast"`
	SouthWest Location `json:"southwest"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

// GeoCodingService provides methods for accessing GoogleMaps geocode endpoints.
type GeoCodingService struct {
	sling *sling.Sling
}

// newGeoCodingService returns a new GeoCodingService.
func newGeoCodingService(sling *sling.Sling) *GeoCodingService {
	return &GeoCodingService{
		sling: sling.Path("geocode/"),
	}
}

// GeoCodeParams are the parameters for GeoCodingService.GeoCode
type GeoCodeParams struct {
	Key     string `url:"key,omitempty"`
	Address string `url:"address,omitempty"`
}

func (s *GeoCodingService) GeoCode(params *GeoCodeParams) ([]GeoCode, *http.Response, error) {
	geoCodes := new(GeoCodes)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("json").QueryStruct(params).Receive(geoCodes, apiError)
	return geoCodes.Results, resp, relevantError(err, *apiError)
}
