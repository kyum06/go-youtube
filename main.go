package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	var yt YOUTUBE

	e.GET("/", func(c echo.Context) error {

		URL := c.QueryParam("url")
		url, _ := yt.Get(URL)

		res, _ := http.Get(url)
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)
		return c.Blob(200, "audio/webm", body)

	})

	e.Start(":80")
}
