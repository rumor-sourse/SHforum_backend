package logic

import (
	"SHforum_backend/interval/dao/mysql"
	"SHforum_backend/interval/dao/redis"
	"SHforum_backend/interval/models"
	"SHforum_backend/interval/models/response"
	"SHforum_backend/pkg/jwt"
	"SHforum_backend/pkg/rabbitmq"
	"SHforum_backend/pkg/snowflake"
	"go.uber.org/zap"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		zap.L().Error("用户已存在", zap.Error(err))
		return err
	}
	// 生成UID
	userId := snowflake.GenID("user")
	// 构造一个User实例
	user := &models.User{
		UserID:   userId,
		Username: p.Username,
		Password: p.Password,
		Email:    p.Email,
		Role:     p.Role,
	}
	rcode, err := redis.GetCode(p.Email)
	if rcode != p.Code {
		zap.L().Error("验证码错误", zap.Error(err))
		return err
	}
	// 保存用户信息
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (userResp *response.UserResponse, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的指针，拿到userID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	userResp = &response.UserResponse{
		UserID: user.UserID,
		Name:   user.Username,
		Token:  token,
	}
	return
}

func UpdateUserByID(userID int64, currentUserID int64, p *models.ParamUpdateUser) (err error) {
	//判断是否为管理员或本人
	curruser, err := mysql.GetUserById(currentUserID)
	if err != nil {
		return err
	}
	if curruser.Role != models.Admin && curruser.UserID != userID {
		zap.L().Error("UpdateUserByID failed：不是管理员或本人", zap.Error(err))
	}
	//TODO 修改密码的逻辑
	err = mysql.UpdateUserByID(userID, p.Username)
	if err != nil {
		return err
	}
	return
}

// SendMessage 发送消息
func SendMessage(fromUserID int64, toUserID int64, p *models.ParamSendMessage) (err error) {
	if err = mysql.SendMessage(fromUserID, toUserID, p.Title, p.Content); err != nil {
		return err
	}
	return nil
}

// SendCode 发送邮箱验证码
func SendCode(email string, code string) {
	MQSendCodeMessage(email, code)
}

func MQSendCodeMessage(email string, code string) {
	rmq := rabbitmq.NewRabbitMQSimple("send_code")
	defer rmq.Destroy()
	rmq.PublishSendCodeMessage(email, code)
}

func MQReceiveCodeMessage() {
	rmq := rabbitmq.NewRabbitMQSimple("send_code")
	defer rmq.Destroy()
	rmq.ConsumeCodeMessage()
}

func Follow(userId int64, followeduser int64) (err error) {
	return mysql.Follow(userId, followeduser)
}

func UnFollow(userId int64, followeduser int64) (err error) {
	return mysql.UnFollow(userId, followeduser)
}

func GetFanList(userId int64) (fans []int64, err error) {
	list, err := mysql.GetFanList(userId)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		fans = append(fans, v.FanUser)
	}
	return
}
