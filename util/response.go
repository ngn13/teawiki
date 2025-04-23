package util

import (
	"net/http"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/consts"
)

func render(c *fiber.Ctx, code int, view string, _data ...fiber.Map) error {
	var data fiber.Map = nil

	if len(_data) > 0 {
		data = _data[0]
	} else {
		data = fiber.Map{}
	}

	if code != 200 {
		data["dir"] = "/"
	}

	data["source"] = consts.SOURCE
	data["version"] = consts.VERSION

	data[view+"_view"] = true
	data["code"] = code

	data["repo"] = c.Locals("repo")
	data["conf"] = c.Locals("config")

	return c.Status(code).Render(view, data)
}

func Ok(c *fiber.Ctx, view string, data fiber.Map) error {
	if data["path"] != nil {
		data["dir"] = path.Dir(data["path"].(string))
	} else {
		data["dir"] = "/"
	}

	return render(c, http.StatusOK, view, data)
}

func NotFound(c *fiber.Ctx, data ...fiber.Map) error {
	return render(c, http.StatusNotFound, "error", data...)
}

func BadRequest(c *fiber.Ctx, data ...fiber.Map) error {
	return render(c, http.StatusBadRequest, "error", data...)
}

func Send(c *fiber.Ctx, file string) error {
	conf := c.Locals("config").(*config.Config)
	fp := path.Join(conf.RepoPath, file)

	if Exists(fp) {
		return c.SendFile(fp)
	}

	return NotFound(c, fiber.Map{
		"path": file,
	})
}
