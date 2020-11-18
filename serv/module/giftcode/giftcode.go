package giftcode

import (
	"errors"
	"fmt"
	"strings"
	"time"

	. "cm-v5/serv/module"
	Package "cm-v5/serv/module/packages"
	Subscription "cm-v5/serv/module/subscription"
	SubscriptionLog "cm-v5/serv/module/subscription_log"
	Transaction "cm-v5/serv/module/transaction"
	"gopkg.in/mgo.v2/bson"
)

func CheckGiftCode(userId, code string) (interface{}, error) {
	var Message struct {
		Valid int `json:"valid" `
		Used  int `json:"used" `
	}

	session, db, err := GetCollection()
	if err != nil {
		return Message, errors.New("an_error_occurred")
	}

	defer session.Close()

	//Kiem tra gift code hop le
	var where = bson.M{
		"code":      code,
		"expire_at": bson.M{"$gte": time.Now().Unix()},
		"status":    1,
		"type_code": TYPE_GIFT_CODE,
	}

	var GiftCode struct {
		Package_id  int `json:"package_id" `
		Hit_use     int `json:"hit_use" `
		Number_uses int `json:"number_uses" `
	}

	err = db.C(COLLECTION_GIFT_CODE).Find(where).One(&GiftCode)
	if err != nil {
		return Message, errors.New("not_valid")
	}

	if GiftCode.Hit_use >= GiftCode.Number_uses {
		return Message, errors.New("not_valid")
	}

	Message.Valid = 1
	if userId == "" {
		return Message, nil
	}

	//kiem tra user da dung gift code hay chua
	var Tran Transaction.TransactionObjectStruct
	Tran.User_id = userId
	Tran.Active_code_id = code
	err = Tran.GetByUserAndPackage()
	if err != nil {
		return Message, errors.New("not_valid")
	}

	//Da dung gift code tra ve kq
	if Tran.Id != "" {
		Message.Used = 1
		return Message, nil
	}

	return Message, nil
}


func CheckPromotionCode(userId, code string, package_id int) (interface{}, error) {
	var Message struct {
		Valid int `json:"valid" `
		Used  int `json:"used" `
		Price int `json:"price" `
		Duration_hour int `json:"duration_hour" ` // Theo giờ
		Type_code string `json:"type_code" `
	}

	session, db, err := GetCollection()
	if err != nil {
		return Message, errors.New("an_error_occurred")
	}

	defer session.Close()

	//Kiem tra gift code hop le
	var where = bson.M{
		"code":      code,
		"expire_at": bson.M{"$gte": time.Now().Unix()},
		"status":    1,
		"type_code": TYPE_VOUCHER_CODE,
		"package_id": package_id,
		
	}

	var GiftCode struct {
		Package_id  int `json:"package_id" `
		Hit_use     int `json:"hit_use" `
		Number_uses int `json:"number_uses" `
		Price int `json:"price" `
		Duration int `json:"duration" ` // Theo giờ
		Type_promotion_code int `json:"type_promotion_code" ` // Neu la promotion code => code giam tien = 1 || code tang gio su dung = 2 (gift code k can quan tam)
	}

	err = db.C(COLLECTION_GIFT_CODE).Find(where).One(&GiftCode)
	if err != nil {
		return Message, errors.New("not_valid")
	}

	if GiftCode.Hit_use >= GiftCode.Number_uses {
		return Message, errors.New("not_valid")
	}

	Message.Valid = 1

	if GiftCode.Type_promotion_code == 1 {
		Message.Type_code = "price"
	} else if GiftCode.Type_promotion_code == 2  {
		Message.Type_code = "duration"	
	}

	Message.Price = GiftCode.Price
	Message.Duration_hour = GiftCode.Duration

	if userId == "" {
		return Message, nil
	}


	//kiem tra user da dung gift code hay chua
	var Tran Transaction.TransactionObjectStruct
	Tran.User_id = userId
	Tran.Active_code_id = code
	err = Tran.GetByUserAndPackage()
	if err != nil {
		return Message, errors.New("not_valid")
	}

	//Da dung gift code tra ve kq
	if Tran.Id != "" {
		Message.Used = 1
		return Message, nil
	}

	return Message, nil
}


