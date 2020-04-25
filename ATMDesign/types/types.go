package types

type Persons struct {
	Person *User
}

//Using a Name instead of Debit Card number
type User struct {
	TotalBalance int
	Pin          int
}
