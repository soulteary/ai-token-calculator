package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed public
var staticFS embed.FS
var Version string

// Credit: https://github.com/gin-contrib/static/issues/19
type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(efs),
	}
}

func main() {
	var debugMode bool
	if strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG"))) == "on" {
		debugMode = true
	} else {
		debugMode = false
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(gzip.Gzip(gzip.DefaultCompression))

	if debugMode {
		route.Use(static.Serve("/", static.LocalFile("./public", false)))
	} else {
		route.Use(static.Serve("/", EmbedFolder(staticFS, "public")))
	}

	route.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	srv := &http.Server{
		Addr:              "0.0.0.0:" + strconv.Itoa(GetPort()),
		Handler:           route,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Program start error: %s\n", err)
		}
	}()
	log.Println("github.com/soulteary/ai-token-calculator has started 🚀")
	if Version != "" {
		log.Printf("Version: %s\n", Version)
	}

	<-ctx.Done()

	stop()
	log.Println("The program is closing, if you want to end it immediately, please press `CTRL+C`")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Program was forced to close: %s\n", err)
	}
}

func GetPort() int {
	defaultPort := 8080
	portStr := os.Getenv("PORT")

	if portStr == "" {
		log.Printf("The PORT environment variable is empty, using the default port: %d\n", defaultPort)
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("The PORT environment variable is not a valid integer, using the default port: %d\n", defaultPort)
		return defaultPort
	}

	if port < 1 || port > 65535 {
		log.Printf("The PORT environment variable is not a valid port number, using the default port: %d\n", defaultPort)
		return defaultPort
	}

	log.Printf("Using the port: %d\n", port)
	return port
}
