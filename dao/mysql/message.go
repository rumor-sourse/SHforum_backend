package mysql

import "SHforum_backend/models"

// SendMessage 发送消息
func SendMessage(user1 int64, user2 int64, title string, content string) (err error) {
	mes := &models.Message{
		Title:       title,
		Content:     content,
		SendUser:    user1,
		ReceiveUser: user2,
		HadRead:     "0",
	}
	result := db.Debug().Create(mes)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
