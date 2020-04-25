package main

import (
	//"github.com/DesignQuestions-Go/ATMDesign/pkg/pb"
	"context"
	"fmt"
	cli "github.com/DesignQuestions-Go/ATMDesign/pkg/client"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"time"
)

var option int
var name string
var money int32
var pin int32

func main() {
	client, conn, err := cli.NewServer()
	if err != nil {
		log.Errorf("NewServer() - error building the grpc client. Reason: %v", err)
		return
	}
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999

	ctx := metadata.AppendToOutgoingContext(context.Background())
atmloop:
	for {
		fmt.Printf("\nThe options available are \n 1.%v\n 2.%v\n 3.%v\n 4.%v\n 5.%v\n", "Register", "WithDrawCash", "CheckBalance", "DepositMoney", "Exit")
		fmt.Println("Enter Your Option ")
		fmt.Scan(&option)
		switch option {
		case 1:
			fmt.Println("Enter you name to register and  money to deposit")
			fmt.Scanf("%s %d \n", &name, &money)
			pin = rand.Int31n(int32(max - min + 1))
			if err := cli.Register(ctx, client, name, pin, money); err == nil {
				log.Infof("your pin is %v \n", pin)
			}
		case 2:
			fmt.Println("Enter you name , pin and cash to withdraw")
			fmt.Scanf("%s %d %d\n", &name, &pin, &money)
			cli.WithDrawCash(ctx, client, name, pin, money)
		case 3:
			fmt.Println("Enter you name , pin ")
			fmt.Scanf("%s %d \n", &name, &pin)
			cli.CheckBalance(ctx, client, name, pin)
		case 4:
			fmt.Println("Enter you name , pin, money ")
			fmt.Scanf("%s %d %d\n", &name, &pin, &money)
			cli.DepositMoney(ctx, client, name, pin, money)
		case 5:
			log.Infoln("Terminating the client")
			break atmloop
		}

	}
	defer conn.Close()

}
