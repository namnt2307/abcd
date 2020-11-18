package subscription

type SubcriptionObjectStruct struct {
	Start_date                  string
	Expiry_date                 string
	Is_trialed                  int
	Id                          string
	Method                      string
	Suspended_date              string
	Created_at                  string
	Updated_at                  string
	Status                      int
	Current_package_id          int
	Last_success_transaction_id string
	Last_transaction_id         string
	Next_package_id             int
	Paypal_agreement_id         string
	User_id                     string
	Recurring                   int
	Type                        int
}
