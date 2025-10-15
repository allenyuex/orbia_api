package utils

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake ID 生成器
// 格式: 41位时间戳 + 10位机器ID + 12位序列号 = 63位（不使用符号位）
const (
	// 时间起始点（2024-01-01 00:00:00 UTC）
	epoch int64 = 1704067200000

	// 机器ID位数
	workerIDBits = 10
	// 序列号位数
	sequenceBits = 12

	// 最大机器ID
	maxWorkerID = -1 ^ (-1 << workerIDBits) // 1023
	// 最大序列号
	maxSequence = -1 ^ (-1 << sequenceBits) // 4095

	// 时间戳左移位数
	timestampShift = workerIDBits + sequenceBits // 22
	// 机器ID左移位数
	workerIDShift = sequenceBits // 12
)

// IDGenerator Snowflake ID 生成器
type IDGenerator struct {
	mu         sync.Mutex
	workerID   int64
	sequence   int64
	lastMillis int64
}

var (
	defaultGenerator *IDGenerator
	once             sync.Once
)

// InitIDGenerator 初始化默认ID生成器
func InitIDGenerator(workerID int64) error {
	if workerID < 0 || workerID > maxWorkerID {
		return fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}

	once.Do(func() {
		defaultGenerator = &IDGenerator{
			workerID: workerID,
		}
	})

	return nil
}

// GetDefaultGenerator 获取默认生成器
func GetDefaultGenerator() *IDGenerator {
	if defaultGenerator == nil {
		// 如果没有初始化，使用默认机器ID 1
		_ = InitIDGenerator(1)
	}
	return defaultGenerator
}

// NextID 生成下一个ID
func (g *IDGenerator) NextID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < g.lastMillis {
		return 0, fmt.Errorf("clock moved backwards")
	}

	if now == g.lastMillis {
		// 同一毫秒内，序列号递增
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for now <= g.lastMillis {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 新的毫秒，序列号归零
		g.sequence = 0
	}

	g.lastMillis = now

	// 组装ID
	id := ((now - epoch) << timestampShift) |
		(g.workerID << workerIDShift) |
		g.sequence

	return id, nil
}

// NextIDString 生成下一个ID字符串
func (g *IDGenerator) NextIDString() (string, error) {
	id, err := g.NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

// GenerateOrderID 生成订单ID（带前缀）
func GenerateOrderID() (string, error) {
	id, err := GetDefaultGenerator().NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ORD%d", id), nil
}

// GenerateID 生成通用ID字符串
func GenerateID() (string, error) {
	return GetDefaultGenerator().NextIDString()
}
