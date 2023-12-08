package rabbitmq

import (
	"SHforum_backend/dao/mysql"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type CreatePostMessage struct {
	UserID  int64   `json:"user_id"`
	Fans    []int64 `json:"fans"`
	Message string  `json:"message"`
}

// PublishCreatePostMessage 传递关注的用户创建了贴子的消息
func (r *RabbitMQ) PublishCreatePostMessage(userID int64, fans []int64, message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mes := CreatePostMessage{
		UserID:  userID,
		Fans:    fans,
		Message: message,
	}
	jsonbody, err := json.Marshal(mes)
	//调用channel 发送消息到队列中
	err = r.channel.PublishWithContext(ctx,
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonbody,
		})
	r.failOnErr(err, "Failed to publish a message")
}

// ConsumeCreatePostMessage 消费端收到已创建贴子的消息
func (r *RabbitMQ) ConsumeCreatePostMessage() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
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
