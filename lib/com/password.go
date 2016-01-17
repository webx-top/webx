package com

// 生成哈希值
func Hash(str string) string {
	return Sha256(str)
}

// 盐值加密
func Salt() string {
	return Hash(RandStr(64))
}

// 创建密码
func MakePassword(password string, salt string) string {
	return Hash(salt + password)
}

// 检查密码(密码原文，数据库中保存的哈希过后的密码，数据库中保存的盐值)
func CheckPassword(rawPassword string, hashedPassword string, salt string) bool {
	return MakePassword(rawPassword, salt) == hashedPassword
}
