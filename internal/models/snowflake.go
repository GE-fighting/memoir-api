package models

import (
	"fmt"
	"sync"
	"time"
)

const (
	// 2024-01-01 00:00:00 UTC作为起始时间
	epoch int64 = 1704067200000

	// 位分配
	nodeBits     uint8 = 10                        // 节点ID占用10位
	sequenceBits uint8 = 12                        // 序列号占用12位
	nodeMax      int64 = -1 ^ (-1 << nodeBits)     // 节点ID最大值(1023)
	sequenceMax  int64 = -1 ^ (-1 << sequenceBits) // 序列号最大值(4095)

	// 时间戳、节点ID、序列号的移位
	timestampShift uint8 = nodeBits + sequenceBits // 时间戳左移22位
	nodeShift      uint8 = sequenceBits            // 节点ID左移12位
)

// Snowflake 雪花算法ID生成器
type Snowflake struct {
	mu        sync.Mutex // 互斥锁
	timestamp int64      // 上次生成ID的时间戳
	node      int64      // 节点ID
	sequence  int64      // 当前序列号
}

// NewSnowflake 创建一个新的雪花算法ID生成器
func NewSnowflake(node int64) (*Snowflake, error) {
	if node < 0 || node > nodeMax {
		return nil, fmt.Errorf("node ID 必须在 0-%d 之间", nodeMax)
	}

	return &Snowflake{
		timestamp: 0,
		node:      node,
		sequence:  0,
	}, nil
}

// NextID 生成下一个ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前时间戳
	now := time.Now().UnixNano() / 1000000 // 转为毫秒
	if now < epoch {
		now = epoch
	}

	// 处理时钟回拨问题
	if s.timestamp > now {
		// 等待至上次生成ID的时间
		time.Sleep(time.Duration(s.timestamp-now) * time.Millisecond)
		now = s.timestamp
	}

	// 同一毫秒内，序列号递增
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMax
		// 序列号溢出，等待下一毫秒
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同毫秒内，序列号重置
		s.sequence = 0
	}

	s.timestamp = now

	// 组装ID（时间戳 | 节点ID | 序列号）
	id := ((now - epoch) << timestampShift) | (s.node << nodeShift) | s.sequence

	return id
}

// 全局默认实例
var (
	defaultNode, _ = NewSnowflake(1) // 默认节点ID为1
)

// GenerateID 使用默认生成器生成ID
func GenerateID() int64 {
	return defaultNode.NextID()
}
