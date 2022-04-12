package main

import (
	"bridgeApi/app"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"github.com/fatih/color"
	"github.com/panjf2000/ants"
)

// var (
// 	defaultListenAddress = "127.0.0.1:9000"
// 	defaultProxyUrl      = "https://apis.ankr.com/9ba1afc9df64426f9123fd6d9f8ac9a5/8e1a6b2b88490d4e20818137e607c759/avax/archive/main"
// 	defaultRelayUrl      = "https://relay.flashbots.net"
// 	defaultRedisUrl      = "localhost:6379"

// 	version = "dev" // is set during build process
// )

// var versionPtr = flag.Bool("version", false, "just print the program version")
// var listenAddress = flag.String("listen", getEnvOrDefault("LISTEN_ADDR", defaultListenAddress), "Listen address")
// var proxyUrl = flag.String("proxy", getEnvOrDefault("PROXY_URL", defaultProxyUrl), "URL for default JSON-RPC proxy target (eth node, Infura, etc.)")
// var redisUrl = flag.String("redis", getEnvOrDefault("REDIS_URL", defaultRedisUrl), "URL for Redis (use 'dev' to use integrated in-memory redis)")

// // Flags for using the relay
// var relayUrl = flag.String("relayUrl", getEnvOrDefault("RELAY_URL", defaultRelayUrl), "URL for relay")
// var relaySigningKey = "a2233805bfbbe29d313ac95e67ed2b13bf9c798afc92e88e664b815f8eb57afa"

func main() {

	defer ants.Release()
	a := new(app.App)
	var wg sync.WaitGroup

	// system config
	numCpu := runtime.NumCPU()
	usedCpu := numCpu
	runtime.GOMAXPROCS(usedCpu)
	fmt.Println("")
	fmt.Println(color.YellowString("  ----------------- System Info -----------------"))
	fmt.Println(color.CyanString("\t    Number CPU cores available: "), color.GreenString(strconv.Itoa(numCpu)))
	fmt.Println(color.MagentaString("\t    Used of CPU cores: "), color.YellowString(strconv.Itoa(usedCpu)))
	fmt.Println(color.MagentaString(""))

	// Tasks
	var fiberHttpTask = func() {
		a.Initialize(":3001")
		wg.Done()
	}

	var web3Task = func() {
		a.InitializeWeb3()
		wg.Done()
	}

	// var rpcServerTask = func() {
	// 	var key *ecdsa.PrivateKey
	// 	var err error

	// 	flag.Parse()

	// 	// Perhaps print only the version
	// 	if *versionPtr {
	// 		fmt.Printf("rpc-endpoint %s\n", version)
	// 		return
	// 	}

	// 	log.Printf("rpc-endpoint %s\n", version)

	// 	if relaySigningKey == "" {
	// 		log.Fatal("Cannot use the relay without a signing key.")
	// 	}

	// 	pkHex := strings.Replace(relaySigningKey, "0x", "", 1)
	// 	if pkHex == "dev" {
	// 		log.Println("Creating a new dev signing key...")
	// 		key, err = crypto.GenerateKey()
	// 	} else {
	// 		key, err = crypto.HexToECDSA(pkHex)
	// 	}

	// 	if err != nil {
	// 		log.Fatal("Error with relay signing key:", err)
	// 	}

	// 	log.Printf("Signing key: %s\n", crypto.PubkeyToAddress(key.PublicKey).Hex())

	// 	// Start the endpoint
	// 	s := server.NewRpcEndPointServer(version, *listenAddress, *proxyUrl, *relayUrl, key)
	// 	s.Start()
	// 	wg.Done()
	// }

	// add task to ants pool
	wg.Add(1)
	_ = ants.Submit(fiberHttpTask)

	wg.Add(1)
	_ = ants.Submit(web3Task)

	// wg.Add(1)
	// _ = ants.Submit(rpcServerTask)

	// wait all tasks
	wg.Wait()
	//fmt.Printf("running goroutines: %d\n", ants.Running())
	//fmt.Printf("finish all tasks.\n")
}

// func getEnvOrDefault(key string, defaultValue string) string {
// 	ret := os.Getenv(key)
// 	if ret == "" {
// 		ret = defaultValue
// 	}
// 	return ret
// }
