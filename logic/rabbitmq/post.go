package rabbitmq

import (
	"SHforum_backend/dao/mysql"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

// PublishCreatePostMessage 传递关注的用户创建了贴子的消息
func (r *RabbitMQ) PublishCreatePostMessage(message string) {
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
	//2.发送消息
	err = r.channel.PublishWithContext(ctx,
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "Failed to publish a message")
}

// ConsumeCreatePostMessage 消费端收到已创建贴子的消息
func (r *RabbitMQ) ConsumeCreatePostMessage(userID int64, fan int64) {
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
		false,
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
	messages, err := r.channel.Consume(
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
	go func() {
		fmt.Println("ConsumeCreatePostMessage() start")
		for mes := range messages {
			fmt.Println("mes.Body:", string(mes.Body))
			zap.L().Info("收到消息", zap.String("message", string(mes.Body)))
			//将消息存入数据库
			err = mysql.SendMessage(userID, fan, "新帖子提醒", string(mes.Body))
			if err != nil {
				zap.L().Error("mysql.SendMessage() failed", zap.Error(err))
				return
			}
		}
	}()
	<-forever
}
