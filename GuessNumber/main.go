package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
	"v3/game"
	"v3/leaderbord"
	"v3/logger"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Guess Number")
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.NewSize(600, 650))

	gameLogger := logger.NewLogger()
	gameService := game.NewGameService(gameLogger)

	title := widget.NewLabel("Угадай число от 0 до 10000")
	title.TextStyle = fyne.TextStyle{Bold: true}

	infoLabel := widget.NewLabel("")
	historyLabel := widget.NewLabel("")
	historyScroll := container.NewVScroll(historyLabel)
	historyScroll.SetMinSize(fyne.NewSize(0, 325))

	guessEntry := widget.NewEntry()
	guessEntry.SetPlaceHolder("Введите число")

	startButton := widget.NewButton("Новая игра", func() {
		gameService.StartNewGame()
		infoLabel.SetText("У вас " + strconv.Itoa(gameService.State.MaxTry) + " попыток!")
		historyLabel.SetText("")
		guessEntry.Enable()
	})

	guessButton := widget.NewButton("Проверить", func() {
		guess, err := strconv.Atoi(guessEntry.Text)
		if err != nil || guess < 0 || guess > 10000 {
			infoLabel.SetText("Ошибка! Введите число от 0 до 10000")
			return
		}

		result := gameService.ProcessGuess(guess)
		infoLabel.SetText(result)

		historyText := ""
		for _, entry := range gameService.State.History {
			historyText += entry + "\n"
		}
		historyLabel.SetText(historyText)

		if gameService.State.IsGameover() {
			guessEntry.Disable()
		}
		guessEntry.SetText("")
	})

	guessEntry.OnSubmitted = func(s string) {
		guessButton.OnTapped()
	}

	gameService.StartNewGame()
	infoLabel.SetText("У вас " + strconv.Itoa(gameService.State.MaxTry) + " попыток!")

	recordsButton := widget.NewButton("Рекорды", func() {
		lb, err := leaderbord.LoadLeaderBoard()
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		recordsTable := [][]string{}
		for i, record := range lb.Records {
			t, _ := time.Parse(time.RFC3339, record.Timestamp)
			formatTime := t.Format("02.01.2006 15:04")

			recordsTable = append(recordsTable, []string{
				strconv.Itoa(i + 1),
				strconv.Itoa(record.Attempts),
				formatTime,
			})
		}

		tableHeader := container.NewGridWithColumns(3,
			widget.NewLabelWithStyle("Место", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Попытки", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Дата", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)

		table := widget.NewTable(
			func() (int, int) { return len(recordsTable), 3 },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(tci widget.TableCellID, o fyne.CanvasObject) {
				label := o.(*widget.Label)
				label.SetText(recordsTable[tci.Row][tci.Col])
				label.Alignment = fyne.TextAlignCenter
			},
		)
		table.SetColumnWidth(0, 125)
		table.SetColumnWidth(1, 125)
		table.SetColumnWidth(2, 200)
		table.OnSelected = func(id widget.TableCellID) {
			table.Unselect(id)
		}
		content := container.NewBorder(
			tableHeader,
			nil, nil, nil,
			table,
		)

		d := dialog.NewCustom("Таблица рекордов", "Закрыть", content, myWindow)
		d.Resize(fyne.NewSize(500, 350))
		d.Show()
	})

	content := container.NewVBox(
		recordsButton,
		title,
		startButton,
		widget.NewLabel("Ваше число:"),
		guessEntry,
		guessButton,
		infoLabel,
		widget.NewLabel("История попыток:"),
		historyScroll,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
