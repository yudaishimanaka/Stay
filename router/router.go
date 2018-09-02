package router

import (
	"net/http"
	"log"
	"os"
	"io"
	"strconv"
	"time"

	"../models"
	"../arpScan"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"github.com/go-xorm/xorm"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

const (
	EventUpdate = iota
)

type Config struct {
	App AppConfig `toml:"app"`
}

type AppConfig struct {
	Interface 	string 		  `toml:"interface"`
	Network   	string 		  `toml:"network"`
	ArpInterval time.Duration `toml:"arp_interval"`
	ArpTimeOut  time.Duration `toml:"arp_timeout"`
}

func Init(engine *xorm.Engine) *gin.Engine {
	// app config parse
	var conf Config
	toml.DecodeFile("config.toml", &conf)

	// routing start
	r := gin.New()
	w := melody.New()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "./assets")
	r.Static("/userIcon", "./userIcon")

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

		userGroup.GET("/viewAll", func(c *gin.Context) {
			var users []models.User
			var response []models.User
			var errMsg string

			err := engine.Find(&users)
			if err != nil {
				errMsg = "query failed (select * from user)"
				c.JSON(http.StatusBadRequest, errMsg)
			} else {
				// mac addr verification
				hwAddrList, _ := arpScan.ArpScan(conf.App.Network, conf.App.Interface, conf.App.ArpTimeOut)
				for _, hwAddr := range hwAddrList {
					for _, user := range users {
						if hwAddr == user.HwAddr {
							response = append(response, user)
						} else {
							continue
						}
					}
				}
				c.JSON(http.StatusCreated, response)
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
							file, _, err := c.Request.FormFile("icon")
							if err != nil {
								file, _ := os.Open("./assets/images/default.jpg")
								out, _ := os.Create("./userIcon/"+user.UserName)
								io.Copy(out, file)
								out.Close()
								file.Close()
							} else {
								out, _ := os.Create("./userIcon/"+user.UserName)
								io.Copy(out, file)
								out.Close()
								file.Close()
							}

							user.IconPath = "userIcon/"+user.UserName

							// insert request
							_, err = engine.Insert(&user)
							if err != nil {
								errMsg = "query failed (insert error)"
								c.JSON(http.StatusBadRequest, errMsg)
							} else {
								c.JSON(http.StatusCreated, "user register successfully")
							}
						}
					}
				}
			}
		})

		userGroup.DELETE("/delete/:userName", func(c *gin.Context) {
			var user models.User
			var errMsg string

			// get request
			userName := c.Params.ByName("userName")
			os.Remove("./userIcon/"+userName)

			_, err := engine.Where("user_name = '"+userName+"'").Delete(&user)
			if err != nil {
				errMsg = "query failed (delete user)"
				c.JSON(http.StatusBadRequest, errMsg)
			} else {
				c.JSON(http.StatusCreated, "user delete successfully")
			}
		})
	}

	// web socket router
	w.HandleConnect(func(s *melody.Session) {
		log.Printf("websocket connection open.")
	})

	w.HandleDisconnect(func(s *melody.Session) {
		log.Printf("websocket connection close.")
	})

	// loop per arp interval
	w.HandleMessage(func(s *melody.Session, msg []byte) {
		for {
			msg = []byte(strconv.Itoa(EventUpdate))
			w.Broadcast(msg)
			time.Sleep(time.Second * conf.App.ArpInterval)
		}
	})

	return r
}
