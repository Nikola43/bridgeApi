package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/antoniodipinto/ikisocket"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nikola43/bridgeApi/controllers"
	//database "github.com/nikola43/bridgeApi/database"
	//middlewares "github.com/nikola43/bridgeApi/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var httpServer *fiber.App

type App struct{}

func (a *App) Initialize(port string) {
	/*
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}

		PROD := os.Getenv("PROD")

		MYSQL_USER := os.Getenv("MYSQL_USER")
		MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
		MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")

		S3_ACCESS_KEY := os.Getenv("S3_ACCESS_KEY")
		S3_SECRET_KEY := os.Getenv("S3_SECRET_KEY")
		S3_ENDPOINT := os.Getenv("S3_ENDPOINT")
		S3_BUCKET_NAME := os.Getenv("S3_BUCKET_NAME")

		X_API_KEY := os.Getenv("X_API_KEY")
		API_CHAINSTACK_URL := os.Getenv("API_CHAINSTACK_URL")
		API_KEY_CHAINSTACK := os.Getenv("API_KEY_CHAINSTACK")

		if PROD == "0" {
			MYSQL_USER = os.Getenv("MYSQL_USER_DEV")
			MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD_DEV")
			MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE_DEV")

			S3_ACCESS_KEY = os.Getenv("S3_ACCESS_KEY_DEV")
			S3_SECRET_KEY = os.Getenv("S3_SECRET_KEY_DEV")
			S3_ENDPOINT = os.Getenv("S3_ENDPOINT_DEV")
			S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME_DEV")

			X_API_KEY = os.Getenv("X_API_KEY_DEV")
			API_CHAINSTACK_URL = os.Getenv("API_CHAINSTACK_URL")
			API_KEY_CHAINSTACK = os.Getenv("API_KEY_CHAINSTACK")

		}
		_ = S3_SECRET_KEY
		_ = S3_ACCESS_KEY
		_ = S3_ENDPOINT
		_ = X_API_KEY
		_ = S3_BUCKET_NAME
		_ = API_KEY_CHAINSTACK
		_ = API_CHAINSTACK_URL

		InitializeDatabase(
			MYSQL_USER,
			MYSQL_PASSWORD,
			MYSQL_DATABASE)

		//database.Migrate()
		//fakedatabase.CreateFakeData()
	*/

	e := InitializeHttpServer().Listen(port)
	if e != nil {
		log.Fatal(e)
	}
}

func (a *App) InitializeWeb3() {

	rpcUrl := "wss://api.avax-test.network/ext/bc/C/ws"
	privateKey := "d4e91ac43134265cc9d905e04be7db37329dc2dddcf69bbdeef5543dc05c0651"

	//todo add rpc

	out := make(chan string)
	web3manager.Web3ManagerInstance = web3manager.NewWsWeb3Client(
		rpcUrl,
		privateKey)

	contractAddress1 := "0x9e20af05ab5fed467dfdd5bb5752f7d5410c832e"

	var addresses []string
	addresses = append(addresses, contractAddress1)

	err := web3manager.Web3ManagerInstance.ListenBridgesEventsV2(addresses, out)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleRoutes(api fiber.Router) {
	//api.Use(middleware.Logger())
	routes.AuthRoutes(api)
	routes.NodesRoutes(api)
}

func InitializeHttpServer() *fiber.App {
	httpServer = fiber.New(fiber.Config{
		BodyLimit: 2000 * 1024 * 1024, // this is the default limit of 4MB
	})

	httpServer.Use(middlewares.XApiKeyMiddleware)

	httpServer.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	httpServer.Get("/", controllers.ON)

	ws := httpServer.Group("/ws")

	// Setup the middleware to retrieve the data sent in first GET request
	ws.Use(middlewares.WebSocketUpgradeMiddleware)

	// Pull out in another function
	// all the ikisocket callbacks and listeners
	setupSocketListeners()

	ws.Get("/:id", ikisocket.New(func(kws *ikisocket.Websocket) {
		models.SocketInstance = kws

		// Retrieve the user id from endpoint
		userId := kws.Params("id")

		// Add the connection to the list of the connected clients
		// The UUID is generated randomly and is the key that allow
		// ikisocket to manage Emit/EmitTo/Broadcast
		models.SocketClients[userId] = kws.UUID

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("user_id", userId)

		//Broadcast to all the connected users the newcomer
		// kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true)
		//Write welcome message
		kws.Emit([]byte(fmt.Sprintf("Socket connected")))
	}))

	fmt.Println(color.YellowString("  ----------------- Websockets -----------------"))
	fmt.Println(color.CyanString("\t    Websocket URL: "), color.GreenString("ws://127.0.0.1:3000/ws"))

	api := httpServer.Group("/api") // /api
	v1 := api.Group("/v1")          // /api/v1
	HandleRoutes(v1)

	/*
		err := httpServer.Listen(port)
		if err != nil {
			log.Fatal(err)
		}
	*/

	return httpServer
}

func InitializeDatabase(user, password, database_name string) {
	connectionString := fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		user,
		password,
		database_name,
	)

	DB, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	database.GormDB, err = gorm.Open(mysql.New(mysql.Config{Conn: DB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}
}

// Setup all the ikisocket listeners
// pulled out main function
func setupSocketListeners() {

	// Multiple event handling supported
	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		fmt.Println(color.GreenString("  New On Connection "), color.CyanString("User: "), color.YellowString(ep.Kws.GetStringAttribute("user_id")))
		fmt.Println("")
	})

	// On message event
	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		socketUserId := ep.Kws.GetStringAttribute("user_id")
		fmt.Println(color.YellowString("  New On Message Event "), color.CyanString("User: "), color.YellowString(socketUserId))
		fmt.Println(string(ep.Data))
		fmt.Println("")

		models.EmitToSocketId("MatchOpponent", "Looking", "Looking", ep.Kws.UUID)
		time.Sleep(5 * time.Second)
		models.EmitToSocketId("MatchOpponent", "Found", "0x987fa83115c1212A6C7044636ff098a3C9e98Ed2", ep.Kws.UUID)
	})

	// On disconnect event
	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		fmt.Println(color.RedString("  New On Disconnect Event "), color.CyanString("User: "), color.YellowString(ep.Kws.GetStringAttribute("user_id")))
		fmt.Println("")
		delete(models.SocketClients, ep.Kws.GetStringAttribute("user_id"))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		fmt.Println(color.RedString("  New On Close Event "), color.CyanString("User: "), color.YellowString(ep.Kws.GetStringAttribute("user_id")))
		fmt.Println("")

		delete(models.SocketClients, ep.Kws.GetStringAttribute("user_id"))
	})

	// On error event
	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		fmt.Println(color.RedString("  New On Error Event "), color.CyanString("User: "), color.YellowString(ep.Kws.GetStringAttribute("user_id")))
		fmt.Println(color.CyanString("\tUser: "), color.YellowString(ep.Kws.GetStringAttribute("user_id")))
		fmt.Println("")
	})
}
