package router

import (
	"net/http"

	"../models"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

func Init(engine *xorm.Engine) *gin.Engine {
	r := gin.New()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{})
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
	}

	return r
}