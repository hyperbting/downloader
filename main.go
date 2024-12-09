package main

import (
	"downloader/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	"log"
	"os/exec"
)

func main() {

	// Create a new engine
	engine := mustache.New("./views", ".mustache")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files from the "js" directory
	app.Static("/js", "./views/js")

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.SendString("ok")
	})

	app.Get("/download", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title":      "Hello, World!",
			"SubmitPath": "./download",
		})
	})

	app.Post("/download", func(c *fiber.Ctx) error {
		p := new(models.DownloadTarget)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		//log.Println(p.Group)
		//log.Println(p.Number)
		//log.Println(p.Name)
		log.Println(p)

		p.Sanitize()

		// Create the command to execute the script
		cmd := exec.Command("bash", "/ref/download.sh", p.Group, p.Number, p.Name)

		// Set the working directory
		cmd.Dir = "/ref"
		log.Println(cmd.String())
		// Capture the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			return c.JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"result": "ok", "output": string(output)})
	})

	app.Post("/download2", func(c *fiber.Ctx) error {
		p := new(models.DownloadTarget)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		p.Sanitize()

		//if err != nil {
		//	return c.JSON(fiber.Map{"error": err.Error()})
		//}

		return c.JSON(fiber.Map{"result": "ok", "output": string("")})
	})

	app.Listen(":3000")
}
