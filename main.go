package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

type Request struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    any               `json:"body"`
	payload []byte
}

// Validate  ...
func (p *Request) Validate() error {
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
		p.payload = []byte{}
	}

	if p.Body != nil {
		// if body is already a byte array, no need to convert
		if _, ok := p.Body.([]byte); ok {
			p.payload = p.Body.([]byte)
		}

		// is body is a string, convert to byte array
		if _, ok := p.Body.(string); ok {
			p.payload = []byte(p.Body.(string))
		}

		p.payload = helpers.ToBytes(p.Body)
	}

	return nil
}

func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "${time}  ${method} ${path} ${status} ${latency}\n",
		TimeFormat: "2006/01/02 15:04:05",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/proxy", func(c *fiber.Ctx) error {
		logs.Info("request received: %s", string(c.BodyRaw()))

		var req Request
		if err := c.BodyParser(&req); err != nil {
			logs.Error("parsing request, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if err := req.Validate(); err != nil {
			logs.Error("validating request, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		res, err := doRequest(req)
		if err != nil {
			logs.Error("sending request, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		var resp any
		if err := json.Unmarshal(res, &resp); err != nil {
			logs.Error("unmarshalling response, %s", err.Error())
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		logs.Info("response: %s", string(helpers.ToBytes(resp)))
		return c.JSON(resp)
	})

	if err := app.Listen(":8080"); err != nil {
		logs.Error("could not proxy-man's server")
	}
}

func doRequest(req Request) ([]byte, error) {
	switch req.Method {
	case "GET":
		return handleGetRequest(req.Url, req.Headers)
	case "POST", "PUT", "DELETE":
		return handlePostRequest(req.Method, req.Url, req.Headers, req.payload)
	default:
		return nil, fmt.Errorf("unsupported method: %s", req.Method)
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
