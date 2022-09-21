package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"goStreaming/api"
	"golang.org/x/sync/errgroup"
	"html/template"
	"log"
	"net/http"
	"os"
)

var g errgroup.Group

func main() {
	mainServer := &http.Server{
		Addr:    ":8000",
		Handler: mainRouter(),
	}
	fileServer := &http.Server{
		Addr:    ":8100",
		Handler: fsRouter(),
	}
	g.Go(func() error {
		return mainServer.ListenAndServe()
	})
	g.Go(func() error {
		return fileServer.ListenAndServe()
	})
	log.Println("Service Running..")
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
func fsRouter() http.Handler {
	r := echo.New()
	r.GET("/vid", func(c echo.Context) error {
		f, err := os.Open("./files/Easy.mp4")
		if err != nil {
			return err
		}
		//if f.Name() !=""{
		//	c.Response().Header().Set(echo.HeaderContentLength, "Content-Length")
		//}
		return c.Stream(http.StatusOK, "video/mp4", f)
	})
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

	r.GET("", func(c echo.Context) error {
		fileAdmins := []string{
			"index.html",
		}
		ts, err := template.ParseFiles(fileAdmins...)
		if err != nil {
			fmt.Println(" page passing fail")
			return err
		}
		type PageData struct {
			Name string
		}
		data := PageData{"Espoir"}
		w := c.Response().Writer
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err, "")
		}

		return c.File("home.html")
	})
	return r
}
