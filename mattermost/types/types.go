package types

type Team struct{
	ID string `json:"id"`
	Name string `json:"name"`
	DisplayName string `json:"display_name"`
	Type string `json:"type"`
}

type Chanels []Chanel

type Chanel struct{
	ID string `json:"id"`
	Name string `json:"name"`
}