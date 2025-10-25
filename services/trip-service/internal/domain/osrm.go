package domain

import tripv1 "ride-sharing/shared/gen/go/trip/v1"

type OsrmApiResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OsrmApiResponse) ToProto() *tripv1.Route {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([]*tripv1.Coordinate, len(geometry))
	for i, coord := range geometry {
		coordinates[i] = &tripv1.Coordinate{
			Latitude:  coord[0],
			Longitude: coord[1],
		}
	}

	return &tripv1.Route{
		Geometry: []*tripv1.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}
