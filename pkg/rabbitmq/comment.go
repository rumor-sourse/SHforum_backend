package rabbitmq

import (
	"SHforum_backend/interval/dao/mysql"
	"SHforum_backend/interval/models"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"time"
)

// nolint
type CreateCommentMessage struct {
	models.Comment
	models.Post
}

type UpdateCommentLikeMessge struct {
	CommentID string
	LikeCount int64
}

// PublishCreateCommentMessage 传递评论信息
func (r *RabbitMQ) PublishCreateCommentMessage(p *models.Comment, post *models.Post) {
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
	message := CreateCommentMessage{
		Comment: *p,
		Post:    *post,
	}
	jsonbody, err := json.Marshal(message)
	if err != nil {
		return
	}
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

// ConsumeCreateCommentMessageByMysql 消费者
func (r *RabbitMQ) ConsumeCreateCommentMessageByMysql() {
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
	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true, 表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列是否阻塞
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")
	//启用协程处理消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var mes CreateCommentMessage
			err = json.Unmarshal(d.Body, &mes)
			if mes.ParentCommentID == -1 {
				//顶级评论，给贴子作者发消息
				if !mes.IsAdminComment {
					//非贴子作者才发消息
					content := fmt.Sprintf("用户%s给你的帖子《%s》评论了", mes.UserName, mes.Title)
					err := mysql.SendMessage(mes.UserID, mes.AuthorID, "新评论提醒", content)
					if err != nil {
						return
					}
				}
			} else {
				//非顶级评论，给回复的作者发消息
				parentcomment, err := mysql.GetCommentByID(mes.ParentCommentID)
				if err != nil {
					return
				}
				content := fmt.Sprintf("用户%s回复了你的评论:%s", mes.UserName, parentcomment.Content)
				err = mysql.SendMessage(mes.UserID, parentcomment.UserID, "新回复提醒", content)
				if err != nil {
					return
				}
			}
		}
	}()
	<-forever
}

// PublishUpdateCommentLikeMessge 传递更新评论点赞信息
func (r *RabbitMQ) PublishUpdateCommentLikeMessge(CommentID string, likeCount int64) {
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
	mes := UpdateCommentLikeMessge{
		CommentID: CommentID,
		LikeCount: likeCount,
	}
	jsonbody, err := json.Marshal(mes)
	if err != nil {
		return
	}
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

// ConsumeUpdateCommentLikeMessgeByMysql 消费者
func (r *RabbitMQ) ConsumeUpdateCommentLikeMessgeByMysql() {
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
	//消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true, 表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列是否阻塞
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")
	//启用协程处理消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//更新评论点赞数
			var mes UpdateCommentLikeMessge
			err = json.Unmarshal(d.Body, &mes)
			if err != nil {
				return
			}
			//转换为int64
			id, err := strconv.ParseInt(mes.CommentID, 10, 64)
			if err != nil {
				return
			}
			err = mysql.UpdateCommentLikeCount(id, mes.LikeCount)
			if err != nil {
				return
			}
		}
	}()
	<-forever
}
