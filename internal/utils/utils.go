package utils

// NormalizeVarArgs 标准化可变参数，根据参数数量返回不同形式的结果
func NormalizeVarArgs[T any](data ...T) any {
	if l := len(data); l == 1 {
		return data[0]
	} else if l > 1 {
		return data
	} else {
		return nil
	}
}

// GetOption 获取可选的参数
func GetOption[T any](opts ...T) (bool, T) {
	if len(opts) == 0 {
		var zero T
		return false, zero
	} else {
		return true, opts[0]
	}
}
