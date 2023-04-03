package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var (
	ErrCanNotOpen  = errors.New("can not open file")
	ErrCanNotRead  = errors.New("can not read file")
	ErrCanNotParse = errors.New("can not parse file")
)
var products []Product

func main() {
	err := obtainData()
	if err != nil {
		panic(err)
	}
}

func obtainData() error {
	file, err := os.Open("../products.json")
	if err != nil {
		return ErrCanNotOpen
	}
	defer file.Close()

	myDecoder := json.NewDecoder(file)

	if err := myDecoder.Decode(&products); err != nil {
		return ErrCanNotRead
	}
	return nil
}
