package domain

type OrderRepository interface {
	CreateOrder(request Order)
}
