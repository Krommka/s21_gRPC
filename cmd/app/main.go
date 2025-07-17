package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Go_Team00.ID_376234-Team_TL_barievel/configs"
	"Go_Team00.ID_376234-Team_TL_barievel/db/postgres"
	"Go_Team00.ID_376234-Team_TL_barievel/internal/analyzer"
	"Go_Team00.ID_376234-Team_TL_barievel/internal/client"
	"Go_Team00.ID_376234-Team_TL_barievel/internal/server"
	"Go_Team00.ID_376234-Team_TL_barievel/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	store, err := postgres.NewStore(dbCtx, *cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to database")

	k := flag.Float64("k", 3.0, "Коэффициент аномалий (сколько STD от среднего)")
	flag.Parse()

	if *k <= 0 {
		log.Fatal("Коэффициент должен быть положительным")
	}

	anomalyDetector := analyzer.NewAnalyzer(*k, cfg.Anomaly.LogFrequency)

	uc := usecase.NewEntryUsecase(anomalyDetector, store)

	grpcClient := client.NewClient(cfg, uc)
	defer grpcClient.Close()

	grpcServer := server.NewServer(cfg)
	defer grpcServer.Shutdown()

	go func() {
		log.Printf("Starting grpc server on port %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Println("Starting grpc client")
	clientCtx, clientCancel := context.WithCancel(context.Background())
	go grpcClient.RunWithReconnect(clientCtx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Println("Shutting down")

	clientCancel()
	grpcClient.Close()
	grpcServer.Shutdown()

	dbShutdownCtx, dbShutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbShutdownCancel()
	store.Disconnect(dbShutdownCtx)

	log.Println("Shutdown complete")

}
