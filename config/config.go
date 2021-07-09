package config

// ConfDatabase is a struct for Database configures. Provides DSN (Data Source Name):
//    [username[:password]@][protocol[(address)]]/dbname
// Refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details.
type ConfDatabase struct {
	Username string
	Password string
	Protocol string
	Address  string
	DBName   string
}

// ConfHttpServer is a struct for serving configures.
type ConfHttpServer struct {
	HttpAddress string
}

type ConfExamResultQuery struct {
	PeriodInMinutes uint // 多久查一次成绩
	MinSleepSeconds uint // 查询两个人之间至少隔多久
}

// ConfSMTP 就是 SMTP 配置，和 server.SMTP 一摸一样，这里就直接取个别名了
type ConfSMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Conf is a struct wraps all configures.
// field XXX -> <type ConfXXX struct>
type Conf struct {
	Database        ConfDatabase
	ExamResultQuery ConfExamResultQuery
	HttpServer      ConfHttpServer
	SMTP            ConfSMTP
}
