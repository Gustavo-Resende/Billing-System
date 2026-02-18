package repository

// EventStore define o contrato para armazernar eventos de domínio.
// Por enquanto interface genérica, será refinada na implementação.
type EventStore interface {
	Append(event interface{}) error
}
