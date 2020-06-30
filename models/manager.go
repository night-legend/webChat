package models

import (
	"fmt"
	"time"
	"webchat/common"
	"webchat/database"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"github.com/docker/docker/pkg/pubsub"
)

var (
	buffer  = 1024
	timeout = time.Second * 5
)

type DataCenterManager struct {
	Count     int
	TimeOut   time.Duration
	Redis     *redis.Client
	Publisher *pubsub.Publisher
}

func NewDataCenterManager() *DataCenterManager {
	database.DB.AutoMigrate(&SessionMessage{})

	if redisClient, err := NewRedisClient(); err != nil {
		panic(err)
	} else {
		return &DataCenterManager{
			Count:     10,
			TimeOut:   time.Hour,
			Redis:     redisClient,
			Publisher: pubsub.NewPublisher(timeout, buffer),
		}
	}
	return nil
}

// msg SessionMessage
func (dataCenter *DataCenterManager) HandlerMessage(ctx *gin.Context, userID string) error {
	var msg SessionMessage
	if err := ctx.ShouldBind(&msg); err != nil {
		fmt.Println("parse error", err)
	}

	dataCenter.Distribution(msg)
	// dataCenter.Save(msg)
	return nil
}

func (dataCenter *DataCenterManager) Distribution(msg SessionMessage) error {
	identify := msg.getUserIdentify()

	var user User
	err := dataCenter.Redis.Get(identify).Scan(&user)

	fmt.Printf("+++++++++identify: %+v\n, err: %+v, user: %+v\n", user, err)

	ws, ok := ManageEnv.WebsocketManager.Connects[identify]

	if err == nil && ok {
		ws.WriteJSON(msg.MessageBody)
	} else {
		fmt.Println("not found ws")
	}
	return nil
}

func (dataCenter *DataCenterManager) Save(msg SessionMessage) error {
	var message SessionMessage

	if err := database.DB.Where("source_id = ? and destination_id = ?", msg.SourceID, msg.DestinationID).Find(&message).Error; err != nil {
		message.SourceID = msg.SourceID
		message.DestinationID = msg.DestinationID
		message.MessageBody = msg.MessageBody

		if err := database.DB.Create(&message).Error; err != nil {
			return err
		}
	}

	message.MessageBody = append(message.MessageBody, msg.MessageBody...)
	database.DB.Model(&message).Update("message_body", message.MessageBody)
	return nil
}

func (*DataCenterManager) GetMessage(ctx *gin.Context, userID, destID string) error {
	var messages []*SessionMessage

	if err := database.DB.Where("source_id = ? and dest_id = ?", userID, destID).Find(&messages).Error; err != nil {
		return err
	}
	common.HttpSuccessResponse(ctx, messages)
	return nil
}
