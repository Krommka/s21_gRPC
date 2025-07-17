package server

import (
	"log"
	"math/rand"
	"time"

	"Go_Team00.ID_376234-Team_TL_barievel/api/gen/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	pb.UnimplementedFrequencyServiceServer
}

func newHandler() *handler {
	return &handler{}
}

// TransmitFrequencies реализует стриминг частот с нормальным распределением
func (h *handler) TransmitFrequencies(in *emptypb.Empty, stream grpc.ServerStreamingServer[pb.FrequencyMessage]) error {
	sessionID := uuid.New().String()

	// Генерируем случайные параметры распределения
	// Мат. ожидание из интервала [-10, 10]
	mean := rand.Float64()*20 - 10

	// Стандартное отклонение из интервала [0.3, 1.5]
	stdDev := rand.Float64()*1.2 + 0.3

	log.Printf("New connection established. Session ID: %s, Mean: %.6f, StdDev: %.6f",
		sessionID, mean, stdDev)

	// Создаем генератор нормального распределения
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Client disconnected. Session ID: %s", sessionID)
			return status.FromContextError(stream.Context().Err()).Err()
		default:
			// Генерируем частоту по нормальному распределению
			// Используем Box-Muller трансформацию для генерации нормально распределенных значений
			frequency := rng.NormFloat64()*stdDev + mean

			msg := &pb.FrequencyMessage{
				SessionId: sessionID,
				Frequency: frequency,
				Timestamp: timestamppb.Now(),
			}

			if err := stream.Send(msg); err != nil {
				log.Printf("Error sending message for session %s: %v", sessionID, err)

				if s, ok := status.FromError(err); ok {
					// Если это уже gRPC ошибка, возвращаем как есть
					return s.Err()
				}

				if stream.Context().Err() != nil {
					// Если контекст отменен, значит проблема с клиентом
					return status.FromContextError(stream.Context().Err()).Err()
				}

				// В остальных случаях считаем это сетевой проблемой
				return status.Errorf(codes.Unavailable, "failed to send message: %v", err)
			}

			// Небольшая задержка между сообщениями (100ms)
			// Это предотвращает чрезмерную нагрузку и делает стрим более реалистичным
			time.Sleep(100 * time.Millisecond)
		}
	}
}
