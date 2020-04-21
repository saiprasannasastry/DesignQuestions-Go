package parkingLot

import (
	"fmt"
	"github.com/goParkingLot/Parking"
	"github.com/goParkingLot/Slot"
)

func CreateLot(n int) *Parking.Parking{
	p :=new(Parking.Parking)
	p.Slots = make([]*Slot.Slot,n)
	p.Capacity =n
	fmt.Printf("created slot with capacity :%v\n", n)
	return p
}

func GetSLotsWIthReg(p *Parking.Parking ,reg string){
	for k,v := range p.Slots{
		if v !=nil && v.Car.RegPlate==reg {
			fmt.Printf("The slot for %v slot is  is :%v\n", reg, k+1)
			return
		}
	}
	fmt.Printf("No slot found with %v ",reg)
}



