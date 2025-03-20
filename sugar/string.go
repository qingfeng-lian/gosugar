package sugar

// =============================================================================
// 字符串相关处理
// =============================================================================

// StringSub  按照字符截取字符串，从0开始截取指定长度的字符串
func StringSub(text string, subLen int) string {
	// 将字符串转换为[]rune以便正确处理多字节字符
	runes := []rune(text)
	// 确保不会在字符中间截断
	truncatedRunes := runes
	if len(runes) > subLen {
		truncatedRunes = runes[:subLen]
	}

	// 将[]rune转回string
	return string(truncatedRunes)
}

// StringInSlice 检查字符串是否在切片中
func StringInSlice(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
