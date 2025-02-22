package maze

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"log/slog"
	"strconv"
)

type Payload struct {
	Width  string `json:"width"`
	Height string `json:"height"`
	Header string `json:"header"`
}

func WsGenerateMaze(conn *websocket.Conn) {
	payload := Payload{}
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = json.Unmarshal(message, &payload)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		width, err := strconv.Atoi(payload.Width)
		height, err := strconv.Atoi(payload.Height)

		maze := NewMaze(width, height)

		sig := make(chan Maze, 10)
		go func() {
			maze.Create(sig)
			maze.FindExit(sig)
			maze.FindExit(sig)
			maze.Clear(sig)
			close(sig)
		}()

		go func() {
			for {
				select {
				case mazeUpdate, open := <-sig:
					if !open {
						break
					}
					slog.Info("sending maze")
					buffer := new(bytes.Buffer)
					RenderMaze(&mazeUpdate, buffer)
					imageBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

					err := conn.WriteMessage(messageType,
						[]byte(
							fmt.Sprintf(`
					<div id="content">
					<img src="data:image/png;base64,%s"  alt="a maze"/>
					</div>`, imageBase64)))

					if err != nil {
						slog.Error(err.Error())
						return
					}
				}
			}
		}()
	}
}
