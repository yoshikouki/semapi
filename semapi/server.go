package semapi

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/yoshikouki/semapi/api"
	"github.com/yoshikouki/semapi/middleware"
	"github.com/yoshikouki/semapi/model"
)

func Launch(cfg Config) error {
	// get config
	conf, err := NewConfig(cfg)
	if err != nil {
		return err
	}

	// run Redis
	rdb, err := NewRedis(
		conf.RedisHost,
		conf.RedisPort,
		conf.RedisPassword,
		conf.RedisDB,
	)
	if err != nil {
		return err
	}

	model, err := model.NewModel(rdb)
	if err != nil {
		return err
	}

	return serverRun(conf, model)
}

// Run HTTP server
func serverRun(conf Config, model *model.Model) error {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(middleware.Model(model))
	e.Validator = middleware.NewCustomValidator()

	api.DefineEndpoints(e)

	port := fmt.Sprintf(":%d", conf.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err)
		return err
	}
	return nil
}
