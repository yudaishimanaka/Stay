package router

import (
	"net/http"
	"log"
	//"encoding/binary"

	"../models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

func Init(engine *xorm.Engine) *gin.Engine {
	r := gin.New()
	w := melody.New()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{})
	})

	r.GET("/ws", func(c *gin.Context) {
		w.HandleRequest(c.Writer, c.Request)
	})

	userGroup := r.Group("/user")
	{
		userGroup.GET("/view/:userName", func(c *gin.Context) {
			userName := c.Params.ByName("userName")
			view := models.User{UserName: userName}
			has, err := engine.Get(&view)
			if err != nil {
				panic(err)
			} else {
				if has == true {
					c.JSON(http.StatusCreated, view)
				} else {
					c.JSON(http.StatusBadRequest, "don't find user")
				}
			}
		})

		userGroup.POST("/register", func(c *gin.Context) {
			var user models.User
			err := c.BindJSON(&user)
			if err != nil {
				c.JSON(http.StatusBadRequest, "has need username")
			} else {
				_, err := engine.Insert(&user)
				if err != nil {
					panic(err)
				} else {
					c.JSON(http.StatusCreated, "user register success")
				}
			}
		})

		userGroup.POST("/update/:userName", func(c *gin.Context) {

		})

		userGroup.DELETE("/delete/:userName", func(c *gin.Context) {

		})
	}

	// web socket router
	w.HandleConnect(func(s *melody.Session) {
		log.Printf("websocket connection open.")
	})

	w.HandleDisconnect(func(s *melody.Session) {
		log.Printf("websocket connection close.")
	})

	w.HandleMessage(func(s *melody.Session, msg []byte) {
		w.Broadcast(msg)
	})

	return r
}
