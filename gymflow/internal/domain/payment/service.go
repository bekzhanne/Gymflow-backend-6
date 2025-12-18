package payment

import "errors"

type Service interface {
	CreatePayment(userID uint, req CreatePaymentRequest) (*Payment, error)
	ListPayments(userID uint) ([]Payment, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePayment(userID uint, req CreatePaymentRequest) (*Payment, error) {
	// 1. Проверяем, существует ли платёж для бронирования
	existing, err := s.repo.FindByBookingID(req.BookingID)
	if err == nil && existing != nil {
		return nil, errors.New("payment already exists for this booking")
	}

	// 2. Создаём платёж со статусом pending
	payment := &Payment{
		UserID:    userID,
		BookingID: req.BookingID,
		Amount:    req.Amount,
		Method:    req.Method,
		Status:    "pending",
	}

	// 3. Сохраняем
	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}


func (s *service) ListPayments(userID uint) ([]Payment, error) {
	return s.repo.ListByUser(userID)
}
