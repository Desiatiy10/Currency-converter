package model

type Entity interface {
	GetType() string
}

func (c Currency) GetType() string {
	return "Currency"
}

func (c Conversion) GetType() string {
	return "Conversion"
}
