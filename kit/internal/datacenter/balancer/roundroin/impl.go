package roundroin

import (
	"sync"
	"sync/atomic"

	"github.com/KingTrack/gin-kit/kit/internal/datacenter/balancer/define"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
	"github.com/pkg/errors"
)

type Balancer struct {
	instances []instance.Instance
	mu        sync.RWMutex
	counter   int64
}

func New() define.IBalancer {
	return &Balancer{}
}

func (b *Balancer) Pick(meta map[string]string, skip *instance.Instance) (*instance.Instance, error) {
	if len(b.instances) == 0 {
		return nil, errors.New("no available instances")
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	// 在读锁保护下，b.instances不会被其他goroutine修改
	candidates := make([]instance.Instance, 0, len(b.instances))
	if skip != nil {
		for _, v := range b.instances {
			if v.IsEqualEndpoint(skip) {
				continue
			}
			candidates = append(candidates, v)
		}
	} else {
		candidates = b.instances
	}

	if len(candidates) == 0 {
		return nil, errors.New("no available instances when retry")
	}

	fallbackCandidates := make([]instance.Instance, 0, len(candidates))
	fallbackWeight := 0
	matchedCandidates := make([]instance.Instance, 0, len(candidates))
	matchedWeight := 0
	for _, v := range candidates {
		if v.Weight <= 0 {
			continue
		}
		fallbackCandidates = append(fallbackCandidates, v)
		fallbackWeight += v.Weight
		if len(meta) > 0 && !instance.MatchMeta(meta, v.Meta) {
			continue
		}
		matchedCandidates = append(matchedCandidates, v)
		matchedWeight += v.Weight
	}
	if fallbackWeight == 0 {
		return nil, errors.New("no instances available with positive weight")
	}

	if matchedWeight == 0 {
		matchedWeight = fallbackWeight
		matchedCandidates = fallbackCandidates
	}

	// 加权轮询
	counter := atomic.AddInt64(&b.counter, 1)
	weightedIndex := int(counter) % matchedWeight
	currentWeight := 0
	for _, v := range matchedCandidates {
		currentWeight += v.Weight
		if weightedIndex < currentWeight {
			return &v, nil
		}
	}

	// 理论上不会执行到这
	return &matchedCandidates[0], nil
}

func (b *Balancer) Update(instances []instance.Instance) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.instances = instances
}
