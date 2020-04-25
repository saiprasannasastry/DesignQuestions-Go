package svc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DesignQuestions-Go/ATMDesign/pkg/pb"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	//	"encoding/json"
)

type BankServiceServer struct {
	db *sql.DB
}
type atm struct {
	name  string
	pin   int32
	money int32
}

func NewBankServer(database *sql.DB) pb.BankServiceServer {
	return &BankServiceServer{db: database}
}

func (r *BankServiceServer) RegisterTBank(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	tx, _ := r.db.Begin()
	var resp string
	rows := tx.QueryRow("select  name from atm where name = ?", req.GetName())

	err := rows.Scan(&req.Name)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Errorln(err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		log.Infof("No Row present Proceeding to insert")
		_, err := tx.Exec("INSERT IGNORE INTO atm(name,balance,pin) VALUES(?,?,?)", req.GetName(), req.GetMoney(), req.GetPin())
		if err != nil {
			tx.Rollback()
			log.Errorln(err)
			return nil, err
		}
		resp = fmt.Sprintf("The User %s has been recorded to the bank", req.GetName())
	} else {
		resp = fmt.Sprintf("The User %s is already present in the bank", req.GetName())
	}

	tx.Commit()
	return &pb.RegisterResponse{Id: resp}, nil
}

func (r *BankServiceServer) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	tx, _ := r.db.Begin()
	for k, v := range req.Req {
		rows := tx.QueryRow("select name,pin from atm where name=? and pin=?", &k, &v)
		err := rows.Scan(&k, &v)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			log.Errorln(err)
			return nil, err
		}

		if err == sql.ErrNoRows {
			tx.Rollback()
			log.Errorf("No Row present to authenticate")
			return &pb.AuthenticateResponse{Authenticated: false}, err

		}
		tx.Commit()
		return &pb.AuthenticateResponse{Authenticated: true}, nil
	}
	return nil, nil
}
func (r *BankServiceServer) Deposit(ctx context.Context, req *pb.DepositMoney) (*pb.BankBalanceResponse, error) {
	tx, _ := r.db.Begin()
	s := atm{}
	row := tx.QueryRow("select * from atm where name =?", req.GetName())
	row.Scan(&s.name, &s.money, &s.pin)
	_, err := tx.Exec("Update  atm set balance=?", s.money+req.GetMoney())
	if err != nil {
		tx.Rollback()
		log.Errorln(err)
		return nil, err
	}
	tx.Commit()
	return &pb.BankBalanceResponse{Money: s.money + req.GetMoney()}, nil

}

func (r *BankServiceServer) BankBalance(ctx context.Context, req *pb.BankBalanceRequest) (*pb.BankBalanceResponse, error) {
	tx, _ := r.db.Begin()
	s := atm{}
	rows := tx.QueryRow("select balance from atm where name =?", req.GetName())
	rows.Scan(&s.money)
	tx.Commit()
	return &pb.BankBalanceResponse{Money: s.money}, nil
}

func (r *BankServiceServer) WithDraw(ctx context.Context, req *pb.DepositMoney) (*pb.BankBalanceResponse, error) {
	tx, _ := r.db.Begin()

	s := atm{}
	rows := tx.QueryRow("select *from atm where name =?", req.GetName())
	rows.Scan(&s.name, &s.money, &s.pin)
	if s.money-req.GetMoney() < 0 {
		tx.Rollback()
		return nil, fmt.Errorf("You do not have sufficient balance")

	}
	_, err := tx.Exec("Update  atm set balance=?", s.money-req.GetMoney())
	if err != nil {
		tx.Rollback()
		log.Errorln(err)
		return nil, err
	}
	tx.Commit()
	return &pb.BankBalanceResponse{Money: s.money - req.GetMoney()}, nil
}
