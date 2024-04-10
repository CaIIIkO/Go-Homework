package model

import (
	"errors"
	"time"
)

type DataInputAccept struct {
	ID              int
	IdPVZ           int
	DateStorage     time.Time
	IdClient        int
	Price           float64
	Weight          float64
	TypeOfPackaging PackageType
}

type PackageType string

const PACKET PackageType = "packet"
const BOX PackageType = "box"
const FILM PackageType = "film"

func GetTypePackage(packageType PackageType) (Package, error) {
	switch packageType {
	case PACKET:
		return PackageS{maxWeight: 10, price: 5}, nil
	case BOX:
		return PackageS{maxWeight: 30, price: 20}, nil
	case FILM:
		return PackageS{maxWeight: 0, price: 1}, nil
	default:
		return nil, errors.New("ошибка! выбранной упаковки не существует")
	}
}

type Package interface {
	MaxWeight() float64
	Price() float64
	Validate(float64) bool
}

type PackageS struct {
	maxWeight float64
	price     float64
}

func (p PackageS) MaxWeight() float64 {
	return p.maxWeight
}
func (p PackageS) Price() float64 {
	return p.price
}
func (p PackageS) Validate(weight float64) bool {
	if weight > p.MaxWeight() && p.MaxWeight() != 0 {
		return false
	}
	return true
}
