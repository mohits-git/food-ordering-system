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
	Total         float64
	Tax           float64
	PaymentStatus PaymentStatus
}

func NewInvoice(id int, orderID int, total, tax float64, paymentStatus PaymentStatus) Invoice {
	return Invoice{
		ID:            id,
		OrderID:       orderID,
		Total:         total,
		Tax:           tax,
		PaymentStatus: paymentStatus,
	}
}

func (i *Invoice) Validate() bool {
	if i.OrderID <= 0 || i.Total < 0 || i.Tax < 0 {
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
	return i.Total + i.Tax
}
