package ui

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"v1/currency"
)

type CLI struct {
	service currency.Service
	reader  *bufio.Reader
}

func NewCLI(service currency.Service) *CLI {
	return &CLI{
		service: service,
		reader:  bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) RUN() {
	fmt.Println("Welcome to currency converter")

	for {
		if err := c.showMenu(); err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
		if !c.askContinue() {
			break
		}
	}
}

func (c *CLI) showMenu() error {
	currencies := c.service.ListCurrencies()
	sort.Strings(currencies)

	fmt.Print("\nДоступные валюты:")
	fmt.Println()
	currencyMap := make(map[int]string)
	i := 1
	for _, cur := range currencies {
		fmt.Printf("%d. %s\n", i, cur)
		currencyMap[i] = cur
		i++
	}

	amount, err := c.readAmount()
	if err != nil {
		return err
	}

	currencyCode, err := c.readCurrency(currencyMap)
	if err != nil {
		return err
	}

	result, err := c.service.Convert(currency.ConversionRequest{
		Amount: amount,
		From:   "USD Доллар США",
		To:     currencyCode,
	})
	if err != nil {
		return err
	}

	toRes := currencyCode[:3]
	fmt.Printf("\nРезультат: %.2f USD = %.2f %s (Курс: 1 USD = %.4f %s)\n",
		result.Amount, result.ConvertedAmount, toRes, result.Rate, toRes)

	return nil
}

func (c *CLI) readAmount() (float64, error) {
	fmt.Print("\nВведите сумму в USD: ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	amount, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return 0, fmt.Errorf("невозможно преобразовать сумму")
	}

	if amount <= 0 {
		return 0, fmt.Errorf("сумма должна быть положительной")
	}

	return amount, nil
}

func (c *CLI) readCurrency(currencyMap map[int]string) (string, error) {
	fmt.Print("\nВыберите номер валюты: ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", fmt.Errorf("невозможно преобразовать номер")
	}

	code, ok := currencyMap[choice]
	if !ok {
		return "", fmt.Errorf("валюта не найдена")
	}

	return code, nil
}

func (c *CLI) askContinue() bool {
	fmt.Print("\nХотите выполнить еще одно преобразование? (y/n): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return false
	}

	return strings.ToLower(strings.TrimSpace(input)) == "y"
}
