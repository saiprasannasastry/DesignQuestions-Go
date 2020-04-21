package main

import (
	//"fmt"
	"github.com/goParkingLot/Car"
	"github.com/goParkingLot/parkingLot"
	//"github.com/goParkingLot/types"
)

func main()  {
	created := parkingLot.CreateLot(6)
	car :=&Car.Car{RegPlate: "ka-41",Color: "Red"}
	created.Park(car)
	created.Park(&Car.Car{"Ka-42","Green"})
	created.Park(&Car.Car{"Ka-43","Yellow"})
	created.Park(&Car.Car{"Ka-44","Red"})
	created.Park(&Car.Car{"Ka-45","Green"})
	created.Park(&Car.Car{"Ka-46","Green"})
	//fmt.Println(Parking.Parking{})
	created.Park(&Car.Car{"Ka-47","orange"})
	created.Leave(4)
	//created.Status()
	created.Park(&Car.Car{"Ka-48","Green"})
	created.Status()
	created.RegWithColor("Green")
	created.SlotWithColor("Green")
	//created.SlotWithReg("KA-41")
	//created.SlotWithReg("KA-50")
	parkingLot.GetSLotsWIthReg(created ,"Ka-46")
	parkingLot.GetSLotsWIthReg(created ,"Ka-56")

}


