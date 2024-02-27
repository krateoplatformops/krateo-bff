package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/krateoplatformops/krateo-bff/internal/env"
	"github.com/krateoplatformops/krateo-bff/internal/server/middlewares/cors"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/actions"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/health"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/layout/columns"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/layout/rows"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/verbs"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/widgets/cardtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/server/routes/widgets/formtemplates"
	"github.com/rs/zerolog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	serviceName = "krateo-bff"
)

var (
	Version string
	Build   string
)

func main() {
	// Flags
	kconfig := flag.String(clientcmd.RecommendedConfigPathFlag, "", "absolute path to the kubeconfig file")
	debugOn := flag.Bool("debug", env.Bool("KRATEO_BFF_DEBUG", false), "dump verbose output")
	dumpEnv := flag.Bool("dumpEnv", env.Bool("KRATEO_BFF_DUMP_ENV", false), "dump environment variables")
	corsOn := flag.Bool("cors", env.Bool("KRATEO_BFF_CORS", true), "enable or disable CORS")
	port := flag.Int("port", env.ServicePort("KRATEO_BFF_PORT", 8080), "port to listen on")
	authnNS := flag.String("authn-store-namespace",
		env.String("AUTHN_STORE_NAMESPACE", ""),
		"krateo authn service clientconfig secrets namespace")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Initialize the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level for this log is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugOn {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log := zerolog.New(os.Stdout).With().
		Str("service", serviceName).
		Timestamp().
		Logger()

	if log.Debug().Enabled() {
		evt := log.Debug().
			Str("version", Version).
			Str("build", Build).
			Str("debug", fmt.Sprintf("%t", *debugOn)).
			Str("cors", fmt.Sprintf("%t", *corsOn)).
			Str("port", fmt.Sprintf("%d", *port)).
			Str("authn-store-namespace", *authnNS)
		if *dumpEnv {
			evt = evt.Strs("env-vars", os.Environ())
		}

		evt.Msg("configuration and env vars info")
	}

	var cfg *rest.Config
	var err error
	if len(*kconfig) > 0 {
		cfg, err = clientcmd.BuildConfigFromFlags("", *kconfig)
	} else {
		cfg, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Fatal().Err(err).Msg("resolving kubeconfig for rest client")
	}

	healthy := int32(0)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(routes.Logger(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	if *corsOn {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Auth-Code"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	health.Register(r, health.Options{
		Healty: &healthy, Version: Version, Build: Build, ServiceName: serviceName,
	})
	cardtemplates.Register(r, cfg, *authnNS)
	formtemplates.Register(r, cfg, *authnNS)
	columns.Register(r, cfg, *authnNS)
	rows.Register(r, cfg, *authnNS)
	verbs.Register(r, cfg)
	actions.Register(r, cfg, *authnNS)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 50 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}...)
	defer stop()

	go func() {
		atomic.StoreInt32(&healthy, 1)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("could not listen on %s", server.Addr)
		}
	}()

	// Listen for the interrupt signal.
	log.Info().Msgf("server is ready to handle requests at @ %s", server.Addr)

	if *debugOn {
		chi.Walk(r, func(method string, route string, handler http.Handler, _ ...func(http.Handler) http.Handler) error {
			log.Debug().Msgf("%s %s", method, route)
			return nil
		})
	}

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info().Msg("server is shutting down gracefully, press Ctrl+C again to force")
	atomic.StoreInt32(&healthy, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server gracefully stopped")
}
