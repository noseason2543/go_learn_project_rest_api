package logger

import (
	"fmt"
	"go_learn_project_rest_api/pkgs/utils"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ILogger interface {
	Print() ILogger
	Save()
	SetQuery(c fiber.Ctx)
	SetBody(c fiber.Ctx)
	SetResponse(res any)
}

type Logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitLogger(c fiber.Ctx, res any) ILogger {
	log := &Logger{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		StatusCode: c.Response().StatusCode(),
		Path:       c.Path(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)

	return log
}

func (l *Logger) Print() ILogger {
	utils.Debug(l)
	return l
}

func (l *Logger) Save() {
	data := utils.Output(l)
	fileName := fmt.Sprintf("./assets/logs/log_project_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error when opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

func (l *Logger) SetQuery(c fiber.Ctx) {
	var query any
	if err := c.Bind().Query(&query); err != nil {
		log.Printf("query parser error: %v", err)
	}
	l.Query = query
}

func (l *Logger) SetBody(c fiber.Ctx) {
	var body any
	if err := c.Bind().Body(&body); err != nil {
		log.Printf("body parser error: %v", err)
	}

	switch l.Path {
	case "v1/users/signup":
		l.Body = "never give up"
	default:
		l.Body = body
	}
}

func (l *Logger) SetResponse(res any) {
	l.Response = res
}
