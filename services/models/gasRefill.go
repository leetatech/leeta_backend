package models

type GasRefill struct {
	ID            string              `json:"id" bson:"id"`
	Guest         bool                `json:"guest" bson:"guest"`
	GuestBioData  GuestBioData        `json:"guest_bio_data,omitempty" bson:"guest_bio_data"`
	CustomerID    string              `json:"customer_id" bson:"customer_id"`
	RefillDetails RefillDetails       `json:"refill_details" bson:"refill_details"`
	ShippingInfo  ShippingInfo        `json:"shipping_info,omitempty" bson:"shipping_info"`
	AmountPaid    float64             `json:"amount_paid" bson:"amount_paid"`
	DeliveryFee   float64             `json:"delivery_fee" bson:"delivery_fee"`
	ServiceFee    float64             `json:"service_fee" bson:"service_fee"`
	TotalCost     float64             `json:"total_cost" bson:"total_cost"`
	Status        RefillRequestStatus `json:"status" bson:"status"`
	StatusTs      int64               `json:"status_ts" bson:"status_ts"`
	Ts            int64               `json:"ts" bson:"ts"`
} // @name GasRefill

type GuestBioData struct {
	SessionID string `json:"session_id,omitempty" bson:"session_id"`
	DeviceID  string `json:"device_id" bson:"device_id"`
	FirstName string `json:"first_name,omitempty" bson:"first_name"`
	LastName  string `json:"last_name,omitempty" bson:"last_name"`
	Email     string `json:"email,omitempty" bson:"email"`
	Phone     string `json:"phone,omitempty" bson:"phone"`
} // @name GuestBioData

// ShippingInfo This object is only required from a guest/customer who is ordering for someone else
// If the guest/customer is ordering for himself, then we need to collect their address
type ShippingInfo struct {
	ForMe            bool    `json:"for_me" bson:"for_me"`
	RecipientName    string  `json:"recipient_name,omitempty" bson:"recipient_name"`
	RecipientPhone   string  `json:"recipient_phone,omitempty" bson:"recipient_phone"`
	RecipientEmail   string  `json:"recipient_email,omitempty" bson:"recipient_email"`
	RecipientAddress Address `json:"recipient_address,omitempty" bson:"recipient_address"`
} // @name ShippingInfo

type RefillDetails struct {
	OrderItems []CartItem          `json:"order_items" bson:"order_items"`
	Status     RefillRequestStatus `json:"status" bson:"status"`
	StatusTs   int64               `json:"status_ts" bson:"status_ts"`
	Ts         int64               `json:"ts" bson:"ts"`
} // @name RefillDetails

type RefillRequestStatus string

const (
	RefillCancelled RefillRequestStatus = "cancelled"
	RefillAccepted  RefillRequestStatus = "accepted"
	RefillRejected  RefillRequestStatus = "rejected"
	RefillPending   RefillRequestStatus = "pending"
	RefillFulFilled RefillRequestStatus = "fulfilled"
)
