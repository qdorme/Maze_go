package maze

import (
	"fmt"
	"log/slog"
	"math/rand"
)

type Maze struct {
	Width    int
	Height   int
	Grid     [][]int
	Distance [][]uint
	Exit     []Cell
}

type Cell struct {
	X, Y int
}

func NewMaze(width, height int) *Maze {

	grid := make([][]int, height)
	dist := make([][]uint, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]int, width)
		dist[i] = make([]uint, width)
	}

	return &Maze{
		Width:    width,
		Height:   height,
		Grid:     grid,
		Distance: dist,
		Exit:     make([]Cell, 0),
	}
}

func (m *Maze) Create(sig chan Maze) {
	stack := Stack{
		Cells: make([]Cell, m.Width*m.Height),
		Index: -1,
	}
	// start with 0,0
	cell := stack.Push(Cell{0, 0})
	m.visiteCell(cell)
	for {
		err := m.chooseNeighbour(&stack)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		linkCells(m, &stack)
		if sig != nil {
			sig <- *m
		}
		// time.Sleep(10 * time.Millisecond)
	}

}

func (m *Maze) chooseNeighbour(stack *Stack) error {
popAgain:
	cell, err := stack.Pop()
	if err != nil {
		return err
	}
	neighbours := unvisitedNeighbours(m, cell)
	if len(neighbours) == 0 {
		goto popAgain
	}

	randomIndex := rand.Intn(len(neighbours))
	stack.Push(cell)
	stack.Push(neighbours[randomIndex])
	return nil
}

func unvisitedNeighbours(m *Maze, cell Cell) []Cell {
	cells := []Cell{}
	if cell.X > 0 && !m.IsVisited(Cell{cell.X - 1, cell.Y}) {
		cells = append(cells, Cell{cell.X - 1, cell.Y})
	}
	if cell.X < m.Width-1 && !m.IsVisited(Cell{cell.X + 1, cell.Y}) {
		cells = append(cells, Cell{cell.X + 1, cell.Y})
	}
	if cell.Y > 0 && !m.IsVisited(Cell{cell.X, cell.Y - 1}) {
		cells = append(cells, Cell{cell.X, cell.Y - 1})
	}
	if cell.Y < m.Height-1 && !m.IsVisited(Cell{cell.X, cell.Y + 1}) {
		cells = append(cells, Cell{cell.X, cell.Y + 1})
	}
	return cells
}

func (m *Maze) IsVisited(cell Cell) bool {
	return m.Grid[cell.Y][cell.X]&0b10000 == 0b10000
}

func (m *Maze) visiteCell(cell Cell) {
	m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b10000
}

func (m *Maze) FindExit(sig chan Maze) {
	cell := make([]Cell, 0)
	var value uint = 1
	if len(m.Exit) > 0 {
		cell = append(cell, m.Exit[0])
		dists := make([][]uint, m.Height)
		for i := 0; i < m.Height; i++ {
			dists[i] = make([]uint, m.Width)
		}
		m.Distance = dists
		m.Distance[m.Exit[0].Y][m.Exit[0].X] = value
	} else {
		cell = append(cell, Cell{})
		m.Distance[0][0] = value
	}

	for {
		value++
		cells := m.FindConnected(cell)
		if len(cells) == 0 {
			break
		}
		cell = cells
		for _, c := range cell {
			m.Distance[c.Y][c.X] = value
		}
		sig <- *m
	}

	exit := Cell{}
	for i := 0; i < m.Height; i++ {
		if m.Distance[i][0] > m.Distance[exit.Y][exit.X] {
			exit.X = 0
			exit.Y = i
		}
		if m.Distance[i][m.Width-1] > m.Distance[exit.Y][exit.X] {
			exit.X = m.Width - 1
			exit.Y = i
		}
	}
	for i := 0; i < m.Width; i++ {
		if m.Distance[0][i] > m.Distance[exit.Y][exit.X] {
			exit.X = i
			exit.Y = 0
		}
		if m.Distance[m.Height-1][i] > m.Distance[exit.Y][exit.X] {
			exit.X = i
			exit.Y = m.Height - 1
		}
	}
	m.Exit = append(m.Exit, exit)
	slog.Info(fmt.Sprintf("exit found at (%d,%d)", exit.X, exit.Y))

}

func (m *Maze) FindConnected(cells []Cell) []Cell {
	cellsOut := make([]Cell, 0)
	for _, cell := range cells {
		if m.Grid[cell.Y][cell.X]&0b0001 == 0b0001 && m.Distance[cell.Y-1][cell.X] == 0 {
			m.Distance[cell.Y-1][cell.X] = 1
			cellsOut = append(cellsOut, Cell{cell.X, cell.Y - 1})
		}
		if m.Grid[cell.Y][cell.X]&0b0010 == 0b0010 && m.Distance[cell.Y][cell.X+1] == 0 {
			m.Distance[cell.Y][cell.X+1] = 1
			cellsOut = append(cellsOut, Cell{cell.X + 1, cell.Y})
		}
		if m.Grid[cell.Y][cell.X]&0b0100 == 0b0100 && m.Distance[cell.Y+1][cell.X] == 0 {
			m.Distance[cell.Y+1][cell.X] = 1
			cellsOut = append(cellsOut, Cell{cell.X, cell.Y + 1})
		}
		if m.Grid[cell.Y][cell.X]&0b1000 == 0b1000 && m.Distance[cell.Y][cell.X-1] == 0 {
			m.Distance[cell.Y][cell.X-1] = 1
			cellsOut = append(cellsOut, Cell{cell.X - 1, cell.Y})
		}
	}
	return cellsOut
}

func (m *Maze) Clear(sig chan Maze) {
	m.Distance = make([][]uint, m.Height)
	for i := 0; i < m.Height; i++ {
		m.Distance[i] = make([]uint, m.Width)
	}

	for _, cell := range m.Exit {
		if cell.X == 0 {
			m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b1000
		} else if cell.X == m.Width-1 {
			m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b0010
		} else if cell.Y == 0 {
			m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b0001
		} else {
			m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b0100
		}
	}
	sig <- *m

}

func linkCells(m *Maze, stack *Stack) {
	newCell, cell, err := stack.Lasts()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	// slog.Info(fmt.Sprintf("linking cell (%d,%d) with cell (%d,%d)", cell.X, cell.Y, newCell.X, newCell.Y))
	m.visiteCell(newCell)
	if cell.X == newCell.X {
		if cell.Y == newCell.Y-1 {
			m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b0100
			m.Grid[newCell.Y][newCell.X] = m.Grid[newCell.Y][newCell.X] | 0b0001
		} else {
			m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b0001
			m.Grid[newCell.Y][newCell.X] = m.Grid[newCell.Y][newCell.X] | 0b0100
		}
	} else if cell.X == newCell.X-1 {
		m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b0010
		m.Grid[newCell.Y][newCell.X] = m.Grid[newCell.Y][newCell.X] | 0b1000
	} else {
		m.Grid[cell.Y][cell.X] = m.Grid[cell.Y][cell.X] | 0b1000
		m.Grid[newCell.Y][newCell.X] = m.Grid[newCell.Y][newCell.X] | 0b0010
	}
}
