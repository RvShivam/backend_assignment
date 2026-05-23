package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var errDuplicateSKU = errors.New("sku already exists")

type Store struct {
	mu       sync.RWMutex
	products map[string]*Product
	skuIndex map[string]string
	ordered  []string
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]*Product),
		skuIndex: make(map[string]string),
	}
}

func newID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// copyProduct returns a deep copy so callers can safely use it after releasing the lock.
func copyProduct(p *Product) Product {
	cp := *p
	cp.ImageURLs = make([]string, len(p.ImageURLs))
	copy(cp.ImageURLs, p.ImageURLs)
	cp.VideoURLs = make([]string, len(p.VideoURLs))
	copy(cp.VideoURLs, p.VideoURLs)
	return cp
}

func (s *Store) Create(req CreateProductRequest) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.skuIndex[req.SKU]; exists {
		return Product{}, errDuplicateSKU
	}

	imageURLs := req.ImageURLs
	if imageURLs == nil {
		imageURLs = []string{}
	}
	videoURLs := req.VideoURLs
	if videoURLs == nil {
		videoURLs = []string{}
	}

	p := &Product{
		ID:        newID(),
		Name:      req.Name,
		SKU:       req.SKU,
		ImageURLs: imageURLs,
		VideoURLs: videoURLs,
		CreatedAt: time.Now(),
	}
	s.products[p.ID] = p
	s.skuIndex[p.SKU] = p.ID
	s.ordered = append(s.ordered, p.ID)

	return copyProduct(p), nil
}

func (s *Store) List(limit, offset int) ([]ProductSummary, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.ordered)
	if offset >= total {
		return []ProductSummary{}, total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	summaries := make([]ProductSummary, 0, end-offset)
	for _, id := range s.ordered[offset:end] {
		p := s.products[id]
		thumb := ""
		if len(p.ImageURLs) > 0 {
			thumb = p.ImageURLs[0]
		}
		summaries = append(summaries, ProductSummary{
			ID:           p.ID,
			Name:         p.Name,
			SKU:          p.SKU,
			ImageCount:   len(p.ImageURLs),
			VideoCount:   len(p.VideoURLs),
			ThumbnailURL: thumb,
			CreatedAt:    p.CreatedAt,
		})
	}

	return summaries, total
}

func (s *Store) Get(id string) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok {
		return Product{}, false
	}
	return copyProduct(p), true
}

func (s *Store) AddMedia(id string, imageURLs, videoURLs []string) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[id]
	if !ok {
		return Product{}, false
	}
	p.ImageURLs = append(p.ImageURLs, imageURLs...)
	p.VideoURLs = append(p.VideoURLs, videoURLs...)
	return copyProduct(p), true
}
