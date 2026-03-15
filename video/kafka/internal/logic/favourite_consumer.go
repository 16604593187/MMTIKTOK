package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mini-tiktok/video/kafka/internal/svc"
	"mini-tiktok/video/model"
	"mini-tiktok/video/model/KafkaMessage"
)

type FavoriteConsumer struct {
	reader *kafka.Reader
	svcCtx *svc.ServiceContext
}

func NewFavoriteConsumer(svcCtx *svc.ServiceContext) *FavoriteConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{svcCtx.Config.KafkaConfig.Host},
		Topic:    svcCtx.Config.KafkaConfig.Topic,
		GroupID:  "favorite_consumer_group",
		MinBytes: svcCtx.Config.KafkaConfig.MinBytes,
		MaxBytes: svcCtx.Config.KafkaConfig.MaxBytes,
		MaxWait:  time.Second * 1,
	})

	return &FavoriteConsumer{
		reader: reader,
		svcCtx: svcCtx,
	}
}

func (c *FavoriteConsumer) Start() {
	logx.Info("点赞批量落盘消费者已启动...")
	ctx := context.Background()

	// 🌟 核心改动 1：将批处理池分为“待插入”和“待删除”两组
	var insertRecords []*model.Favorite
	var deleteRecords []*model.Favorite 
	var batchMessages []kafka.Message

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		fetchCtx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
		m, err := c.reader.FetchMessage(fetchCtx)
		cancel()

		if err == nil {
			var msgInfo KafkaMessage.MsgInfo
			if err := json.Unmarshal(m.Value, &msgInfo); err == nil && msgInfo.Model == "favorite" {

				cols := strings.Split(msgInfo.Columns, ",")

				if len(cols) >= 3 { 
					userId, _ := strconv.ParseUint(cols[0], 10, 64)
					videoId, _ := strconv.ParseUint(cols[1], 10, 64)
					createTime, _ := strconv.ParseInt(cols[2], 10, 64)

					favRecord := &model.Favorite{
						UserId:     userId,
						VideoId:    videoId,
						CreateTime: createTime,
					}

					//根据 Op 动作将记录分流
					if msgInfo.Op == "delete" {
						deleteRecords = append(deleteRecords, favRecord)
					} else {
						insertRecords = append(insertRecords, favRecord)
					}
					
					batchMessages = append(batchMessages, m)
				}
			}
		}

		// 触发落盘条件：总消息数 >= 100 或者 触发 1 秒定时器
		totalLen := len(insertRecords) + len(deleteRecords)
		select {
		case <-ticker.C:
			if totalLen > 0 {
				c.flushAndCommit(ctx, &insertRecords, &deleteRecords, &batchMessages)
			}
		default:
			if totalLen >= 100 {
				c.flushAndCommit(ctx, &insertRecords, &deleteRecords, &batchMessages)
			}
		}
	}
}

// flushAndCommit 在一个数据库事务中处理批量插入和批量删除
func (c *FavoriteConsumer) flushAndCommit(ctx context.Context, inserts *[]*model.Favorite, deletes *[]*model.Favorite, msgs *[]kafka.Message) {
	// 开启事务，保证这一批次的增删操作要么全成功，要么全失败
	err := c.svcCtx.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		
		//处理批量插入 (使用 IGNORE 避免重复点赞报错)
		if len(*inserts) > 0 {
			err := tx.Clauses(clause.OnConflict{
				DoNothing: true, // 如果联合索引冲突，直接忽略，不做任何操作
			}).CreateInBatches(*inserts, len(*inserts)).Error
			if err != nil {
				return err
			}
		}

		//处理批量物理删除
		// GORM 针对联合主键的批量删除较弱，这里由于有事务包裹，循环删除速度依然极快
		if len(*deletes) > 0 {
			for _, d := range *deletes {
				// 执行 DELETE FROM favorite WHERE user_id = ? AND video_id = ?
				if err := tx.Where("user_id = ? AND video_id = ?", d.UserId, d.VideoId).Delete(&model.Favorite{}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		logx.Errorf("批量落盘 MySQL 失败: %v", err)
		return // 发生错误，不提交 Offset
	}

	// 提交 Kafka Offset
	if err := c.reader.CommitMessages(ctx, *msgs...); err != nil {
		logx.Errorf("MySQL 写入成功，但提交 Kafka Offset 失败: %v", err)
	} else {
		logx.Infof("成功落盘并提交 Offset: 插入 %d 条, 删除 %d 条", len(*inserts), len(*deletes))
	}

	// 清空切片，复用内存
	*inserts = (*inserts)[:0]
	*deletes = (*deletes)[:0]
	*msgs = (*msgs)[:0]
}