package services

import (
	math "math/rand/v2"
	"slices"

	"ride-sharing/services/driver-service/internal/utils"
	driverv1 "ride-sharing/shared/gen/go/driver/v1"
	"ride-sharing/shared/util"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type driverInMap struct {
	Driver driverv1.Driver
}

type driverService struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewDriverService() *driverService {
	return &driverService{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *driverService) RegisterDriver(driverId string, packageSlug string) (*driverv1.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(utils.PredefinedRoutes))
	randomRoute := utils.PredefinedRoutes[randomIndex]

	randomPlate := utils.GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)

	// we can ignore this property for now, but it must be sent to the frontend.
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &driverv1.Driver{
		Id:             driverId,
		Geohash:        geohash,
		Location:       &driverv1.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lando Norris",
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
	}

	s.drivers = append(s.drivers, &driverInMap{Driver: *driver})

	return driver, nil
}

func (s *driverService) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.drivers = slices.DeleteFunc(s.drivers, func(d *driverInMap) bool {
		return d.Driver.Id == driverId
	})
}
