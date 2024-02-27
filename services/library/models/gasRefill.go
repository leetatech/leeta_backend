package models

type GasRefill struct {
	ID            string              `json:"id" bson:"id"`
	Guest         bool                `json:"guest" bson:"guest"`
	GuestBioData  GuestBioData        `json:"guest_bio_data,omitempty" bson:"guest_bio_data"`
	CustomerID    string              `json:"customer_id" bson:"customer_id"`
	RefillDetails RefillDetails       `json:"refill_details" bson:"refill_details"`
	ShippingInfo  ShippingInfo        `json:"shipping_info,omitempty" bson:"shipping_info"`
	Status        RefillRequestStatus `json:"status" bson:"status"`
	StatusTs      int64               `json:"status_ts" bson:"status_ts"`
	Ts            int64               `json:"ts" bson:"ts"`
} // @name GasRefill

type GuestBioData struct {
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
	ProductID  string          `json:"product_id" bson:"product_id"`
	VendorID   string          `json:"vendor_id,omitempty" bson:"vendor_id"`
	Weight     float32         `json:"weight" bson:"weight"`
	AmountPaid float64         `json:"amount_paid" bson:"amount_paid"`
	GasType    ProductCategory `json:"gas_type" bson:"gas_type"`
} // @name RefillDetails

type RefillRequestStatus string

const (
	RefillCancelled RefillRequestStatus = "cancelled"
	RefillAccepted  RefillRequestStatus = "accepted"
	RefillRejected  RefillRequestStatus = "rejected"
	RefillPending   RefillRequestStatus = "pending"
	RefillFulFilled RefillRequestStatus = "fulfilled"
)
