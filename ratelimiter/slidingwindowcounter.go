package ratelimiter

import (
	"sync"
	"time"
)

type windowCounts struct{
	currentWindow int64
	currentCount int
	previousCount int
}

type SlidingWindowCounterLimiter struct{
	mu sync.Mutex
	windows map[string]*windowCounts
	windowSize time.Duration
	threshold int
}

func NewSlidingWindowCounterLimiter(windowSize time.Duration, threshold int) *SlidingWindowCounterLimiter  {
	return &SlidingWindowCounterLimiter{
		windows: make(map[string]*windowCounts),
		windowSize: windowSize,
		threshold: threshold,
	}
}

func (l *SlidingWindowCounterLimiter) Allow(key string) bool  {
	l.mu.Lock()
	defer l.mu.Unlock()

	windowSecs := int64(l.windowSize.Seconds())
	now := time.Now()
	nowWindow := now.Unix() / windowSecs

	wc, exists := l.windows[key]
	if !exists {
		wc = &windowCounts{currentWindow: nowWindow}
		l.windows[key] = wc
	}

	switch {

	case wc.currentWindow == nowWindow:
		//same window, nothing change
	case wc.currentWindow == nowWindow - 1:
		wc.previousCount, wc.currentCount = wc.currentCount, 0
		wc.currentWindow = nowWindow
	default:
		//gap of more than one window -  both counts are stable
		wc.previousCount, wc.currentCount = 0, 0
		wc.currentCount = int(nowWindow)
	}

	elapsed := float64(now.Unix()%windowSecs) / float64(windowSecs)
	estimated := float64(wc.previousCount)*(1-elapsed) + float64(wc.currentCount)

	if estimated >= float64(l.threshold){

		return false
	}
	wc.currentCount++
	return true

}