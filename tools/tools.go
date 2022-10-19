package tools

import (
	"crypto/md5"
	"fmt"
)

// GetMd5 把密码进行盐值加密处理
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
