package service

import (
	"examresult/model"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

// Noticer 是一个通知渠道：
// 调用 Notice 方法进行通知
//
// 除了全局的 noticeHandler 单例，其他的任何 Noticer 不可写数据库，只能用给定的 result，不可更改其中内容然后保存。
type Noticer interface {
	Notice(result model.ExamResult) bool // 对成绩进行进行通知，如通知成功则返回 true
}

//////////////////////
// 👇全局的通知处理器  //
/////////////////////

// noticeHandler 将成绩通知分发给各个具体的通知器
type noticeHandler struct {
	noticers []Noticer
}

// Notice 分发成绩通知，然后写数据库标记已通知
func (n noticeHandler) Notice(result model.ExamResult) bool {
	chOk := make(chan bool, 3)
	var workers int32 = int32(len(n.noticers))

	for _, n := range n.noticers {
		n := n
		go func() {
			chOk <- n.Notice(result)
		}()
	}

	for ok := range chOk {
		atomic.AddInt32(&workers, -1)
		if ok {
			result.Noticed = true
			model.DB.Save(&result)
			if atomic.LoadInt32(&workers) <= 0 {
				break
			}
		}
	}

	lg := log.WithFields(log.Fields{
		"ExamResultID": result.ID,
		"StudentID":    result.StudentID,
	})
	if result.Noticed {
		lg.Info("service noticeHandler: notice success")
	} else {
		lg.Warn("service noticeHandler: notice failed")
	}

	return result.Noticed
}

// AddNoticer 注册一个新的通知渠道
func (n *noticeHandler) AddNoticer(nn Noticer) {
	n.noticers = append(n.noticers, nn)
}

// NoticeHandler 全局通知处理器
// 用来将一个成绩通知分发给各个具体的通知器
var NoticeHandler *noticeHandler

func initNoticerHandler() {
	log.Info("init NoticeHandler")
	NoticeHandler = &noticeHandler{}
}
