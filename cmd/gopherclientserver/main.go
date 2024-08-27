package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	fun "github.com/AlexGrtsnk/go-prof-diploma-2/internal/functions"
	grp "github.com/AlexGrtsnk/go-prof-diploma-2/internal/grpc"
)

// глобальные переменные флагов
var (
	// версия билда
	BuildVersion string
	// время и дата билда
	BuildDate string
	// коммит билда
	BuildCommit string
)

func main() {
	fmt.Printf("version=%s, date=%s, commit=%ss\n", BuildVersion, BuildDate, BuildCommit)
	isUsingGRPC := false
	srv, enableHTTPS := fun.Run()
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if isUsingGRPC {
		grp.RunGRPCServer()
		return
	}
	if !enableHTTPS {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		if err := srv.ListenAndServeTLS("certificate", "key"); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)

		}
	}
	<-idleConnsClosed

}
