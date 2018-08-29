package router

import (
	"net/http"
	"log"
	"os"
	"io"

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
			var errMsg string
			user := models.User{}
			userName := c.Params.ByName("userName")
			has, err := engine.Where("user_name = '"+userName+"'").Get(&user)
			if err != nil {
				errMsg = "query failed (find user)"
				c.JSON(http.StatusBadRequest, errMsg)
			} else {
				if has == true {
					c.JSON(http.StatusCreated, user)
				} else {
					errMsg = "user not found"
					c.JSON(http.StatusBadRequest, errMsg)
				}
			}
		})

		userGroup.POST("/register", func(c *gin.Context) {
			var user models.User
			var errMsg string

			// bind request
			err := c.Bind(&user)
			if err != nil {
				errMsg = "query failed (can not bind request)"
				c.JSON(http.StatusBadRequest, errMsg)
			} else {
				// check user same name
				checkExistenceUser := "select count(*) from user where user_name = '"+user.UserName+"'"
				counts, err := engine.SQL(checkExistenceUser).Count()
				if err != nil {
					errMsg = "query failed (select error)"
					c.JSON(http.StatusBadRequest, errMsg)
				} else {
					if counts > 0 {
						errMsg = "a user with the same name already exists, please change it to a different name"
						c.JSON(http.StatusBadRequest, errMsg)
					} else {
						// hwAddr null check
						if user.HwAddr == "" {
							errMsg = "insert failed (mac_address is null)"
							c.JSON(http.StatusBadRequest, errMsg)
						} else {
							// save user icon
							file, header, _ := c.Request.FormFile("icon")
							file.Close()
							out, _ := os.Create("./userIcon/"+header.Filename)
							defer out.Close()
							io.Copy(out, file)
							user.IconPath = "userIcon/"+header.Filename

							// insert request
							_, err = engine.Insert(&user)
							if err != nil {
								errMsg = "query failed (insert error)"
								c.JSON(http.StatusBadRequest, errMsg)
							} else {
								c.JSON(http.StatusCreated, "user register success")
							}
						}
					}
				}
			}
		})

		userGroup.POST("/iconUpdate/:userName", func(c *gin.Context) {
			// receive and get file
			file, err := c.FormFile("icon")
			if err != nil {
				log.Printf("receive error")
			}

			log.Println(file)
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
