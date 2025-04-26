package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/util"
)

func verify(c *fiber.Ctx, conf *config.Config, header string) bool {
	hasher := hmac.New(sha256.New, []byte(conf.WebhookSecret))
	hasher.Write(c.Body())
	return hex.EncodeToString(hasher.Sum(nil)) == c.Get(header)
}

func handle(c *fiber.Ctx, event string) {
	reload_chan := c.Locals("reload").(chan bool)

	switch event {
	case "push":
		reload_chan <- true
	}
}

func POST_Webhook(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Config)
	platform := c.Params("platform")

	if conf.WebhookSecret == "" {
		return util.NotFound(c)
	}

	switch strings.ToLower(platform) {
	case "github":
		if !verify(c, conf, "x-hub-signature-256") {
			return util.BadRequest(c)
		}

		handle(c, c.Get("x-github-event"))
		return c.Status(202).SendString("Accepted")

	case "gitea":
		if !verify(c, conf, "HTTP_X_GITEA_SIGNATURE") {
			return util.BadRequest(c)
		}

		handle(c, c.Get("X-Gitea-Event"))
		return c.Status(202).SendString("Accepted")

	case "forgejo":
		if !verify(c, conf, "HTTP_X_FORGEJO_SIGNATURE") {
			return util.BadRequest(c)
		}

		handle(c, c.Get("X-Forgejo-Event"))
		return c.Status(202).SendString("Accepted")
	}

	return util.BadRequest(c)
}
