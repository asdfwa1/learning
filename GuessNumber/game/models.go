package game

type GameState struct {
	SecretNum int
	MaxTry    int
	TryStep   int
	History   []string
	IsCorrect bool
}

func (g *GameState) IsGameover() bool {
	return g.TryStep <= 0 || g.IsWon()
}

func (g *GameState) IsWon() bool {
	return g.IsCorrect
}
