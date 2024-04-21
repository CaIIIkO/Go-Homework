package fixtures

import (
	"homework-3/internal/pkg/repository"
	"homework-3/tests/states"
)

type PvzBuilder struct {
	instance *repository.Pvz
}

func Pvz() *PvzBuilder {
	return &PvzBuilder{instance: &repository.Pvz{}}
}

func (b *PvzBuilder) ID(v int64) *PvzBuilder {
	b.instance.ID = v
	return b
}

func (b *PvzBuilder) Name(v string) *PvzBuilder {
	b.instance.Name = v
	return b
}

func (b *PvzBuilder) Address(v string) *PvzBuilder {
	b.instance.Address = v
	return b
}

func (b *PvzBuilder) Contact(v string) *PvzBuilder {
	b.instance.Contact = v
	return b
}

func (b *PvzBuilder) P() *repository.Pvz {
	return b.instance
}

func (b *PvzBuilder) V() repository.Pvz {
	return *b.instance
}

func (b *PvzBuilder) Valid() *PvzBuilder {
	return Pvz().ID(states.ID1).Name(states.Name1).Address(states.Address1).Contact(states.Contact1)
}
