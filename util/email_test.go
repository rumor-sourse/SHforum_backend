package util

import "testing"

func TestSendEmailWithText(t *testing.T) {
	SendEmailWithCode([]string{"chunkai_wang@qq.com"}, "123456")
}
