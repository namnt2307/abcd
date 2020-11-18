package giftcode

const (
	COLLECTION_GIFT_CODE = "gift_code"

	PACKAGE_ID_TRIAL_LAUCHING = 1
	PREFIX_CODE_TRIAL         = "VIEONVIP"
	TYPE_GIFT_CODE            = 1
	TYPE_VOUCHER_CODE         = 2
)

type DataStruct struct {
	Message string `json:"message" `
	Error   string `json:"error" `
}
