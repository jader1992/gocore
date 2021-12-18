package contract

const IDKey = "gocore:id"

type IDService interface {
	NewID() string
}