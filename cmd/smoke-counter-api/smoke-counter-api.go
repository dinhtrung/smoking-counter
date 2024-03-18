package main

import (
	"embed"
	"flag"
	app "github.com/dinhtrung/smoking-counter/internal"
	authImpl "github.com/dinhtrung/smoking-counter/internal/app/api-gateway/services/impl"
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/services/impl"
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/web/rest"
	authJwt "github.com/dinhtrung/smoking-counter/pkg/fiber/authjwt"
	authApi "github.com/dinhtrung/smoking-counter/pkg/fiber/authjwt/web/rest"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"net/http"
	"strings"
)

var configFile string

// Embed the entire web directory
//
//go:embed web/*
var embedDirStatic embed.FS

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&configFile, "config", "configs/api-gateway.yaml", "API Gateway configuration file")
	flag.Parse()

	// 1 - set default settings for components.
	app.BuntDBConfig()

	// 2 - override defaults with configuration file and watch changes
	app.ConfigInit(configFile)
	// app.ConfigWatch(configFile)

	app.BuntDBInit()

	// 3 - bring up components
	// + inject UserServiceDummy into application
	userRepo := authImpl.NewUserRepositoryBuntDB(app.Config.MustString("buntdb.path"))
	userSvc := authImpl.NewUserServiceBuntDB(userRepo)
	// 4 - setup the web server
	srv := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
		Prefork:   false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				return c.Status(code).JSON(e)
			}
			return c.Status(code).JSON(fiber.Map{"error": code, "message": err.Error()})
		},
	})
	configureFiber(srv)

	authJwt.USER_RESOURCE = authApi.NewDefaultUserResource(userSvc, userRepo)
	authJwt.ACCOUNT_RESOURCE = authApi.NewDefaultAccountResource(userSvc)
	authJwt.SetupAuthJWT(srv, app.Config.MustString("security.jwt-secret"), app.Config.Strings("security.skip-auth")...)
	authJwt.SetupRoutes(srv)

	// main application endpoints
	smokeAPI := rest.NewSmokeRestAPI(impl.NewSmokeServiceBuntDB(app.BuntDB))
	srv.Get("/api/smokes", smokeAPI.GetAll)
	srv.Post("/api/smokes", smokeAPI.Create)
	srv.Delete("/api/smokes", smokeAPI.Delete)

	// bring the server up
	if app.Config.String("https.listen") != "" {
		log.Fatal(srv.ListenTLS(app.Config.String("https.listen"), app.Config.MustString("https.cert"), app.Config.MustString("https.key")))
	} else {
		log.Fatal(srv.Listen(app.Config.String("http.listen")))
	}
}

// configureFiber start the fiber with common settings
func configureFiber(srv *fiber.App) {
	staticAsset := filesystem.New(filesystem.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api") || strings.HasPrefix(c.Path(), "/services")
		},
		PathPrefix: "web",
		Root:       http.FS(embedDirStatic),
	})
	srv.Use(staticAsset)
	// + logging
	srv.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
	}))

	srv.Use(recover.New())
	// management endpoints
	srv.Get("/management/info", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusOK, "ok")
	})
	srv.Get("/management/health", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusOK, "ok")
	})
}
