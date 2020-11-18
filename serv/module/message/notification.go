package message

import (
	. "cm-v5/serv/module"
)

func SaveDataPushToken(push_token, access_token, user_id, platform string) error {
	return nil
	
	// Connect MySQL
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return err
	}
	defer DbMysql.Close()

	//Update data Mysql
	// sqlRaw := fmt.Sprintf(`
	// 	UPDATE cas_access_token
	// 	SET push_token = '%s'
	// 	WHERE user_id = '%s' AND platform = '%s' AND access_token = '%s'`, push_token, user_id, platform, access_token)
	_, err = DbMysql.Query(`
	UPDATE cas_access_token 
	SET push_token = ?
	WHERE user_id = ? AND platform = ? AND access_token = ?`, push_token, user_id, platform, access_token)
	if err != nil {
		return err
	}

	return err
}
