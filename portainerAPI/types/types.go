package types

type Containers []Containeer

//Структура контейнера с докера
type Containeer struct {
	Id     string
	Names  [1]string
	Image  string
	State  string
	Status string
	Logs string
}
