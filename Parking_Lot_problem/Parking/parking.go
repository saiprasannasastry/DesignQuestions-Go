package Parking

import (
	"fmt"
	"github.com/goParkingLot/Parking_Lot_problem/Car"
	"github.com/goParkingLot/Parking_Lot_problem/Slot"
	//log "github.com/sirupsen/logrus"
)

type Parking struct {

	Slots []*Slot.Slot
	Capacity int
}
func (p Parking)NextFreeSlot()(int){
	for k,v :=range p.Slots{
		if v ==nil{
			return k
		}
	}
	fmt.Println("No free slots")
	return -1
}

func (p Parking)Park(car *Car.Car) {

	freeSlot :=p.NextFreeSlot()
	if freeSlot != -1 {
		p.Slots[freeSlot] = &Slot.Slot{Car: car}
		fmt.Printf(" Allocated Slot is :% v\n",freeSlot+1)
	}
	//fmt.Println(p.Slots)

}
func (p Parking)Leave(n int){
	p.Slots[n-1] = nil
	fmt.Printf("Slot %v is free \n",n)
}

func (p Parking) Status(){
	fmt. Println("Slot No.    Reg No.  Color")
	for k,v:= range p.Slots{
		if v!=nil{
			fmt.Printf("%v.          %v    %v\n",k+1, p.Slots[k].Car.RegPlate,p.Slots[k].Car.Color)
		}
	}
}
func (p Parking)RegWithColor(col string){
	s := []string{}
	for _,v := range p.Slots{
		if v !=nil && v.Car.Color==col{
			s = append(s,v.Car.RegPlate)
		}
	}
	fmt.Printf("Reg plate with given color : %v\n",s)
}

func (p Parking) SlotWithColor(col string){
	s := []int{}
	for k,v := range p.Slots{
		if v !=nil && v.Car.Color==col{
			s = append(s,k+1)
		}
	}
	fmt.Printf("Slot No with given color : %v\n",s)
}
//func (p Parking) SlotWithReg(reg string){
//	s :=1
//	count :=0
//	for k,v := range p.Slots{
//		if v !=nil && v.Car.RegPlate==reg {
//			s = k
//			count ++
//		}
//		}
//	if count ==0{
//		fmt.Println("No slots found matching the given Reg Plate")
//		return
//	}
//	fmt.Printf("Slot No with given Reg No : %v\n",s)
//}
