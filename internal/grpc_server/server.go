package grpcserver

import (
	"errors"
	"fmt"
	"net"
	"time"

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
	if configs.GRPCPort == "" {
		return errors.New("GRPCPort is empty, can`t init grpcServer")
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
			coins, err := s.dataService.GetTop50CoinMarketCurrency()
			if err != nil {
				return err
			}

			err = stream.Send(&pb.GetCryptoCoinsResponse{Coins: coins, Count: 50})

			//next interval set
			timerChan = time.After(time.Minute * 1)
		}
	}
}
