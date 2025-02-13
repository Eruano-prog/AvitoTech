package entity

type AccountInfo struct {
	Received  []Operation
	Sent      []Operation
	Coins     int
	Inventory map[string]int
}

type Operation struct {
	Id       int
	FromUser string
	ToUser   string
	Amount   int
}
