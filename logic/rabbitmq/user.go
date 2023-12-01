package rabbitmq

import (
	"SHforum_backend/dao/redis"
	"SHforum_backend/util"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// NewRabbitMQSimple 创建RabbitMQ简单模式实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

type SendCodeMessage struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// PublishSendCodeMessage 传递email和code
func (r *RabbitMQ) PublishSendCodeMessage(email string, code string) {
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
	message := SendCodeMessage{
		Email: email,
		Code:  code,
	}
	jsonbody, err := json.Marshal(message)
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

// ConsumeCodeMessage 消费者
func (r *RabbitMQ) ConsumeCodeMessage() {
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

	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			var mes SendCodeMessage
			err = json.Unmarshal(d.Body, &mes)
			err = util.SendEmailWithCode([]string{mes.Email}, mes.Code)
			if err != nil {
				return
			}
			redis.SaveCode(mes.Email, mes.Code)
		}
	}()
}
