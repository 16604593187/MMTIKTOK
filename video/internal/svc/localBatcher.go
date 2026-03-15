package svc

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// LocalBatcher 本地微批处理
type LocalBatcher struct {
	msgChan   chan kafka.Message // 无锁接收通道
	writer    *kafka.Writer      // Kafka 底层写入器
	batchSize int                // 聚合阈值：满多少条发一次
	linger    time.Duration      // 时间阈值：最多等多久发一次
	ctx       context.Context
}

// NewLocalBatcher 初始化并启动后台发射协程
func NewLocalBatcher(ctx context.Context, writer *kafka.Writer) *LocalBatcher {
	b := &LocalBatcher{
		// 初始化一个容量为 5000 的 Channel
		msgChan:   make(chan kafka.Message, 5000),
		writer:    writer,
		batchSize: 100,             // 100 条发一次
		linger:    time.Second * 1, //  1 秒发一次
		ctx:       ctx,
	}

	go b.startLoop()
	return b
}

// Push 业务层调用的非阻塞投递接口
func (b *LocalBatcher) Push(msg kafka.Message) error {
	select {
	case b.msgChan <- msg:
		return nil // 投递成功
	default:
		err := errors.New("本地 Kafka 缓冲队列已满，消息被降级丢弃")
		logx.Errorf("严重警告：%v", err)
		return err // 返回错误，通知业务层进行补偿回滚
	}
}

//后台定时打包与发送
func (b *LocalBatcher) startLoop() {
	var batch []kafka.Message
	ticker := time.NewTicker(b.linger)
	defer ticker.Stop()

	for {
		select {
		case msg := <-b.msgChan:
			batch = append(batch, msg)
			//聚合阈值
			if len(batch) >= b.batchSize {
				b.flush(batch)
				batch = make([]kafka.Message, 0, b.batchSize)
			}
		case <-ticker.C:
			//时间
			if len(batch) > 0 {
				b.flush(batch)
				batch = make([]kafka.Message, 0, b.batchSize)
			}
		case <-b.ctx.Done():
			// 处理剩余
			if len(batch) > 0 {
				b.flush(batch)
			}
			return
		}
	}
}

func (b *LocalBatcher) flush(batch []kafka.Message) {
	err := b.writer.WriteMessages(context.Background(), batch...)
	if err != nil {
		logx.Errorf("批量写入 Kafka 失败: %v", err)
	} else {
		logx.Infof("成功将 %d 条互动消息打包请求送入 Kafka！", len(batch))
	}
}
