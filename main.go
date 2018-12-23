package main

import (
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	repos map[string]*git.Repository
)

func main() {
	config, err := reloadConfig()
	if err != nil {
		log.Fatalf("Could not parse config file %v", err)
	}
	repos := make(map[string]*git.Repository)
	checkoutRepos(*config, repos)

	app := iris.Default()
	app.Use(recover.New())
	app.Use(logger.New())

	// this strange construct is used by the router to verify if the repo exist
	app.Macros().Get("string").RegisterFunc("repoExists", func() func(string) bool {
		return func(name string) bool {
			return repos[name] != nil
		}
	})

	app.Get("/list", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"repositories": config.Repositories})
	})

	app.Get("/{name:string repoExists()}/branches", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		branches, err := listRemoteRefs(repos, name, "heads")
		if err != nil {
			ctx.StatusCode(501)
			ctx.JSON(iris.Map{"status": "error", "description": err})
			return
		}

		ctx.JSON(iris.Map{"branches": branches})
	})

	app.Get("/{name:string repoExists()}/tags", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		branches, err := listRemoteRefs(repos, name, "tags")
		if err != nil {
			ctx.StatusCode(501)
			ctx.JSON(iris.Map{"status": "error", "description": err})
			return
		}

		ctx.JSON(iris.Map{"branches": branches})
	})

	app.Get("/fetch", func(ctx iris.Context) {
		details := fetchRepos(repos)
		ctx.JSON(iris.Map{"status": "OK", "details": details})
	})

	app.Get("/health", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"status": "OK"})
	})

	app.Get("/reload", func(ctx iris.Context) {
		config, err = reloadConfig()
		if err != nil {
			ctx.StatusCode(501)
			ctx.JSON(iris.Map{"status": "error", "description": err})
			return
		}
		ctx.JSON(iris.Map{"status": "OK"})
	})

	app.Run(iris.Addr(":8080"))
}
