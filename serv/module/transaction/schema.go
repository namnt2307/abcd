package transaction

type TransactionObjectStruct struct {
	Id                  string
	User_id             string
	Start_date          string
	Expiry_date         string
	Type                string
	Is_trialed          int
	Platform            string
	Currency            string
	Amount              int
	Setup_fee           string
	Payment_method      string
	Readable_id         string
	Status              int
	Active_code_id      string
	Package_id          string
	Paypal_agreement_id string
	Paypal_ipn_id       string
	Paypal_payer_id     string
	Created_at          string
	Updated_at          string
}