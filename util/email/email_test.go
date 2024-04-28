package email

import (
	"testing"
)

func TestSendEmailWithText(t *testing.T) {
	err := SendEmailWithCode([]string{"chunkai_wang@qq.com"}, "123456")
	if err != nil {
		return
	}
}
