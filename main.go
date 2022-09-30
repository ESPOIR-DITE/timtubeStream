package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"goStreaming/api"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2

var g errgroup.Group

func main() {

	fileServer := &http.Server{
		Addr:    ":8100",
		Handler: fsRouter(),
	}
	fmt.Println("Service Running..")
	err := fileServer.ListenAndServe()
	if err != nil {
		fmt.Println(err, " Service stop running..")
		log.Fatal(err)
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
func fsRouter() http.Handler {
	r := echo.New()
	r.GET("/vid/:id", func(c echo.Context) error {
		videoId := c.Param("id")
		f, err := os.Open("./files/" + videoId + ".mp4")
		if err != nil {
			fmt.Println(err, " error reading the following file: ", videoId)
			return err
		}
		return c.Stream(http.StatusOK, "video/mp4", f)
	})
	r.GET("/swagger/*", echoSwagger.WrapHandler)
	return r
}
func mainRouter() http.Handler {
	r := echo.New()
	r.GET("/vid", func(c echo.Context) error {
		videoId := c.QueryParam("videoId")
		req, _ := http.NewRequest("GET", api.BASE_URL+"video/video-data/stream/"+videoId, nil)
		client := &http.Client{}
		req.Header.Set("Range", c.Request().Header.Get("Range"))
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		c.Response().Header().Set(echo.HeaderContentLength, resp.Header.Get("Content-Length"))
		contentRange := resp.Header.Get("Content-Range")
		if contentRange != "" {
			c.Response().Header().Set("Content-Range", contentRange)
		}
		defer resp.Body.Close()
		return c.Stream(resp.StatusCode, "video/mp4", resp.Body)
	})

	return r
}
