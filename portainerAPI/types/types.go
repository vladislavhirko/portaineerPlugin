package types

type Containers []Containeer

type Containeer struct {
	Id     string
	Names  [1]string
	Image  string
	State  string
	Status string
}
