package main

import (
	"fmt"
	"github.com/goParkingLot/ATMDesign/ATM"
	"github.com/goParkingLot/ATMDesign/types"
)

func main() {

	person1 := &types.Persons{Person: &types.User{10000, 1234}}
	person2 := &types.Persons{Person: &types.User{100000, 1235}}
	person3 := &types.Persons{Person: &types.User{1000000, 1235}}
	person4 := &types.Persons{Person: &types.User{1000000, 1224}}

	personsList := new(ATM.Atm)
	personsList.Add("test1", person1)
	personsList.Add("test2", person2)
	personsList.Add(" test3", person3)
	personsList.Add("test4", person4)
	//Step 1 Authenticate
	//Step2 options, withdraw cash or check balance
	//step3 ask for balance statment
	fmt.Println("What is your name?")
	var name string
	var pin int
	var option int
	var cash int

	fmt.Scan(&name)
	verify := personsList.Verify(name)
	if verify == nil {
		fmt.Println("What is your pin")
		fmt.Scan(&pin)
		if err := personsList.Authenticate(name, pin); err != nil {
			fmt.Println(err)
			fmt.Println("What is your pin")
			fmt.Scan(&pin)
			if err := personsList.Authenticate(name, pin); err != nil {
				fmt.Println(err)
			}
			fmt.Println("What is your pin")
			fmt.Scan(&pin)
			if err := personsList.Authenticate(name, pin); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("\n The options available are \n 1.%v\n 2.%v\n\n", "BankBalance", "WithDraw Cash")
			fmt.Println("please enter the option number")
			fmt.Scan(&option)
			switch option {
			case 1:
				fmt.Println()
				fmt.Println(personsList.BankBalance(name))

			case 2:
				fmt.Println("Enter the amount to withdraw")
				fmt.Scan(&cash)
				_, err := personsList.WithDrawCash(name, cash)
				if err == nil {
					var option string
					fmt.Printf("\nDo you want to know you balance \n")
					fmt.Println("Type Yes or NO ")
					fmt.Scan(&option)
					if option == "yes" {
						fmt.Printf("your balance is %v\n", personsList.BankBalance(name))
					} else {
						fmt.Printf("\nThankyou\n")
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	}else {
		fmt.Println(verify)
	}
	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt)
	//<-c
	//fmt.Println("Terminating the op...")


}
