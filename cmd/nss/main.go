package main

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	// Hosts
	hosts := map[string]*Host{}

	// Ping subdomain
	ping := echo.New()
	ping.Use(middleware.Logger())
	ping.Use(middleware.Recover())
	hosts["ping.localhost:1323"] = &Host{ping}
	// Routes
	ping.GET("/", hello)

	// Proxy subdomain
	proxy := echo.New()
	proxy.Use(middleware.Logger())
	proxy.Use(middleware.Recover())
	hosts["proxy.localhost:1323"] = &Host{proxy}
	// Proxying localhost:26657 to proxy.localhost:1323
	url, err := url.Parse("http://localhost:26657")
	if err != nil {
		proxy.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url,
		},
	}
	proxy.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	// Website
	site := echo.New()
	site.Use(middleware.Logger())
	site.Use(middleware.Recover())

	hosts["localhost:1323"] = &Host{site}

	site.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Website")
	})

	// Start server
	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
