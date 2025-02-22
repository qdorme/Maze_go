package maze

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
)

const PIXEL_SIZE = 15

func RenderMaze(Maze *Maze, out io.Writer) {

	img := image.NewRGBA(image.Rect(-PIXEL_SIZE, -PIXEL_SIZE, Maze.Width*PIXEL_SIZE+1+PIXEL_SIZE, Maze.Height*PIXEL_SIZE+1+PIXEL_SIZE))

	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	maxCells := Maze.Width * Maze.Height

	for h := 0; h < Maze.Height; h++ {
		for w := 0; w < Maze.Width; w++ {

			if Maze.Distance[h][w] > 0 {
				alpha := uint8(255 * (1 - float64(Maze.Distance[h][w])/float64(maxCells)))
				drawSquare(img, w, h, alpha)
			}

			if Maze.Grid[h][w]&0b0001 == 0 {
				drawLine(img, w, h)
			}
			if Maze.Grid[h][w]&0b1000 == 0 {
				drawColumn(img, w, h)
			}
			if Maze.Grid[h][w]&0b0100 == 0 {
				drawLine(img, w, h+1)
			}
			if Maze.Grid[h][w]&0b0010 == 0 {
				drawColumn(img, w+1, h)
			}
		}
	}

	// Encodage en PNG
	if err := png.Encode(out, img); err != nil {
		panic(err)
	}
}

func drawSquare(img *image.RGBA, w, h int, alpha uint8) {
	W := w * PIXEL_SIZE
	H := h * PIXEL_SIZE
	purple := color.RGBA{alpha, 0, alpha, 255}
	draw.Draw(img, image.Rect(W, H, W+PIXEL_SIZE, H+PIXEL_SIZE), &image.Uniform{purple}, image.Point{}, draw.Src)

}

func drawLine(img *image.RGBA, w, h int) {
	W := w * PIXEL_SIZE
	for i := W; i <= W+PIXEL_SIZE; i++ {
		img.Set(i, h*PIXEL_SIZE, color.Black)
	}
}

func drawColumn(img *image.RGBA, w, h int) {
	H := h * PIXEL_SIZE
	for i := H; i <= H+PIXEL_SIZE; i++ {
		img.Set(w*PIXEL_SIZE, i, color.Black)
	}
}
