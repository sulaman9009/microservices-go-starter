package domain

import driverv1 "ride-sharing/shared/gen/go/driver/v1"

type DriverService interface {
	RegisterDriver(driverId string, packageSlug string) (*driverv1.Driver, error)
	UnregisterDriver(driverId string)
}
