package service

import (
	"examresult/model"
	"examresult/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// RegisterModel 是 Register 注册的请求参数
// emmm，其实这就是个 model.Student 😂
type RegisterModel struct {
	School   string `form:"school" json:"school" xml:"school" binding:"required"`
	Sid      string `form:"sid" json:"sid" xml:"sid" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	Email    string `form:"email" json:"email" xml:"email" binding:"required"`
}

// Register 注册一个学生
// 以 school + sid 唯一确定一个学生，
// 首先通过教务系统验证用户存在性，
// 然后查询成绩，以前的就不提醒了
// 最后写数据库：存在则更新信息，不存在新建
//
// 返回保存到数据库中的 Student model，还有一个额外的字符串（这里暂定表示学生名字）
func Register(r RegisterModel) (student model.Student, studentName string, err error) {
	s := model.Student{
		School:   r.School,
		Sid:      r.Sid,
		Password: r.Password,
		Email:    r.Email,
	}

	qzcli, err := auth(s)
	if err != nil {
		return s, "", err
	}

	exams, err := getExamResults(s, qzcli)
	if err != nil {
		return s, "", err
	}

	s, err = save(s, exams)
	return s, qzcli.Auth.UserRealName, err
}

// Register 的 Http 路由:
//    POST /register
//    Multipart/Urlencoded Form:
//        school=ncepu
//        sid=201800000000
//        password=123456
//        email=johnsmith@example.com
//    Response JSON:
//        {"success": "张三"}
//        {"error": "error message"}
// E.g.
//    curl '0.0.0.0:8080/register' -X POST -d 'school=ncepu&sid=201800000000&password=123456&email=johnsmith@example.com'
func initRegister() {
	log.Info("init router /register")
	server.HttpRouter.POST("/register", func(c *gin.Context) {
		var req RegisterModel
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, name, err := Register(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": name})
	})
}

func auth(s model.Student) (*server.QzClient, error) {
	qzcli, err := server.NewQzClient(s)
	if err != nil {
		log.WithFields(log.Fields{
			"s":   s,
			"err": err,
		}).Warn("service register: QZ auth error")
	}
	return qzcli, err
}

func getExamResults(s model.Student, qzcli *server.QzClient) ([]model.ExamResult, error) {
	exams, err := qzcli.GetExamResults()

	// 注册时已有的成绩不予提醒
	for i := 0; i < len(exams); i++ {
		exams[i].Noticed = true
	}

	if err != nil {
		log.WithFields(log.Fields{
			"s":   s,
			"err": err,
		}).Error("service register: QZ get exam results error")
	}
	return exams, err
}

func save(s model.Student, exams []model.ExamResult) (model.Student, error) {
	s, err := saveStudent(s)
	if err != nil {
		return s, err
	}

	saveExamResults(s, exams)

	return s, err
}

func saveStudent(s model.Student) (model.Student, error) {
	s, err := model.SaveStudent(s)
	if err != nil || s.ID == 0 {
		log.WithFields(log.Fields{
			"s":   s,
			"err": err,
		}).Error("service register: SaveStudent error")
	}
	return s, err
}

func saveExamResults(s model.Student, exams []model.ExamResult) {
	for _, result := range exams {
		// 重写学生字段，确保和刚保存到数据库中的一致
		result.StudentID = s.ID
		result.Student = s

		if _, err := model.FirstOrCreateExamResult(result); err != nil {
			log.WithFields(log.Fields{
				"s":    s,
				"exam": result,
				"err":  err,
			}).Error("service register: FirstOrCreateExamResult error")
		}
	}
}
