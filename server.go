package gin

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Config struct {
	Host string
	Port string
}

func newConfig() Config {
	port := os.Getenv("MINSTACK_PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		Host: os.Getenv("MINSTACK_HOST"),
		Port: port,
	}
}

func NewServer(lc fx.Lifecycle) *gin.Engine {
	cfg := newConfig()
	r := gin.Default()

	if origin, ok := os.LookupEnv("MINSTACK_CORS_ORIGIN"); ok {
		corsConfig := cors.DefaultConfig()
		if origin == "*" {
			corsConfig.AllowOriginFunc = func(_ string) bool { return true }
		} else {
			corsConfig.AllowOrigins = strings.Split(origin, ",")
		}
		corsConfig.AddAllowHeaders("Authorization")
		r.Use(cors.New(corsConfig))
	}

	addr := cfg.Host + ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return r
}
