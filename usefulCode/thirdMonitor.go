package main

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${uri} ${remote_ip} ${host} ${referer} ${user_agent} ${status} ${latency_human}\n",
	}))
	e.GET("/ck", func(c echo.Context) error {
		encodedUrl := c.QueryParam("target")
		if encodedUrl == "" {
			return c.String(http.StatusBadRequest, "target url is empty.")
		}

		targetUrl, err := url.QueryUnescape(encodedUrl)
		if err != nil {
			return c.String(http.StatusBadRequest, "target url err:"+err.Error())
		}
		return c.Redirect(http.StatusFound, targetUrl)
	})
	e.Logger.Fatal(e.Start(":9032"))
}
