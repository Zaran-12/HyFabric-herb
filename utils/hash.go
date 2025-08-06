package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

/* *
* 用于接收处理的内容，并对该内容进行哈希计算
 */
func SHA256(str string) string {
	// 移除空格和换行符
	cleanedStr := strings.TrimSpace(str)
	fmt.Printf("原始字符串: '%s'\n", str)
	fmt.Printf("清理后的字符串: '%s'\n", strings.TrimSpace(str))
	hash := sha256.New()
	hash.Write([]byte(cleanedStr))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)

}
