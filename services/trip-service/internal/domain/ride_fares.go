package domain

import (
	tripv1 "ride-sharing/shared/gen/go/trip/v1"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string // ex: van, luxury, sedan
	TotalPriceInCents float64
	ExpiresAt         int64
}

func (r *RideFareModel) ToProto() *tripv1.RideFare {
	return &tripv1.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}
func ToRideFaresProto(fares []*RideFareModel) []*tripv1.RideFare {
	var protoFares []*tripv1.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	return protoFares
}
