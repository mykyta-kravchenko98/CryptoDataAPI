package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/configs"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/services"
	pb "github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/cryptodata_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCryptoDataServiceServer
	dataService services.DataService
}

// Start server
func Init(dataService services.DataService, configs configs.ServerConfig) error {
	// Создание нового мультиплексора ServeMux
	mux := runtime.NewServeMux()

	// Registrate gRPC-heandler
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterCryptoDataServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf(":%s", configs.GRPCPort), opts)
	if err != nil {
		log.Fatalf("Cant registrate gRPC server: %v", err)
	}

	// Создание обработчика для добавления заголовков CORS
	handler := corsHandler(mux)

	// Запуск HTTP-сервера
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configs.RESTPort), handler))
	}()

	if configs.GRPCPort == "" {
		return errors.New("GRPCPort is empty, can't init gRPC server")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterCryptoDataServiceServer(s, &grpcServer{dataService: dataService})

	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *grpcServer) GetTop50Coins(req *pb.GetCryptoCoinsRequest, stream pb.CryptoDataService_GetTop50CoinsServer) error {
	ctx := stream.Context()

	timerChan := time.After(0)

	for {
		select {
		case <-ctx.Done():
			//client disconected or error happens
			return ctx.Err()

		case <-timerChan:
			coins, err := s.dataService.GetTop50CoinMarketCurrencyProto()
			if err != nil {
				return err
			}

			err = stream.Send(&pb.GetCryptoCoinsResponse{Coins: coins, Count: 50})

			//next interval set
			timerChan = time.After(time.Minute * 1)
		}
	}
}

// For CORS handling
func corsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// any domain access
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// allow methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// allow headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Heandle OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// send request to next heandler
		handler.ServeHTTP(w, r)
	})
}
