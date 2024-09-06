package handlers

import (
	"log"
	"net/http"

	"github.com/SubhamMurarka/Schotky/Config"
	models "github.com/SubhamMurarka/Schotky/Models"
	services "github.com/SubhamMurarka/Schotky/Services"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	s services.ShortenServices
	r services.ResolveServices
}

func NewHandler(ser services.ShortenServices, res services.ResolveServices) *Handler {
	return &Handler{
		s: ser,
		r: res,
	}
}

func (h *Handler) ShortenUrl(c *fiber.Ctx) error {
	var body models.Request

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot Parse Request"})
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Url"})
	}

	url, err := h.s.ShortUrl(body.URL)
	if err != nil {
		log.Println("unable to create shorturl", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "url not shoertened"})
	}

	resp := models.Response{
		URL:    url,
		Expiry: Config.Cfg.ExpiryTime,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *Handler) ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")
	longUrl, err := h.r.ResolveURL(url)
	if err != nil {
		log.Println("not able to fetch url ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "No url Found"})
	}

	return c.Redirect(longUrl, 301)
}
