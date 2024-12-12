package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ImportDashboardHandler(c *fiber.Ctx) error {

	userRequest := c.Params("url")

	fmt.Println(userRequest)

	if userRequest == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	redirectUrl := "http://localhost:3000/d/ae6mbx951cfswf123/url-analytics?orgId=1&from=now-30d&to=now&timezone=browser&var-short_url=" + userRequest

	return c.Redirect(redirectUrl, http.StatusFound)
}
