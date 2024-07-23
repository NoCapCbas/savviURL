package inMemory

import (
  "sync"
  "time"
  
	"github.com/pkg/errors"

	"github.com/NoCapCbas/url-shortener/urlshortener"
)

type inMemoryRepository struct {
  store map[string]*urlshortener.Redirect
  mu sync.RWMutex
}

func NewInMemoryRepository() urlshortener.RedirectRepository {

  return &inMemoryRepository{
    store: make(map[string]*urlshortener.Redirect),
  }

}

func (r *inMemoryRepository) Find(code string) (*urlshortener.Redirect, error) {
  r.mu.RLock()
  defer r.mu.RUnlock()

  redirect, exists := r.store(code)
  if !exists {
    return nil, errors.Wrap(urlshortener.ErrRedirectNotFound, "repository.Redirect.Find")
  }
  return redirect, nil
}

func (r *inMemoryRepository) Store(redirect *urlshortener.Redirect) (error) {
  r.mu.Lock()
  defer r.mu.Unlock()
  r.store[redirect.Code] = redirect
  return nil
}
