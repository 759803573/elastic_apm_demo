package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgin"
)

var db = make(map[string]string)

func randomSleep(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Millisecond)
}

func randomPanic(n int, err error) {
	if err == nil {
		err = fmt.Errorf("Default Error")
	}
	if rand.Intn(100) > n {
		panic(err)
	}
}

func mockDbCall(tx *apm.Transaction) {
	dbSpan := tx.StartSpan("Mock db query", "db.mysql.query", nil)
	defer func() { dbSpan.End() }()

	randomSleep(100)
	randomPanic(80, fmt.Errorf("DB Error"))
}

func mock_exteral_sevice_call(tx *apm.Transaction) {
	dbSpan := tx.StartSpan("Mock Service Call", "service.mock.call", nil)
	defer func() { dbSpan.End() }()

	randomSleep(50)
	dbInnerSpan := tx.StartSpan("Mock Service db query", "db.mysql.query", dbSpan)
	randomSleep(100)
	dbInnerSpan.End()

	randomSleep(100)
	randomPanic(80, fmt.Errorf("DB Error"))
}

func setupRouter() *gin.Engine {
	r := gin.New()

	r.Use(apmgin.Middleware(r))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/user/:name", func(c *gin.Context) {
		tx := apm.TransactionFromContext(c.Request.Context())
		user := c.Params.ByName("name")
		value := rand.Intn(100)

		mockDbCall(tx)
		mock_exteral_sevice_call(tx)

		if value > 50 {
			tx.Context.SetTag("TagTest", string(value))
			tx.Context.SetTag("TagTest1", string(value))
			tx.Context.SetTag("TagTest2", string(value))
			tx.Context.SetUsername(user)
			tx.Context.SetUserID(string(value))
			tx.Context.SetUserEmail(user + "@alo7.com")
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			tx.Context.SetTag("TagTest1", string(value))
			tx.Context.SetTag("TagTest2", string(value))
			tx.Context.SetTag("TagTest3", string(value))
			c.JSON(http.StatusForbidden, gin.H{"user": user, "status": "no value"})
		}
	})

	return r
}

func main() {
	db["apm"] = "mpa"

	s := &http.Server{
		Addr:           ":8081",
		Handler:        setupRouter(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
