package model

import (
	"encoding/json"
	"fmt"
)

type OrderStatus int

const (
	Unspecified OrderStatus = iota
	PENDING
	PAID
	CANCELED
)

func (os OrderStatus) String() string {
	strings := [...]string{"Unespecifed", "PENDING", "PAID", "CANCELLED"}
	if os < 0 || int(os) >= len(strings) {

		return "Unspecified"
	}
	return strings[os]
}

func (os *OrderStatus) Scan(src any) error {
	if src == nil {
		*os = Unspecified
		return nil
	}

	var statusStr string
	switch val := src.(type) {
	case string:
		statusStr = val
	case []byte:
		statusStr = string(val)
	default:
		return fmt.Errorf("tipo invalido para OrderStatus: %T", src)
	}

	// Mapeia o texto vindo do banco para o seu iota correspondente
	switch statusStr {
	case "PENDING":
		*os = PENDING
	case "PAID":
		*os = PAID
	case "CANCELED", "CANCELLED": // Aceita ambas as grafias por segurança
		*os = CANCELED
	default:
		*os = Unspecified
	}

	return nil
}

// pra sair como texto no json
func (os OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}
