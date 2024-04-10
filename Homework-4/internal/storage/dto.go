package storage

import (
	"homework-4/internal/model"
	"time"
)

type DataDTO struct {
	ID              int
	IdClient        int
	IdPVZ           int
	Weight          float64
	TypeOfPackaging model.PackageType
	Price           float64
	DateStorage     time.Time //date of storage for the PVZ
	DateIssue       time.Time //date of issue to the client
	IsReturn        bool      //return status
	IsIssued        bool      //order issue status
	IsIssuedBack    bool
	IsDelete        bool
}
