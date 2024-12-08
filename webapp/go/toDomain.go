package main

import (
	"github.com/isucon/isucon14/webapp/go/sqlcgen"
)

func NewRides(rides []sqlcgen.Ride) []Ride {
	domainRides := make([]Ride, len(rides))
	for i, ride := range rides {
		var evaluation *int
		if ride.Evaluation.Valid {
			evaluationNum := int(ride.Evaluation.Int32)
			evaluation = &evaluationNum
		}
		domainRides[i] = Ride{
			ID:                   ride.ID,
			UserID:               ride.UserID,
			ChairID:              ride.ChairID,
			PickupLatitude:       int(ride.PickupLatitude),
			PickupLongitude:      int(ride.PickupLongitude),
			DestinationLatitude:  int(ride.DestinationLatitude),
			DestinationLongitude: int(ride.DestinationLongitude),
			Evaluation:           evaluation,
			CreatedAt:            ride.CreatedAt,
			UpdatedAt:            ride.UpdatedAt,
		}
	}
	return domainRides
}

func NewCoupon(coupon sqlcgen.Coupon) Coupon {
	var usedBy *string
	if coupon.UsedBy.Valid {
		usedBy = &coupon.UsedBy.String
	}
	return Coupon{
		UserID:    coupon.UserID,
		Code:      coupon.Code,
		Discount:  int(coupon.Discount),
		CreatedAt: coupon.CreatedAt,
		UsedBy:    usedBy,
	}
}

func NewChair(chair sqlcgen.Chair) Chair {
	return Chair{
		ID:          chair.ID,
		OwnerID:     chair.OwnerID,
		Name:        chair.Name,
		Model:       chair.Model,
		IsActive:    chair.IsActive,
		AccessToken: chair.AccessToken,
		CreatedAt:   chair.CreatedAt,
		UpdatedAt:   chair.UpdatedAt,
	}
}

func NewOwner(owner sqlcgen.Owner) Owner {
	return Owner{
		ID:                 owner.ID,
		Name:               owner.Name,
		AccessToken:        owner.AccessToken,
		ChairRegisterToken: owner.ChairRegisterToken,
		CreatedAt:          owner.CreatedAt,
		UpdatedAt:          owner.UpdatedAt,
	}
}
