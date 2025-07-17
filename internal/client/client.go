package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"Go_Team00.ID_376234-Team_TL_barievel/api/gen/pb"
	"Go_Team00.ID_376234-Team_TL_barievel/configs"
	"Go_Team00.ID_376234-Team_TL_barievel/internal/entities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usecase interface {
	ProcessEntry(ctx context.Context, entry entities.Entry) error
}

type Client struct {
	conn           *grpc.ClientConn
	client         pb.FrequencyServiceClient
	serverPort     string
	reconnectDelay time.Duration
	uc             Usecase
}

func NewClient(cfg *configs.Config, uc Usecase) *Client {
	return &Client{
		serverPort:     cfg.GRPC.Port,
		reconnectDelay: cfg.GRPC.ReconnectDelay,
		uc:             uc,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	conn, err := grpc.NewClient(fmt.Sprintf(":%s", c.serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = pb.NewFrequencyServiceClient(conn)

	log.Printf("Connected to server at %s", fmt.Sprintf(":%s", c.serverPort))
	return nil
}

func (c *Client) StartReceiving(ctx context.Context) error {
	const op = "client.StartReceiving"
	if c.client == nil {
		return fmt.Errorf("%s: client not connected", op)
	}

	stream, err := c.client.TransmitFrequencies(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Println("Started receiving frequency stream...")

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s: context cancelled, stopping stream reception", op)
			return status.FromContextError(ctx.Err()).Err()
		default:
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Printf("%s: server closed the stream", op)
					return nil
				}

				if s, ok := status.FromError(err); ok {
					switch s.Code() {
					case codes.Canceled:
						log.Printf("%s: stream was canceled", op)
						return err
					case codes.DeadlineExceeded:
						log.Printf("%s: stream deadline exceeded", op)
						return err
					case codes.Unavailable:
						log.Printf("%s: server unavailable: %v", op, err)
						return err
					default:
						log.Printf("%s: stream error: %v", op, err)
						return err
					}
				}

				log.Printf("%s: unknown stream error: %v", op, err)
				return err
			}

			// Обрабатываем полученное сообщение
			err = c.uc.ProcessEntry(ctx, entities.Entry{
				SessionId: msg.SessionId,
				Frequency: msg.Frequency,
				Timestamp: msg.Timestamp.AsTime(),
			})
			if err != nil {
				log.Printf("%s: error processing entry: %v", op, err)
			}
		}
	}
}

// RunWithReconnect запускает клиент с автоматическим переподключением
func (c *Client) RunWithReconnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Client shutdown requested")
			return
		default:
			// Пытаемся подключиться
			if err := c.Connect(ctx); err != nil {
				log.Printf("Failed to connect: %v. Retrying in %v...", err, c.reconnectDelay)
				time.Sleep(c.reconnectDelay)
				continue
			}

			// Начинаем получение данных
			if err := c.StartReceiving(ctx); err != nil {
				log.Printf("Stream error: %v. Reconnecting in %v...", err, c.reconnectDelay)
				c.Close() // Закрываем соединение перед переподключением
				time.Sleep(c.reconnectDelay)
				continue
			}
		}
	}
}

// Close закрывает соединение с сервером
func (c *Client) Close() error {
	if c.conn != nil {
		log.Println("Closing connection to server")
		return c.conn.Close()
	}
	return nil
}
