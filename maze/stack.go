package maze

import "errors"

type Stack struct {
	Cells []Cell
	Index int
}

func (s *Stack) Push(cell Cell) Cell {
	s.Index++
	s.Cells[s.Index] = cell
	return cell
}

func (s *Stack) Pop() (Cell, error) {
	if s.Index < 0 {
		return Cell{}, errors.New("Stack is empty")
	}
	cell := s.Cells[s.Index]
	s.Index--
	return cell, nil
}

func (s Stack) Lasts() (Cell, Cell, error) {
	if s.Index < 0 {
		return Cell{}, Cell{}, errors.New("Stack is empty")
	} else if s.Index == 0 {
		return s.Cells[0], Cell{}, nil
	} else {
		return s.Cells[s.Index], s.Cells[s.Index-1], nil
	}
}
