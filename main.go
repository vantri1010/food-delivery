package main

import (
	"food-delivery/component/appctx"
	"food-delivery/component/uploadprovider"
	"food-delivery/middleware"
	"food-delivery/pubsub/localpb"
	"food-delivery/skio"
	"food-delivery/subscriber"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id"` // tag
	Name string `json:"name" gorm:"column:name"`
	Addr string `json:"addr" gorm:"column:addr"`
}

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name"`
	Addr *string `json:"addr" gorm:"column:addr"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

func (RestaurantUpdate) TableName() string {
	return "restaurants"
}

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := os.Getenv("MYSQL_CONN_STRING")

	s3BucketName := os.Getenv("S3BucketName")
	s3Region := os.Getenv("S3Region")
	s3APIKey := os.Getenv("S3APIKey")
	s3SecretKey := os.Getenv("S3SecretKey")
	s3Domain := os.Getenv("S3Domain")
	secretKey := os.Getenv("SYSTEM_SECRET")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)
	ps := localpb.NewPubsub()
	appContext := appctx.NewAppContext(db, s3Provider, secretKey, ps)

	// Setup subscribers
	//subscriber.Setup(appContext, context.Background())
	_ = subscriber.NewEngine(appContext).Start()

	r := gin.Default()
	r.Use(middleware.Recover(appContext))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Static("/static", "./static")

	// POST /restaurants
	v1 := r.Group("/v1")

	setupRoute(appContext, v1)
	setupAdminRoute(appContext, v1)

	rtEngine := skio.NewEngine()
	appContext.SetRealtimeEngine(rtEngine)
	_ = rtEngine.Run(appContext, r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//func startSocketIOServer(engine *gin.Engine) {
//	server, _ := socketio.NewServer(&engineio.Options{
//		Transports: []transport.Transport{websocket.Default},
//	})
//
//	server.OnConnect("/", func(s socketio.Conn) error {
//		s.SetContext("")
//		fmt.Println("Socket connected:", s.ID(), " IP:", s.RemoteAddr())
//
//		go func() {
//			i := 0
//			for {
//				i++
//				s.Emit("test", i)
//				time.Sleep(time.Second)
//
//				if i == 10 {
//					break
//				}
//			}
//		}()
//		return nil
//	})
//
//	server.OnEvent("/", "authenticate", func(s socketio.Conn, token string) {
//		// Validate token
//		// If false: s.Close(), and return
//
//		// If true
//		// => UserId
//		// Fetch db find user by Id
//		// Here: s belongs to who? (user_id)
//		// We need a map[user_id][]socketio.Conn
//		log.Println("socket and token:", s.ID(), token)
//	})
//
//	type A struct {
//		Age int `json:"age"`
//	}
//
//	server.OnEvent("/", "notice", func(s socketio.Conn, msg A) {
//		fmt.Println("notice:", msg.Age)
//		s.Emit("reply", A{msg.Age + 1})
//	})
//
//	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
//		s.SetContext(msg)
//		fmt.Println("msg:", msg)
//		return "recv " + msg
//	})
//
//	server.OnEvent("/", "bye", func(s socketio.Conn) string {
//		last := s.Context().(string)
//		s.Emit("bye", last)
//		s.Close()
//		return last
//	})
//
//	server.OnError("/", func(s socketio.Conn, e error) {
//		fmt.Println("meet error:", e)
//	})
//
//	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
//		fmt.Println("socket closed:", reason)
//		// Remove socket from socket engine (from app context)
//	})
//
//	go server.Serve()
//
//	engine.GET("/socket.io/*any", gin.WrapH(server))
//	engine.POST("/socket.io/*any", gin.WrapH(server))
//
//	engine.StaticFile("/demo/", "./demo.html")
//}
