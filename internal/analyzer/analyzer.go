package analyzer

import (
	"log"
	"math"
	"sync"
	"time"
)

type Analyzer struct {
	count              int
	mean               float64
	m2                 float64
	stdDev             float64
	ready              bool
	mu                 sync.Mutex
	anomalyCoefficient float64
	lastLog            time.Time
	logFreq            time.Duration
}

func NewAnalyzer(
	anomalyCoefficient float64,
	logFrequency time.Duration,
) *Analyzer {
	anomalyAnalyzer := &Analyzer{
		anomalyCoefficient: anomalyCoefficient,
		logFreq:            logFrequency,
	}
	return anomalyAnalyzer
}

func (a *Analyzer) ProcessFrequency(freq float64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.count++
	delta := freq - a.mean
	a.mean += delta / float64(a.count)
	delta2 := freq - a.mean
	a.m2 += delta * delta2

	if a.count >= 2 {
		variance := a.m2 / float64(a.count-1)
		a.stdDev = math.Sqrt(variance)
	}
	now := time.Now()
	if now.Sub(a.lastLog) >= a.logFreq {
		log.Printf("Обработано значений: %d. Среднее=%.4f, STD=%.4f",
			a.count, a.mean, a.stdDev)
		a.lastLog = now
	}

	if a.count >= 100 && !a.ready {
		log.Println("Переключение в режим детектирования аномалий")
		a.ready = true
	}

	if a.ready && a.isAnomaly(freq) {
		log.Printf("ОБНАРУЖЕНА АНОМАЛИЯ: %.4f (среднее=%.4f, отклонение=%.4f)", freq, a.mean, a.stdDev)
		return true
	}
	return false
}

func (a *Analyzer) isAnomaly(value float64) bool {
	if !a.ready || a.stdDev == 0 {
		return false
	}
	lowerBound := a.mean - a.anomalyCoefficient*a.stdDev
	upperBound := a.mean + a.anomalyCoefficient*a.stdDev
	return value < lowerBound || value > upperBound
}
