package snake

import (
	"math/rand"
	"time"
)

// Action represents a turn (or no turn) command for the snake.
type Action int

const (
	DoNothing Action = iota
	RotateLeft
	RotateRight
)

// GameResult indicates the current state of the game.
type GameResult int

const (
	Running GameResult = iota
	GameOver
	Winner
)

// Coordinates defines a cell position on the game board.
type Coordinates struct {
	Row int
	Col int
}

// Snake represents the snake, including its body and moving direction.
type Snake struct {
	Body      []Coordinates // Head is at index 0.
	Direction Coordinates   // Direction vector (row, col).
}

// NewSnake creates a new snake starting at (startRow, startCol) moving right.
func NewSnake(startRow, startCol int) *Snake {
	return &Snake{
		Body:      []Coordinates{{startRow, startCol}},
		Direction: Coordinates{0, 1},
	}
}

// Head returns the current head coordinate.
func (s *Snake) Head() Coordinates {
	return s.Body[0]
}

// Move advances the snake in its current direction. If grow is true, the tail is not removed.
func (s *Snake) Move(grow bool) {
	newHead := Coordinates{
		Row: s.Head().Row + s.Direction.Row,
		Col: s.Head().Col + s.Direction.Col,
	}
	s.Body = append([]Coordinates{newHead}, s.Body...)
	if !grow {
		s.Body = s.Body[:len(s.Body)-1]
	}
}

// ChangeDirection rotates the snake's direction based on the action.
func (s *Snake) ChangeDirection(action Action) {
	if action == RotateLeft {
		s.Direction = Coordinates{-s.Direction.Col, s.Direction.Row}
	} else if action == RotateRight {
		s.Direction = Coordinates{s.Direction.Col, -s.Direction.Row}
	}
}

// Contains returns true if the snake's body occupies the given coordinate.
func (s *Snake) Contains(coord Coordinates) bool {
	for _, part := range s.Body {
		if part == coord {
			return true
		}
	}
	return false
}

// Game holds the overall game state.
type Game struct {
	Width     int
	Height    int
	Snake     *Snake
	Food      Coordinates
	GameState GameResult
}

// NewGame initializes a new game with a snake and random food.
func NewGame(width, height int) *Game {
	snake := NewSnake(height/2, width/2)
	return &Game{
		Width:     width,
		Height:    height,
		Snake:     snake,
		Food:      GenerateFood(width, height, snake),
		GameState: Running,
	}
}

// GenerateFood selects a random board cell that is not occupied by the snake.
func GenerateFood(width, height int, snake *Snake) Coordinates {
	rand.Seed(time.Now().UnixNano())
	for {
		food := Coordinates{
			Row: rand.Intn(height),
			Col: rand.Intn(width),
		}
		if !snake.Contains(food) {
			return food
		}
	}
}

// Update advances the game state based on the snake's movement and an action.
func (g *Game) Update(action Action) {
	if g.GameState != Running {
		return
	}

	// Change direction based on user action.
	g.Snake.ChangeDirection(action)

	// Compute the next head position.
	newHead := Coordinates{
		Row: g.Snake.Head().Row + g.Snake.Direction.Row,
		Col: g.Snake.Head().Col + g.Snake.Direction.Col,
	}

	// Check wall collisions.
	if newHead.Row < 0 || newHead.Row >= g.Height || newHead.Col < 0 || newHead.Col >= g.Width {
		g.GameState = GameOver
		return
	}

	// Determine if the snake is about to eat food.
	grow := false
	if newHead == g.Food {
		grow = true
	}

	// Check self collision.
	collision := false
	if grow {
		// When growing, the tail remains.
		for _, part := range g.Snake.Body {
			if part == newHead {
				collision = true
				break
			}
		}
	} else {
		// When not growing, the tail will be removed.
		if len(g.Snake.Body) > 1 {
			for _, part := range g.Snake.Body[:len(g.Snake.Body)-1] {
				if part == newHead {
					collision = true
					break
				}
			}
		}
	}
	if collision {
		g.GameState = GameOver
		return
	}

	// Move the snake.
	g.Snake.Move(grow)

	// Generate new food if it was eaten.
	if grow {
		g.Food = GenerateFood(g.Width, g.Height, g.Snake)
	}
}
