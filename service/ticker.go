package service

import (
	"examresult/config"
	"examresult/model"
	"examresult/server"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// examResultTicker 定期从教务系统查成绩。
// 更新数据库，发现新的则进行通知
type examResultTicker struct {
	ticker *time.Ticker

	PeriodInMinute int           // 运行周期，多久检查一次
	done           chan struct{} // 用来停止周期运作的 chan
}

func NewExamResultTicker(periodInMinute int) *examResultTicker {
	return &examResultTicker{
		PeriodInMinute: periodInMinute,
		done:           make(chan struct{}),
		ticker:         time.NewTicker(time.Duration(periodInMinute) * time.Minute),
	}
}

// Start 让 examResultTicker 开始周期性工作
func (t *examResultTicker) Start() {
	go func() {
		// 调用时立即执行一次
		time.Sleep(time.Duration(
			int(config.ExamResultQuery.MinSleepSeconds)) * time.Second)
		t.work()

		for {
			select {
			case <-t.done:
				log.Info("examResultTicker Stop!")
				return
			case <-t.ticker.C: // 检查出成绩没
				t.work()
			}
		}
	}()
}

// Stop 停止 CourseTicker
func (t *examResultTicker) Stop() {
	t.done <- struct{}{}
}

// work 就是周期 Tick 了要做的事：查成绩，如有新成绩则进行通知
func (t *examResultTicker) work() {
	log.Info("service examResultTicker: work")

	students, err := model.GetAllStudents()
	if err != nil {
		log.WithError(err).Error("service examResultTicker: GetAllStudents failed")
		return
	}

	for _, s := range students {
		fetchExamResultAndNotice(s)
		// 睡眠一段时间
		minSleep := int(config.ExamResultQuery.MinSleepSeconds)
		randMax := t.PeriodInMinute / len(students) * 30
		time.Sleep(time.Duration(
			minSleep+rand.Intn(randMax)) * time.Second)
	}
}

func fetchExamResultAndNotice(student model.Student) {
	cli, err := server.NewQzClient(student)
	if err != nil {
		log.WithError(err).Warn("service examResultTicker: login Qz failed")
		return
	}

	results, err := cli.GetExamResults()
	if err != nil {
		log.WithError(err).Error("service examResultTicker: GetExamResults failed")
		return
	}

	for _, result := range results {
		result, err = model.FirstOrCreateExamResult(result)
		if err != nil {
			log.WithError(err).Error("service examResultTicker: save exam failed")
		}
		// 没通知就通知
		if result.Noticed == false {
			NoticeHandler.Notice(result)
		}
	}
}

// ExamResultTicker 定期从教务系统查成绩。
//
// MAIN: 这个东西要在 main 里 ExamResultTicker.Start()
var ExamResultTicker *examResultTicker

func initExamResultTicker() {
	log.Info("init ExamResultTicker")
	ExamResultTicker = NewExamResultTicker(int(config.ExamResultQuery.PeriodInMinutes))
}
