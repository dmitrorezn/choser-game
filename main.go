package main

import (
	"errors"
	"fmt"
	"strings"
)

type Scene struct {
	Description string
	Win         bool
	Dead        bool
}

type Game struct {
	Scene       Scene
	LeftOption  option
	RightOption option

	Left  *Game
	Right *Game
}

func (g *Game) Print() string {
	if strings.Count(g.Scene.Description, "%d") == 2 {
		return fmt.Sprintf(g.Scene.Description, g.LeftOption, g.RightOption)
	}

	return g.Scene.Description
}

type option uint8

const (
	first option = iota + 1
	second
	end
)

var errDead = errors.New("DEAD")

func (g *Game) AddLeft(next *Game) *Game {
	g.Left = next

	return next
}
func (g *Game) AddRight(next *Game) *Game {
	g.Right = next

	return next
}

func (g *Game) MoveNext(o option) (*Game, error) {
	var next *Game
	switch o {
	case g.LeftOption:
		next = g.Left
	case g.RightOption:
		next = g.Right
	}
	if next == nil || next.Scene.Dead {
		return next, errDead
	}

	return next, nil
}

func main() {
	start := Scene{
		Description: "Start, go to cage or go to wood? (enter %d or %d)",
	}

	cage := &Game{
		Scene: Scene{
			Description: "Enter cage, yes or no? (%d or %d)",
		},
		LeftOption:  first,
		RightOption: second,
	}

	g := &Game{
		Scene:       start,
		LeftOption:  first,
		RightOption: second,
	}

	left := g.AddLeft(cage)
	left.AddLeft(&Game{
		Scene: Scene{
			Description: "U entered cage\nMake a fire or go deeper? (%d or %d) ",
		},
		LeftOption:  first,
		RightOption: second,
	})
	left.AddRight(&Game{
		Scene: Scene{
			Description: "U wen to the cache",
			Win:         true,
		},
		LeftOption:  first,
		RightOption: second,
	})

	right := g.AddRight(&Game{
		Scene: Scene{
			Description: "U met bear, throw stone or ran away? (%d or %d)",
		},
		LeftOption:  first,
		RightOption: second,
	})
	right.AddRight(&Game{
		Scene: Scene{
			Description: "WIN",
			Win:         true,
		},
	})
	right.AddLeft(&Game{
		Scene: Scene{
			Description: "U were eaten by bear",
			Dead:        true,
		},
	})

	var err error
	for {
		fmt.Println(g.Print())
		var v uint8
		_, err = fmt.Scanf("%d", &v)
		if err != nil {
			fmt.Println("Err Scan")

			continue
		}
		g, err = g.MoveNext(option(v))
		if errors.Is(err, errDead) {
			fmt.Println("USER DEAD")
			break
		}
		if g.Scene.Dead {
			fmt.Println(g.Print())
			fmt.Println("USER DEAD")

			break
		}
		if g.Scene.Win {
			fmt.Println("U WIN")
			break
		}
	}
}
