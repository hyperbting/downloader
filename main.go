package main

import (
	"downloader/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/mustache/v2"
	"net/http"
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

	app.Get("/download2", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title":      "Hello, World!",
			"SubmitPath": "./download2",
		})
	})

	//app.Post("/download", func(c *fiber.Ctx) error {
	//	p := new(models.DownloadTarget)
	//	if err := c.BodyParser(p); err != nil {
	//		return err
	//	}
	//
	//	//log.Println(p.Group)
	//	//log.Println(p.Number)
	//	//log.Println(p.Name)
	//	log.Info(p)
	//
	//	p.Sanitize()
	//
	//	// Create the command to execute the script
	//	cmd := exec.Command("bash", "/ref/download.sh", p.Group, p.Number, p.Name)
	//
	//	// Set the working directory
	//	cmd.Dir = "/ref"
	//	log.Info(cmd.String())
	//	// Capture the output
	//	output, err := cmd.CombinedOutput()
	//	if err != nil {
	//		return c.JSON(fiber.Map{"error": err.Error()})
	//	}
	//
	//	return c.JSON(fiber.Map{"result": "ok", "output": string(output)})
	//})

	app.Post("/download", func(c *fiber.Ctx) error {
		p := new(models.DownloadTarget)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		localBasePath := "/ref"
		p.SetLocalPathBase(localBasePath)

		p.Sanitize()
		log.Info(p)

		// try download main image
		if err := p.TryDownloadMain(); err != nil {
			log.Infof("TryDownloadMain() err:%v\n", err)
			return c.SendStatus(http.StatusBadRequest)
		}

		// download sub images; do not care if err occurred
		if err := p.DownloadSub(); err != nil {
			log.Infof("DownloadSub() err:%v\n", err)
			//return c.SendStatus(http.StatusBadRequest)
		}

		if !p.HadFilesDownloaded() {
			return c.JSON(fiber.Map{"result": "Empty", "output": string("NoFileDownloaded")})
		}

		// move to desired folder
		if err := p.MoveLocalFilesUnderFolder(); err != nil {
			log.Infof("MoveLocalFilesUnderFolder() err:%v\n", err)
			return c.SendStatus(http.StatusBadRequest)
		}

		return c.JSON(fiber.Map{"result": "ok", "output": string("")})
	})

	app.Listen(":3000")
}
