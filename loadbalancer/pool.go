package loadbalancer

import(
	"net/url"
	"sync"
)

type Backend struct {
	URL *url.URL
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

		p.backends = append(p.backends, &Backend{URL: u})
	}
	return p, nil
}

// NextBackend returns the next backend, round robin.
func (p *Pool) NextBackend() *Backend {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.backends) == 0 {
		return nil
	}

	b := p.backends[p.current % len(p.backends)]
	p.current++
	return  b
}