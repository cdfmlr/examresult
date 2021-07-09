package service

import log "github.com/sirupsen/logrus"

func Init() {
	initRegister()
	initRegisterFrontEnd()
	initExamResultTicker()
	initNoticerHandler()
	initEmailNoticer()
}

func Run() {
	log.Info("ExamResultTicker start")
	ExamResultTicker.Start()
}
