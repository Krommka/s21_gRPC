package analyzer

import (
	"testing"
	"time"
)

// Тест 1: Проверка корректности среднего и stdDev
func TestProcessValue_Statistics(t *testing.T) {
	an := NewAnalyzer(3, time.Hour)

	values := []float64{1, 2, 3, 4, 5}
	for _, v := range values {
		an.ProcessFrequency(v)
	}

	expectedMean := 3.0
	expectedStdDev := 1.5811 // √2.5 (проверено вручную)

	if diff := an.mean - expectedMean; diff > 0.01 {
		t.Errorf("Ожидалось среднее %.2f, получено %.2f", expectedMean, an.mean)
	}
	if diff := an.stdDev - expectedStdDev; diff > 0.01 {
		t.Errorf("Ожидалось STD %.4f, получено %.4f", expectedStdDev, an.stdDev)
	}
}

// Тест 2: Проверка перехода в режим готовности после 100 значений
func TestProcessValue_Readiness(t *testing.T) {
	an := NewAnalyzer(3, time.Hour)

	for i := 0; i < 99; i++ {
		an.ProcessFrequency(1.0)
	}
	if an.ready {
		t.Error("Анализатор не должен быть готов при 99 значениях")
	}

	an.ProcessFrequency(1.0)
	if !an.ready {
		t.Error("Анализатор должен быть готов после 100 значений")
	}
}

// Тест 3: Проверка, что аномалия корректно распознаётся
func TestProcessValue_AnomalyDetection(t *testing.T) {
	an := NewAnalyzer(2.0, time.Hour)

	// 100 значений вокруг 10
	for i := 0; i < 100; i++ {
		an.ProcessFrequency(10)
	}

	if !an.ready {
		t.Fatal("Анализатор должен быть готов после 100 значений")
	}

	// Вводим выброс
	anomalyDetected := an.ProcessFrequency(50)

	if !anomalyDetected {
		t.Error("Аномалия не была зафиксирована, хотя должна быть")
	}
}
