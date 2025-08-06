package service

func LoginCheck(username, password string) (bool, map[string]interface{}) {

	mapp := make(map[string]interface{})
	if len(username) < 6 || len(username) > 10 {
		mapp["code"] = "20002"
		mapp["msg"] = "用户名不符合规范，请检查！"
		return false, mapp
	}
	if len(password) < 6 || len(password) > 10 {
		mapp["code"] = "20003"
		mapp["msg"] = "密码不符合规范，请检查！"
		return false, mapp
	}
	return true, nil
}
func RegisterCheck(username, password, RepeatPwd string) (bool, map[string]interface{}) {
	mapp := make(map[string]interface{})
	if len(username) < 6 || len(username) > 10 {
		mapp["code"] = "10002"
		mapp["msg"] = "用户名不符合规范，请检查！"
		return false, mapp
	}
	if len(password) < 6 || len(password) > 10 {
		mapp["code"] = "10003"
		mapp["msg"] = "密码不符合规范，请检查！"
		return false, mapp
	}
	if password != RepeatPwd {
		mapp["code"] = "10004"
		mapp["msg"] = "密码和确认密码不一致，请检查!"
		return false, mapp
	}
	return true, nil
}
