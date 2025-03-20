package sugar

import (
	"errors"
	"time"

	"golang.org/x/exp/rand"
)

// ExponentialBackoff 计算指数退避的时间
// @param attempt: 重试次数
// @param baseDelay: 基础等待时间,单位秒
// @param maxDelay: 最大等待时间，单位秒
// @param isJitter: 是否加入随机化以避免重试同步
// @return delay: 指数退避的时间
// @return error: 错误信息
func ExponentialBackoff(attempt int, baseDelay, maxDelay int, isJitter bool) (time.Duration, error) {
	if baseDelay <= 0 || maxDelay <= 0 {
		return 0, errors.New("无效的参数")
	}
	// 计算指数退避的时间
	baseDelayDur := time.Duration(baseDelay) * time.Second
	maxDelayDur := time.Duration(maxDelay) * time.Second
	delay := baseDelayDur * time.Duration(1<<uint(attempt))
	// 防止超过最大延迟时间
	if delay > maxDelayDur {
		return delay, errors.New("退避时间超过了预期的时间")
	}
	if isJitter {
		// 加入随机化以避免重试同步
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay += jitter
	}
	return delay, nil
}
