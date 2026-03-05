package gin

import (
	"context"
	"fmt"
	"log/slog"
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
	port := firstEnv("MINSTACK_HTTP_PORT", "MINSTACK_PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		Host: os.Getenv("MINSTACK_HOST"),
		Port: port,
	}
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func NewServer(lc fx.Lifecycle, log *slog.Logger) *gin.Engine {
	cfg := newConfig()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Info("route registered",
			"method", httpMethod,
			"path", absolutePath,
			"handler", handlerName,
			"handlers", nuHandlers,
		)
	}
	gin.DebugPrintFunc = func(format string, values ...interface{}) {
		log.Debug(strings.TrimRight(fmt.Sprintf(format, values...), "\n"))
	}

	r := gin.New()
	r.Use(requestLogger(log), recovery(log))

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
			log.Info("listening", "address", addr)
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
