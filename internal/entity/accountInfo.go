package entity

type AccountInfo struct {
	Received  []Operation
	Sent      []Operation
	Coins     int
	Inventory map[string]int
}

type Operation struct {
	User   string
	Amount int
}
