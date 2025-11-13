package kafka

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventConverter struct{}

func (c *EventConverter) ToMessage(event event.Event) *ckafka.Message {
	topic := event.Topic().String()
	return &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          addEventIDMeta(event.Payload(), event.ID()),
		Timestamp:      event.OccurredAt(),
		Key:            []byte(event.AggregateID().String()),
		Opaque:         event, //供回执读取
	}
}

func (c *EventConverter) ToEvent(message *ckafka.Message) event.Event {
	return event.NewStandardEvent(
		getEventIDMeta(message.Value), // unknown id
		kernel.ID(message.Key),
		message.Timestamp,
		event.Topic(*message.TopicPartition.Topic),
		message.Value,
	)
}

func addEventIDMeta(payload []byte, id event.ID) []byte {
	var data map[string]interface{}

	// 解析原始JSON负载
	if err := json.Unmarshal(payload, &data); err != nil {
		// 如果解析失败，创建一个新的空对象
		data = make(map[string]interface{})
	}

	// 添加或更新重试次数字段
	data["event_id"] = id

	// 将map序列化回JSON
	result, err := json.Marshal(data)
	if err != nil {
		return payload // 失败时返回原负载
	}

	return result
}

func getEventIDMeta(payload []byte) event.ID {
	var data map[string]interface{}

	// 解析JSON负载
	if err := json.Unmarshal(payload, &data); err != nil {
		return "" // 解析失败返回""
	}

	// 检查retry_meta字段是否存在
	retryVal, exists := data["event_id"]
	if !exists {
		return "" // 字段不存在返回""
	}

	// 类型断言，安全地提取整数值
	switch v := retryVal.(type) {
	case string:
		return event.ID(v)
	default:
		return "" // 类型不匹配返回0
	}
}
