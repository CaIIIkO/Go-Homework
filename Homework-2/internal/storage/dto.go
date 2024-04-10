package storage

import (
	"time"
)

type DataDTO struct {
	ID       int
	IdClient int
	IdPVZ    int
	//Discription string
	DateStorage  time.Time //date of storage for the PVZ
	DateIssue    time.Time //date of issue to the client
	IsReturn     bool      //return status
	IsIssued     bool      //order issue status
	IsIssuedBack bool
	IsDelete     bool
}

type PvzDTO struct {
	ID      int
	Name    string
	Address string
	Contact string
}
