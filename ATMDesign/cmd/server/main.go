package main

import (
	"database/sql"
	"github.com/DesignQuestions-Go/ATMDesign/pkg/pb"
	"github.com/DesignQuestions-Go/ATMDesign/pkg/svc"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

func main() {

	//doneC := make(chan error)
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
		//	doneC <- err
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	var db *sql.DB

	db, err = sql.Open("mysql", "root:infoblox@tcp(127.0.0.1:3306)/atmdesign")
	if err != nil {
		log.Fatalf("Unable to open db :%v", err)
	}

	srv := svc.NewBankServer(db)
	pb.RegisterBankServiceServer(s, srv)
	log.Infof("Starting the server")
	go func() {
		if err := s.Serve(listener); err != nil {
			db.Close()
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//	s.Serve(listener)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Infof(" Received signal Interrupt stopping the server")
	db.Close()
	s.Stop()
	listener.Close()
}
