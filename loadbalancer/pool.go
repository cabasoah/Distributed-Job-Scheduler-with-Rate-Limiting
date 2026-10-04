package loadbalancer

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL *url.URL
	Alive atomic.Bool
}

type Pool struct {
	backends []*Backend
	mu sync.Mutex
	current int
}

func NewPool(urls []string) (*Pool, error)  {
	p := &Pool{}
	for _, raw := range urls {
		u, err := url.Parse(raw)

		if err != nil {
			return  nil, err
		}
		b := &Backend{URL: u}
		b.Alive.Store(true) //assume healthy until the first check says otherwise
		p.backends = append(p.backends, b)
	}
	return p, nil
}

// NextBackend returns the next healthy backend, round robin.
func (p *Pool) NextBackend() *Backend {
	p.mu.Lock()
	defer p.mu.Unlock()

	n := len(p.backends)
	for i := 0; i < n; i++ {
		b := p.backends[p.current%n]
		p.current++

		if b.Alive.Load(){
			return b
		}
	}
	return  nil //none healthy
}


//StartHealthChecks

func (p *Pool) StartHealthChecks(path string, period time.Duration)  {
	ticker := time.NewTicker(period)

	//gorountiine
	go func() {
		for range ticker.C{
			for _, b:= range p.backends {
				go p.checkBackend(b, path)
			}
		}

	}()
}


func (p *Pool) checkBackend(b *Backend, path string)  {
	resp, err := http.Get(b.URL.String() + path)
	wasAlive := b.Alive.Load()
	nowAlive := err == nil && resp.StatusCode == http.StatusOK

	if err == nil {
		resp.Body.Close()
	}
	b.Alive.Store(nowAlive)

	if wasAlive && !nowAlive {
		log.Printf("backend %s failed health check, taking out of rotation", b.URL)

	} else if !wasAlive && nowAlive {
		log.Printf("backend %s passed health check, back in rotation", b.URL)
	}
}