func Giftcode_UseCode(userId, code string) (DataStruct, error) {
	var data DataStruct
	keyCache := "USE_CODE_NOT_VALID_" + userId
	numFail, _ := mRedisUSC.GetInt(keyCache)
	if numFail >= 5 {
		return data, errors.New("too_many_fail")
	}

	var WriteCacheBlock = func(keyCache string, numFail int) {
		mRedisUSC.SetInt(keyCache, numFail+1, 60*3)
	}

	//Kiem tra code hien tai dang nhap la gift code hay trail code
	idPackage, typeCode, err := CheckTypeCode(code)
	if err != nil || idPackage <= 0 {

		WriteCacheBlock(keyCache, numFail)

		data.Message = "fail"
		return data, errors.New("code_not_valid")
	}

	// neu la trail code thi can xac thuoc prefix + user_id co giong vs code truyen vao k
	codeTrial := PREFIX_CODE_TRIAL + "-" + userId[0:5]
	if typeCode == "trial_code" && (strings.ToLower(code) != strings.ToLower(codeTrial)) {

		WriteCacheBlock(keyCache, numFail)

		data.Message = "fail"
		return data, errors.New("code_not_valid")
	}

	if typeCode == "gift_code" {
		//Cap nhat lai luot su dung truoc khi xac nhan giao dich va kich hoat
		err := UpdateHitUseCode(code, 1)
		if err != nil {
			data.Message = "fail"
			data.Error = err.Error()
			return data, errors.New("an_error_occurred")
		}
	}

	//luu thong tin giao dich
	data, err = SaveInfoTransaction(userId, code, idPackage)
	if err != nil {
		//Cap nhat lai luot su dung da tang neu co loi xay ra
		UpdateHitUseCode(code, -1)

		data.Message = "fail"
		data.Error = err.Error()
		return data, errors.New("an_error_occurred")
	}

	return data, nil
}

func UpdateHitUseCode(code string, val int) error {
	session, db, err := GetCollection()
	if err != nil {
		return err
	}

	defer session.Close()

	var where = bson.M{
		"code": code,
	}
	var update = bson.M{"$inc": bson.M{"hit_use": val}}

	err = db.C(COLLECTION_GIFT_CODE).Update(where, update)
	if err != nil {
		return err
	}

	return nil
}

func CheckTypeCode(code string) (int, string, error) {
	i := strings.Index(code, PREFIX_CODE_TRIAL)
	if i >= 0 {
		return PACKAGE_ID_TRIAL_LAUCHING, "trial_code", nil
	}

	session, db, err := GetCollection()
	if err != nil {
		return 0, "", err
	}

	defer session.Close()

	var where = bson.M{
		"code":      code,
		"expire_at": bson.M{"$gte": time.Now().Unix()},
		"status":    1,
		"type_code": TYPE_GIFT_CODE,
	}

	var GiftCode struct {
		Package_id  int `json:"package_id" `
		Hit_use     int `json:"hit_use" `
		Number_uses int `json:"number_uses" `
	}

	err = db.C(COLLECTION_GIFT_CODE).Find(where).One(&GiftCode)
	if err != nil {
		return 0, "", err
	}

	if GiftCode.Hit_use >= GiftCode.Number_uses {
		return 0, "", errors.New("invalid_use")
	}

	return GiftCode.Package_id, "gift_code", nil
}

func SaveInfoTransaction(userId, code string, idPackage int) (DataStruct, error) {
	var data DataStruct

	// Get Package Detail
	var Pack Package.PackagesObjectStruct
	err := Pack.FetchDetailById(idPackage)
	if err != nil {
		data.Message = "fail"
		data.Error = err.Error()
		return data, errors.New("code_not_valid")
	}

	// Check User valid with core
	var Tran Transaction.TransactionObjectStruct
	Tran.User_id = userId
	Tran.Package_id = fmt.Sprint(Pack.Id)
	Tran.Active_code_id = code
	err = Tran.GetByUserAndPackage()
	if err != nil {
		data.Message = "fail"
		data.Error = err.Error()
		return data, errors.New("package_not_valid")
	}
	if Tran.Id != "" {
		data.Message = "fail"
		return data, errors.New("code_used")
	}

	t := time.Now()
	// Save Transaction
	Tran.Start_date = t.Format("2006-01-02 15:04:05")
	// Tran.Expiry_date = t.AddDate(0, 0, +(Pack.Duration / 24)).Format("2006-01-02 15:04:05")
	Tran.Expiry_date = GetExpiredDateFromDuration(Pack.Duration, Pack.Duration_type, "2006-01-02 15:04:05")
	Tran.Is_trialed = 1
	Tran.Amount = Pack.Price
	Tran.Status = 1

	err = Tran.Save()
	if err != nil {
		data.Message = "fail"
		data.Error = err.Error()
		return data, errors.New("package_not_valid")
	}

	// Save Subscription Log
	var SubLog SubscriptionLog.SubcriptionLogObjectStruct
	SubLog.User_id = userId
	SubLog.Package_id = fmt.Sprint(Pack.Id)
	SubLog.Subscribed_date = Tran.Start_date
	SubLog.Expired_date = Tran.Expiry_date
	SubLog.Is_active = true

	err = SubLog.Save()
	if err != nil {
		data.Message = "fail"
		data.Error = err.Error()
		return data, errors.New("package_not_valid")
	}

	// Save Subscription
	var Sub Subscription.SubcriptionObjectStruct
	Sub.Start_date = Tran.Start_date
	Sub.Expiry_date = Tran.Expiry_date
	Sub.User_id = userId
	Sub.Current_package_id = Pack.Id
	Sub.Next_package_id = Pack.Id
	Sub.Last_success_transaction_id = Tran.Id
	Sub.Last_transaction_id = Tran.Id
	Sub.Status = 1
	err = Sub.Save()
	if err != nil {
		data.Message = "fail"
		data.Error = err.Error()
		return data, errors.New("package_not_valid")
	}

	data.Message = "success"
	return data, nil
}
