package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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
		fmt.Println("Recived The Link Request")
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
    file, err := os.OpenFile("WebScrape.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write the header only if the file is newly created
    if stat, err := file.Stat(); err == nil && stat.Size() == 0 {
        header := []string{"Title", "Description", "Link To Article", "Publication Date", "Category", "Image URL"}
        if err := writer.Write(header); err != nil {
            log.Fatal(err)
        }
    }

    c := colly.NewCollector()

    c.OnXML("//item", func(e *colly.XMLElement) {
        fmt.Println("Found item") // Debugging output
        title := e.ChildText("title")
        description := e.ChildText("description")
        link := e.ChildText("link")
        pubDate := e.ChildText("pubDate")
        category := e.ChildText("category")
        
        // For media:content, handle the namespace correctly
        imageURL := e.ChildAttr("media:content", "url") // Use colon without escape if using OnXML

        // Create a row with the extracted data
        row := []string{title, description, link, pubDate, category, imageURL}
        // fmt.Println("Row data:", row) // Log the row data before writing
        
        // Attempt to write to CSV and check for errors
        if err := writer.Write(row); err != nil {
            log.Println("Error writing to CSV:", err)
        } else {
            fmt.Println("Successfully wrote row:", row) // Confirm successful write
        }
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Response received with status code:", r.StatusCode)
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting:", r.URL.String())
    })

    // Visit the RSS feed URL
    if err := c.Visit(url); err != nil {
        log.Println("Error visiting URL:", err)
    }
}

func HandleError(err error, c *fiber.Ctx) error {
	log.Println(err) // Log the error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
}
