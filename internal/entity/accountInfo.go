package entity

type AccountInfo struct {
	Received  []Operation
	Sent      []Operation
	Coins     int
	Inventory map[string]int
}

type Operation struct {
	ID       int
	FromUser string
	ToUser   string
	Amount   int
}
