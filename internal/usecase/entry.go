package usecase

import (
	"context"
	"fmt"

	"Go_Team00.ID_376234-Team_TL_barievel/internal/entities"
)

type AnomalyRepository interface {
	PutAnomaly(ctx context.Context, msg entities.Entry) error
}

type AnomalyDetector interface {
	ProcessFrequency(freq float64) bool
}

type EntryUsecase struct {
	anomalyDetector   AnomalyDetector
	anomalyRepository AnomalyRepository
}

func NewEntryUsecase(anomalyDetector AnomalyDetector, anomalyRepository AnomalyRepository) *EntryUsecase {
	return &EntryUsecase{anomalyDetector: anomalyDetector, anomalyRepository: anomalyRepository}
}

func (s *EntryUsecase) ProcessEntry(ctx context.Context, entry entities.Entry) error {
	const op = "usecase.EntryUsecase.ProcessEntry"
	isAnomaly := s.anomalyDetector.ProcessFrequency(entry.Frequency)
	if isAnomaly {
		if err := s.anomalyRepository.PutAnomaly(ctx, entry); err != nil {
			return fmt.Errorf("%s: error putting anomaly: %w", op, err)
		}
	}
	return nil
}
