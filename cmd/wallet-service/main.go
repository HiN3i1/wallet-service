package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HiN3i1/wallet-service/db"
	"github.com/HiN3i1/wallet-service/service"
	"github.com/joho/godotenv"
)

var (
	method = flag.String("m", "", "method")
)

func main() {
	flag.Parse()
	godotenv.Load()
	switch *method {
	case "resetdb":
		db.CreateDBClient()
		db.CleanTable()
		db.InitTable()
		fmt.Println("reset db done.")

	case "run-server":
		runAllService()

	default:
		fmt.Println("unknown method:", *method)
	}
}

func runAllService() {
	log.Println("Starting DB Service...")
	service.CreateDBClient()

	log.Println("Starting HTTP API Service...")
	serverHTTP := service.NewAPIServer(":80")

	go func() {
		if err := serverHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e := fmt.Errorf("listen http failed. err=[%v]", err)
			fmt.Println(e)
			return
		}
	}()

	log.Println("Starting HTTPS API Service...")
	serverHTTPs := service.NewAPIServer(":443")

	go func() {
		if err := serverHTTPs.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != nil && err != http.ErrServerClosed {
			e := fmt.Errorf("listen https failed. err=[%v]", err)
			fmt.Println(e)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serverHTTP.Shutdown(ctx); err != nil {
		log.Println("Http Server Shutdown:", err)
	}

	if err := serverHTTPs.Shutdown(ctx); err != nil {
		log.Println("Https Server Shutdown:", err)
	}
}
