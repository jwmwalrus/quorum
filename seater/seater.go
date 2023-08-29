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

// SeedHandler provides the handler
type SeedHandler struct {
	db   *gorm.DB
	list []Seed
}

// SeedHandlerNew returns a new handler
func SeedHandlerNew(db *gorm.DB) (h SeedHandler) {
	h.db = db
	return
}

// Add adds a new seed
func (h *SeedHandler) Add(s Seed) {
	h.list = append(h.list, s)
}

// AddSome adds a bunch of seeds
func (h *SeedHandler) AddSome(l []Seed) {
	for _, l := range l {
		h.Add(l)
	}
}

// RunAll runs all seeds
func (h *SeedHandler) RunAll() (err error) {
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

// Runs a single seed by name
func (h *SeedHandler) RunByName(name string) (err error) {
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
