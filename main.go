package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/logs"
)

type Payload struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// Validate  ...
func (p *Payload) Validate() error {
	if p.Method == "" {
		return fmt.Errorf("method is required")
	}
	if p.Url == "" {
		return fmt.Errorf("url is required")
	}
	if p.Headers == nil {
		p.Headers = make(map[string]string)
	}
	if p.Body == nil {
		p.Body = []byte{}
	}
	return nil
}

func main() {
	app := fiber.New()

	app.Post("/proxy", func(c *fiber.Ctx) error {
		var payload Payload
		if err := c.BodyParser(&payload); err != nil {
			logs.Error("parsing request, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if err := payload.Validate(); err != nil {
			logs.Error("validating request, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		res, err := doRequest(payload.Method, payload.Url, payload.Headers, payload.Body)
		if err != nil {
			logs.Error("sending request, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(res)
	})

	if err := app.Listen(":8080"); err != nil {
		logs.Error("could not proxy-man's server")
	}
}

func doRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	switch method {
	case "GET":
		return handleGetRequest(url, headers)
	case "POST", "PUT", "DELETE":
		return handlePostRequest(method, url, headers, body)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

func handleGetRequest(url string, headers map[string]string) ([]byte, error) {
	res, err := gttp.Get(url, headers)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("error: %d, msg: %s", res.StatusCode, string(res.Body))
	}

	return res.Body, nil
}

func handlePostRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	res, err := gttp.SendRequest(url, method, headers, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("error: %d, msg: %s", res.StatusCode, string(res.Body))
	}

	return res.Body, nil
}
