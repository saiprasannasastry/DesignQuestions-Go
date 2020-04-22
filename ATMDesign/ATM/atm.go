package ATM

import (
	"fmt"
	"github.com/goParkingLot/ATMDesign/types"
)

type Atm struct{
	Details map[string]*types.Persons

}
func (a *Atm)Add(name string,person *types.Persons){
	if a.Details ==nil{
		a.Details= make(map[string]*types.Persons)
		a.Details[name] = person
	}else{
		a.Details[name] = person
	}
	return
}
func (a *Atm)Verify(name string)error{

	if a.Details[name]== nil{
		//equvalent to not able to detect a card
		return fmt.Errorf("sorry ,cant recognize you")
	}
	return nil
}
func (a *Atm)Authenticate(name string ,pin int) error{

	if (a.Details[name].Person.Pin)== pin{
		fmt.Println("Auth successful")
		return nil
	}	else{
		return fmt.Errorf("wrong Pin entered")
	}
}

func (a *Atm) BankBalance(name string) int{
	if a.Details[name].Person.TotalBalance <0{
		return 0
	}
	return a.Details[name].Person.TotalBalance
}

func (a *Atm)WithDrawCash(name string, cash int)(int,error){

	a.Details[name].Person.TotalBalance = a.Details[name].Person.TotalBalance - cash
	if a.Details[name].Person.TotalBalance <0{
		return 0,fmt.Errorf("insufficient balance")
	}
	return a.Details[name].Person.TotalBalance,nil
}