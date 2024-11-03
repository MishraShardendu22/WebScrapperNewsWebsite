package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/gocolly/colly/v2" // Import Colly for web scraping
)

type LinkRequest struct {
	URL string `json:"url"` // Struct to hold the incoming URL
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Print a welcome message
	fmt.Println("This is My News Aggregator Web Scraper Project!!")

	// Create a new Fiber app
	app := fiber.New()

	// Define application routes
	TestRoutes(app)

	// Define the link route
	GetLinkRoute(app)

	// Set the port, with a default of 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start the server
	log.Fatal(app.Listen(":" + port))
}

func TestRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "This is a Test Route"})
	})
}

func GetLinkRoute(app *fiber.App) {
	app.Post("/getLink", func(c *fiber.Ctx) error {
		var linkReq LinkRequest
		err := c.BodyParser(&linkReq)
		if err != nil {
			return HandleError(err, c) // Pass the context to HandleError
		}
		
		// Call the WebScrapeRoute function with the provided URL
		WebScrapeRoute(linkReq.URL)

		return c.Status(200).JSON(fiber.Map{"message": "The Link was received", "url": linkReq.URL})
	})
}

func WebScrapeRoute(url string) {
	file, err := os.Create("WebScrape.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Title", "Description", "Link To Article"}
	if err := writer.Write(header); err != nil {
		log.Fatal(err) // Handle CSV write errors
	}

	c := colly.NewCollector()
	
	c.OnHTML("item", func(e *colly.HTMLElement) {
		title := e.ChildText("title")
		description := e.ChildText("description")
		link := e.ChildText("link")
		row := []string{title, description, link}
		if err := writer.Write(row); err != nil {
			log.Println("Error writing to CSV:", err)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if err := c.Visit(url); err != nil {
		log.Println("Error visiting URL:", err)
	}
}

func HandleError(err error, c *fiber.Ctx) error {
	log.Println(err) // Log the error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
}
