package rabbitmq

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/es"
	"SHforum_backend/models"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type CreatePostMessage struct {
	UserID      int64   `json:"user_id"`
	Fans        []int64 `json:"fans"`
	Message     string  `json:"message"`
	models.Post `json:"post"`
}

// PublishCreatePostMessage 传递关注的用户创建了贴子的消息
func (r *RabbitMQ) PublishCreatePostMessage(userID int64, fans []int64, message string, post models.Post) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mes := CreatePostMessage{
		UserID:  userID,
		Fans:    fans,
		Message: message,
		Post:    post,
	}
	jsonbody, err := json.Marshal(mes)
	//2.发送消息
	err = r.channel.PublishWithContext(ctx,
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonbody,
		})
}

// ConsumeCreatePostMessageByMysql Mysql消费端收到已创建贴子的消息
func (r *RabbitMQ) ConsumeCreatePostMessageByMysql() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"fanout",
		true,
		false,
		//YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil)
	r.failOnErr(err, "Failed to bind a queue")
	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			var mes CreatePostMessage
			err = json.Unmarshal(d.Body, &mes)
			if err != nil {
				return
			}
			for _, fan := range mes.Fans {
				err := mysql.SendMessage(mes.UserID, fan, "新贴子提醒", mes.Message)
				if err != nil {
					return
				}
			}
		}
	}()
	<-forever
}

// ConsumeCreatePostMessageByEs Es消费端收到已创建贴子的消息
func (r *RabbitMQ) ConsumeCreatePostMessageByEs() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"fanout",
		true,
		false,
		//YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil)
	r.failOnErr(err, "Failed to bind a queue")
	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			var mes CreatePostMessage
			err = json.Unmarshal(d.Body, &mes)
			if err != nil {
				return
			}
			//将贴子信息存入es
			err := es.CreatePostIndex(mes.Post)
			if err != nil {
				zap.L().Error("es.CreatePostIndex(mes.Post) failed", zap.Error(err))
				return
			}
		}
	}()
	<-forever
}
