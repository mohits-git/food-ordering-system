package domain

type PaymentStatus string
const (
	Paid       PaymentStatus = "paid"
	Unpaid     PaymentStatus = "unpaid"
	Cancelled  PaymentStatus = "cancelled"
	Processing PaymentStatus = "processing"
	Failed     PaymentStatus = "failed"
)

func (p PaymentStatus) Validate() bool {
	switch p {
	case Paid, Unpaid, Cancelled, Processing, Failed:
		return true
	}
	return false
}

type Invoice struct {
	ID            int
	OrderID       int
	TotalAmount   float64
	TaxAmount     float64
	PaymentStatus PaymentStatus
}

func NewInvoice(id int, orderID int, totalAmount, taxAmount float64, paymentStatus PaymentStatus) *Invoice {
	return &Invoice{
		ID:            id,
		OrderID:       orderID,
		TotalAmount:   totalAmount,
		TaxAmount:     taxAmount,
		PaymentStatus: paymentStatus,
	}
}

func (i *Invoice) Validate() bool {
	if i.OrderID <= 0 || i.TotalAmount < 0 || i.TaxAmount < 0 {
		return false
	}
	return i.PaymentStatus.Validate()
}

func (i *Invoice) UpdatePaymentStatus(status PaymentStatus) bool {
	if !status.Validate() {
		return false
	}
	i.PaymentStatus = status
	return true
}

func (i *Invoice) BillWithTax() float64 {
	return i.TotalAmount + i.TaxAmount
}
