package json

import (
	"encoding/json"
)

type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

// Какая тут проблема?
func ValidateV1(data []byte) (*Transaction, error) {
	temp := make(map[string]interface{})
	json.Unmarshal(data, &temp)

	var tx Transaction
	json.Unmarshal(data, &tx)
	return &tx, nil
}

// ok: один парсинг, прямо в структуру
func ValidateV2(data []byte) (*Transaction, error) {
	var tx Transaction
	json.Unmarshal(data, &tx)
	return &tx, nil
}

// ok*2: пул объектов
var pool = make(chan *Transaction, 100)

func init() {
	for i := 0; i < 100; i++ {
		pool <- &Transaction{}
	}
}

func ValidateV3(data []byte) (*Transaction, error) {
	tx := <-pool
	json.Unmarshal(data, tx)
	defer func() { pool <- tx }()
	return tx, nil
}
