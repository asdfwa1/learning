package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"v3/algorithm"
	"v3/logger"
)

type GameService struct {
	State  *GameState
	Logger *logger.GameLogger
}

func NewGameService(logger *logger.GameLogger) *GameService {
	return &GameService{
		State: &GameState{
			MaxTry: algorithm.FindOptimalTry(10000) + 1,
		},
		Logger: logger,
	}
}

func (gser *GameService) StartNewGame() {
	gser.State.SecretNum = rand.Intn(10000)
	gser.State.TryStep = gser.State.MaxTry
	gser.State.History = []string{}
	gser.State.IsCorrect = false
}

func (gser *GameService) ProcessGuess(guess int) string {
	gser.State.TryStep--
	attemptNum := gser.State.MaxTry - gser.State.TryStep
	history := "Попытка " + strconv.Itoa(attemptNum) + ": " + strconv.Itoa(guess)
	gser.State.History = append(gser.State.History, history)
	if guess == gser.State.SecretNum {
		gser.State.IsCorrect = true
		gser.Logger.LogRes(gser.State.SecretNum, attemptNum, true)
		return fmt.Sprintf("Поздравляем! Вы угадали число %d за %d попыток!",
			gser.State.SecretNum, attemptNum)
	}
	if gser.State.TryStep == 0 {
		gser.Logger.LogRes(gser.State.SecretNum, attemptNum, false)
		return fmt.Sprintf("Неудача! Вы не отгадали число %d за %d попыток!",
			gser.State.SecretNum, attemptNum)
	}
	if guess < gser.State.SecretNum {
		return "Мое число больше! Осталось попыток: " + strconv.Itoa(gser.State.TryStep)
	}
	return "Мое число меньше! Осталось попыток: " + strconv.Itoa(gser.State.TryStep)
}
