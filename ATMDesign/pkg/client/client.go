package client

import (
	"context"
	"github.com/DesignQuestions-Go/ATMDesign/pkg/pb"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
	//"encoding/json"
)

const (
	//	server = "10.196.105.125:50051"
	server = "localhost:50051"
)

type ClientConn interface {
	Close() error
}

func Connect(serverAddr string) (*grpc.ClientConn, error) {
	log.Infof("trying to connect to %s", serverAddr)
	dialOpt := grpc.WithInsecure()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, serverAddr,
		dialOpt,
		grpc.WithBlock(),
	)

	if err != nil {
		log.Errorf("unable to create connection %v", err)
		return nil, err
	}
	log.Infof("Established connection with %s", serverAddr)
	return conn, nil
}

func NewServer() (pb.BankServiceClient, ClientConn, error) {
	conn, err := Connect(server)
	if err != nil {
		return nil, conn, err
	}
	return pb.NewBankServiceClient(conn), conn, nil
}

func Register(ctx context.Context, client pb.BankServiceClient, name string, pin int32, money int32) error {
	resp, err := client.RegisterTBank(ctx, &pb.RegisterRequest{Name: name, Pin: pin, Money: money})
	if err != nil {
		log.Errorf("Failed to register to bank %v", err)
		return err
	}
	log.Infof("the resp is %v ", resp)
	return nil
}

func WithDrawCash(ctx context.Context, client pb.BankServiceClient, name string, pin int32, money int32) {
	reqMap := make(map[string]int32)
	reqMap[name] = pin
	res, err := client.Authenticate(ctx, &pb.AuthenticateRequest{Req: reqMap})
	if err != nil {
		log.Errorf("\nFailed to authenticate:%v\n", err)
	}
	if res != nil && res.Authenticated {
		resp, err := client.WithDraw(ctx, &pb.DepositMoney{Name: name, Money: money})
		if err != nil {
			log.Errorln(err)
		} else {
			log.Infof("your balance is %v\n", resp.Money)
		}
	}
}

func CheckBalance(ctx context.Context, client pb.BankServiceClient, name string, pin int32) {
	reqMap := make(map[string]int32)
	reqMap[name] = pin
	res, err := client.Authenticate(ctx, &pb.AuthenticateRequest{Req: reqMap})
	if err != nil {
		log.Errorf("\nFailed to authenticate:%v\n", err)
	}
	if res != nil && res.Authenticated {
		resp, err := client.BankBalance(ctx, &pb.BankBalanceRequest{Name: name})
		if err != nil {
			log.Errorln(err)
		} else {
			log.Infof("your balance is %v\n", resp.Money)
		}
	}
}

func DepositMoney(ctx context.Context, client pb.BankServiceClient, name string, pin int32, money int32) {
	reqMap := make(map[string]int32)
	reqMap[name] = pin
	res, err := client.Authenticate(ctx, &pb.AuthenticateRequest{Req: reqMap})
	if err != nil {
		log.Errorf("\nFailed to authenticate:%v\n", err)
	}
	if res != nil && res.Authenticated {
		resp, err := client.Deposit(ctx, &pb.DepositMoney{Name: name, Money: money})
		if err != nil {
			log.Errorln(err)
		} else {
			log.Infof("your balance is %v\n", resp.Money)
		}
	}
}
