package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"maze/maze"
	"os"
)

func main() {

	app := fiber.New(fiber.Config{
		AppName: "Maze Generator",
	})

	app.Static("/", "./public")

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/generate", websocket.New(maze.WsGenerateMaze))

	if err := app.Listen("localhost:3000"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
