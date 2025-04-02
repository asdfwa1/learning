package logger

import (
	"fmt"
	"log/slog"
	"os"
	"v3/leaderbord"
)

type GameLogger struct {
	Logger *slog.Logger
}

func NewLogger() *GameLogger {
	file, err := os.OpenFile("game.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Errorf("ошибка в логируемом файле")
	}
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return &GameLogger{
		Logger: slog.New(handler),
	}
}

func (glog *GameLogger) LogRes(secNum, attempts int, isWin bool) {
	if isWin {
		glog.Logger.Info("Игрок победил",
			"secret_number", secNum,
			"attempts_used", attempts,
			"result", "WIN")

		lb, _ := leaderbord.LoadLeaderBoard()
		lb.AddRecord(secNum, attempts, "WIN")
	} else {
		glog.Logger.Info("Игрок не угадал число",
			"secret_number", secNum,
			"attempts_used", attempts,
			"result", "LOSE")
	}
}
