package payment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock Repository
type MockPaymentRepository struct {
	mock.Mock
}



func (m *MockPaymentRepository) Create(payment *Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByID(id uint) (*Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByBookingID(bookingID uint) (*Payment, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Payment), args.Error(1)
}

func (m *MockPaymentRepository) ListByUser(userID uint) ([]Payment, error) {
	args := m.Called(userID)
	return args.Get(0).([]Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(payment *Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

// Tests
func TestCreatePayment_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewService(mockRepo)

	req := CreatePaymentRequest{
		BookingID: 1,
		Amount:    50.0,
		Method:    "card", // Изменено: было PaymentMethod, теперь Method
	}

	mockRepo.On("FindByBookingID", uint(1)).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*payment.Payment")).Return(nil)
	

	

	payment, err := service.CreatePayment(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, req.Amount, payment.Amount)
	assert.Equal(t, "pending", payment.Status)
	mockRepo.AssertExpectations(t)
}

func TestCreatePayment_AlreadyExists(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewService(mockRepo)

	existingPayment := &Payment{
		ID:        1,
		BookingID: 1,
		Status:    "completed",
	}

	req := CreatePaymentRequest{
		BookingID: 1,
		Amount:    50.0,
		Method:    "card",
	}

	mockRepo.On("FindByBookingID", uint(1)).Return(existingPayment, nil)

	payment, err := service.CreatePayment(1, req)

	assert.Error(t, err)
	assert.Nil(t, payment)
	mockRepo.AssertExpectations(t)
}

func TestListPayments_Success(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewService(mockRepo)

	expectedPayments := []Payment{
		{ID: 1, UserID: 1, Amount: 50.0, Status: "completed"},
		{ID: 2, UserID: 1, Amount: 30.0, Status: "completed"},
	}

	mockRepo.On("ListByUser", uint(1)).Return(expectedPayments, nil)

	payments, err := service.ListPayments(1)

	assert.NoError(t, err)
	assert.Len(t, payments, 2)
	assert.Equal(t, 50.0, payments[0].Amount)
	mockRepo.AssertExpectations(t)
}

func TestListPayments_Empty(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	service := NewService(mockRepo)

	emptyPayments := []Payment{}

	mockRepo.On("ListByUser", uint(1)).Return(emptyPayments, nil)

	payments, err := service.ListPayments(1)

	assert.NoError(t, err)
	assert.Len(t, payments, 0)
	mockRepo.AssertExpectations(t)
}