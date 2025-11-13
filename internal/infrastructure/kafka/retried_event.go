package kafka

import (
	"encoding/json"
	"gochat/internal/shared/event"
)

type RetriedEvent struct {
	retryTimes int
	event.Event
}

func wrapEventWithRetryMeta(ev event.Event, retry int) event.Event {
	return &RetriedEvent{
		retryTimes: retry,
		Event:      ev,
	}
}

func (e *RetriedEvent) Payload() []byte {
	return addRetryMeta(e.Event.Payload(), e.retryTimes)
}

// addRetryMeta 向JSON负载中添加重试元数据字段
func addRetryMeta(payload []byte, retry int) []byte {
	var data map[string]interface{}

	// 解析原始JSON负载
	if err := json.Unmarshal(payload, &data); err != nil {
		// 如果解析失败，创建一个新的空对象
		data = make(map[string]interface{})
	}

	// 添加或更新重试次数字段
	data["retry_times"] = retry

	// 将map序列化回JSON
	result, err := json.Marshal(data)
	if err != nil {
		return payload // 失败时返回原负载
	}

	return result
}

func getRetryMeta(payload []byte) int {
	var data map[string]interface{}

	// 解析JSON负载
	if err := json.Unmarshal(payload, &data); err != nil {
		return 0 // 解析失败返回0
	}

	// 检查retry_meta字段是否存在
	retryVal, exists := data["retry_times"]
	if !exists {
		return 0 // 字段不存在返回0
	}

	// 类型断言，安全地提取整数值
	switch v := retryVal.(type) {
	case float64: // JSON数字默认解析为float64
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case float32:
		return int(v)
	default:
		return 0 // 类型不匹配返回0
	}
}
