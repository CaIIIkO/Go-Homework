package model

import "time"

type DataInputAccept struct {
	ID          int
	IdPVZ       int
	DateStorage time.Time
	IdClient    int
}
