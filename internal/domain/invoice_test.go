package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_domain_NewInvoice(t *testing.T) {
	type args struct {
		id            int
		orderID       int
		total         float64
		tax           float64
		paymentStatus PaymentStatus
	}
	tests := []struct {
		name string
		args args
		want Invoice
	}{
		{
			name: "Valid Invoice",
			args: args{
				id:            1,
				orderID:       100,
				total:         250.75,
				tax:           20.25,
				paymentStatus: Paid,
			},
			want: Invoice{
				ID:            1,
				OrderID:       100,
				Total:         250.75,
				Tax:           20.25,
				PaymentStatus: Paid,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInvoice(tt.args.id, tt.args.orderID, tt.args.total, tt.args.tax, tt.args.paymentStatus)
			assert.Equal(t, tt.want, got, "NewInvoice() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_Invoice_Validate(t *testing.T) {
	tests := []struct {
		name string
		inv  Invoice
		want bool
	}{
		{
			name: "Valid Invoice",
			inv: Invoice{
				ID:            1,
				OrderID:       100,
				Total:         250.75,
				Tax:           20.25,
				PaymentStatus: Paid,
			},
			want: true,
		},
		{
			name: "Invalid OrderID",
			inv: Invoice{
				ID:            2,
				OrderID:       0,
				Total:         150.00,
				Tax:           15.00,
				PaymentStatus: Unpaid,
			},
			want: false,
		},
		{
			name: "Negative Total",
			inv: Invoice{
				ID:            3,
				OrderID:       101,
				Total:         -50.00,
				Tax:           5.00,
				PaymentStatus: Cancelled,
			},
			want: false,
		},
		{
			name: "Negative Tax",
			inv: Invoice{
				ID:            4,
				OrderID:       102,
				Total:         200.00,
				Tax:           -10.00,
				PaymentStatus: Processing,
			},
			want: false,
		},
		{
			name: "Invalid Payment Status",
			inv: Invoice{
				ID:            5,
				OrderID:       103,
				Total:         300.00,
				Tax:           30.00,
				PaymentStatus: "unknown",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.inv.Validate()
			assert.Equal(t, tt.want, got, "Invoice.Validate() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_Invoice_BillWithTax(t *testing.T) {
	type fields struct {
		ID            int
		OrderID       int
		Total         float64
		Tax           float64
		PaymentStatus PaymentStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Normal Case",
			fields: fields{
				ID:            1,
				OrderID:       100,
				Total:         250.75,
				Tax:           20.25,
				PaymentStatus: Paid,
			},
			want: 271.00,
		},
		{
			name: "Zero Tax",
			fields: fields{
				ID:            2,
				OrderID:       101,
				Total:         150.00,
				Tax:           0.00,
				PaymentStatus: Unpaid,
			},
			want: 150.00,
		},
		{
			name: "Zero Total and Tax",
			fields: fields{
				ID:            3,
				OrderID:       102,
				Total:         0.00,
				Tax:           0.00,
				PaymentStatus: Cancelled,
			},
			want: 0.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := &Invoice{
				ID:            tt.fields.ID,
				OrderID:       tt.fields.OrderID,
				Total:         tt.fields.Total,
				Tax:           tt.fields.Tax,
				PaymentStatus: tt.fields.PaymentStatus,
			}
			got := inv.BillWithTax()
			assert.Equal(t, tt.want, got, "Invoice.BillWithTax() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_PaymentStatus_Validate(t *testing.T) {
	tests := []struct {
		name string
		ps   PaymentStatus
		want bool
	}{
		{
			name: "Valid Status - Paid",
			ps:   Paid,
			want: true,
		},
		{
			name: "Valid Status - Unpaid",
			ps:   Unpaid,
			want: true,
		},
		{
			name: "Valid Status - Cancelled",
			ps:   Cancelled,
			want: true,
		},
		{
			name: "Valid Status - Processing",
			ps:   Processing,
			want: true,
		},
		{
			name: "Valid Status - Failed",
			ps:   Failed,
			want: true,
		},
		{
			name: "Invalid Status",
			ps:   "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ps.Validate()
			assert.Equal(t, tt.want, got, "PaymentStatus.Validate() = %v, want %v", got, tt.want)
		})
	}
}
