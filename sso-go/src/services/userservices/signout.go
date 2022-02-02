package userservices

import (
	"time"

	"github.com/amitdotkr/sso/sso-go/src/global"
	"github.com/gofiber/fiber/v2"
)

func SignOut(c *fiber.Ctx) error {
	baseUrl := c.BaseURL()
	rc := new(fiber.Cookie)
	rc.Name = "refresh_token"
	// rc.Value = rt
	rc.Expires = time.Now().Add(-(time.Hour * 2))
	rc.HTTPOnly = true
	rc.Secure = global.CookieSecure
	rc.SameSite = "Strict"

	ac := new(fiber.Cookie)
	ac.Name = "access_token"
	// ac.Value = at
	ac.Expires = time.Now().Add(-(time.Hour * 2))
	ac.HTTPOnly = true
	ac.Secure = global.CookieSecure
	ac.SameSite = "Strict"
	c.Cookie(ac)
	c.Cookie(rc)
	// c.ClearCookie()
	// c.Location(baseUrl)
	return c.Status(fiber.StatusOK).Redirect(baseUrl)
	// return c.Redirect(baseUrl, 200)
	// return c.SendStatus(fiber.StatusOK)
}
