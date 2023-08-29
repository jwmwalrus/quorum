package seater

import (
	"fmt"

	"gorm.io/gorm"
)

// Seed provides the basic struct
type Seed struct {
	Name     string
	Run      func(*gorm.DB) error
	Requires []string
	done     bool
}

// Seater provides the handler
type Seater interface {
	// Add adds a new seed
	Add(s Seed)

	// AddSome adds a bunch of seeds
	AddSome(l []Seed)

	// RunAll runs all seeds
	RunAll() error

	// Runs a single seed by name
	RunByName(name string) error
}

// New returns a new handler
func New(db *gorm.DB) Seater {
	return &seedHandler{db: db}
}

type seedHandler struct {
	db   *gorm.DB
	list []Seed
}

// Add implements the Seater interface
func (h *seedHandler) Add(s Seed) {
	h.list = append(h.list, s)
}

// AddSome implements the Seater interface
func (h *seedHandler) AddSome(l []Seed) {
	for _, l := range l {
		h.Add(l)
	}
}

// RunAll implements the Seater interface
func (h *seedHandler) RunAll() (err error) {
	for k, v := range h.list {
		if v.done {
			continue
		}
		for _, r := range v.Requires {
			if err = h.RunByName(r); err != nil {
				return
			}
		}
		if err = v.Run(h.db); err != nil {
			return
		}
		h.list[k].done = true
	}
	return
}

// RunByName implements the Seater interface
func (h *seedHandler) RunByName(name string) (err error) {
	for k, v := range h.list {
		if v.Name != name {
			continue
		}
		if v.done {
			return
		}
		err = v.Run(h.db)
		h.list[k].done = true
		return
	}

	return fmt.Errorf("Seed not found: %v", name)
}
