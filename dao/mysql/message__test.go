package mysql

import "testing"

func TestSendMessage(t *testing.T) {
	SendMessage(6728893088141312, 6758719832461312, "test", "test")
}
