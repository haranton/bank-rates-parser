package scraper

import (
	"bank-rates-parser/internal/models"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// Scraper - основной компонент для скрапинга вкладов
type Scraper struct {
	wd     selenium.WebDriver
	logger *slog.Logger
}

// NewScraper создает новый экземпляр скрапера
func NewScraper(logger *slog.Logger) *Scraper {

	scr := Scraper{
		logger: logger,
	}

	scr.Initialize()

	return &scr
}

// Initialize инициализирует WebDriver
func (s *Scraper) Initialize() {
	fmt.Println("Проверяем подключение к Selenium...")

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://chrome:4444/wd/hub")
	if err != nil {
		msgErr := fmt.Errorf("ошибка подключения: %v", err)
		panic(msgErr)
	}

	s.wd = wd

	fmt.Println("Успешно подключились!")
}

// Close закрывает WebDriver
func (s *Scraper) Close() {
	if s.wd != nil {
		s.wd.Quit()
	}
}

// ScrapeDeposits основной метод для сбора данных о вкладах
func (s *Scraper) ScrapeDeposits() ([]models.BankRate, error) {
	if s.wd == nil {
		return nil, fmt.Errorf("WebDriver не инициализирован")
	}

	// Открываем страницу
	fmt.Println("Загружаем страницу...")
	err := s.wd.Get("https://www.sravni.ru/vklady/")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки страницы: %v", err)
	}

	fmt.Println("Страница загружена!")
	time.Sleep(3 * time.Second)

	// Нажимаем кнопку "Показать еще"
	err = s.clickShowMore()
	if err != nil {
		fmt.Printf("Предупреждение: %v\n", err)
	}

	// Собираем карточки
	cards, err := s.collectCards()
	if err != nil {
		return nil, fmt.Errorf("ошибка сбора карточек: %v", err)
	}

	s.PrintCards(cards)

	return cards, nil
}

// clickShowMore нажимает кнопку "Показать еще"
func (s *Scraper) clickShowMore() error {
	fmt.Println("Ищем и нажимаем кнопку 'Показать еще'...")

	buttons, err := s.wd.FindElements(selenium.ByCSSSelector, "button")
	if err != nil {
		return fmt.Errorf("не нашли кнопки: %v", err)
	}

	for _, btn := range buttons {
		text, _ := btn.Text()
		if strings.Contains(text, "Показать еще") {
			fmt.Printf("Нашли кнопку: %s\n", text)

			err := btn.Click()
			if err != nil {
				return fmt.Errorf("ошибка при нажатии кнопки: %v", err)
			}

			fmt.Println("Нажали кнопку 'Показать еще'")
			time.Sleep(3 * time.Second)
			return nil
		}
	}

	return fmt.Errorf("не нашли кнопку 'Показать еще'")
}

// collectCards собирает данные карточек в структуры
func (s *Scraper) collectCards() ([]models.BankRate, error) {
	elements, err := s.wd.FindElements(selenium.ByCSSSelector, "div.DepositCard_wrapper__jKpqw")
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска карточек: %v", err)
	}

	var cards []models.BankRate

	for _, card := range elements {
		text, err := card.Text()
		if err != nil {
			fmt.Printf("Ошибка чтения карточки: %v\n", err)
			continue
		}

		deposit := s.parseCardData(text)
		if deposit.BankName != "" {
			cards = append(cards, deposit)
		}
	}

	return cards, nil
}

// parseCardData парсит данные из текста карточки
func (s *Scraper) parseCardData(text string) models.BankRate {
	lines := strings.Split(text, "\n")
	var deposit models.BankRate

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		// Определяем банк (обычно первая строка)
		if deposit.BankName == "" && !strings.Contains(line, "%") && !strings.Contains(line, "₽") {
			deposit.BankName = line
			continue
		}

		// Определяем название вклада (обычно вторая строка после банка)
		if deposit.DepositName == "" && deposit.BankName != "" && !strings.Contains(line, "%") && !strings.Contains(line, "₽") {
			deposit.DepositName = line
			continue
		}

		// Определяем ставку
		if deposit.Rate == 0 && strings.Contains(line, "%") {
			// Удаляем все не-цифровые символы, кроме точки и запятой
			cleaned := strings.TrimSpace(line)
			cleaned = strings.ReplaceAll(cleaned, "%", "")
			cleaned = strings.ReplaceAll(cleaned, " ", "")

			// Пробуем преобразовать в float
			if rate, err := strconv.ParseFloat(cleaned, 32); err == nil {
				deposit.Rate = float32(rate)
			} else {
				fmt.Printf("ОШИБКА: Не удалось преобразовать '%s' в число: %v\n", cleaned, err)
			}
			continue
		}

		// Если нашли все данные, выходим
		if deposit.BankName != "" && deposit.DepositName != "" && deposit.Rate != 0 {
			break
		}
	}

	return deposit
}

// PrintCards выводит карточки в форматированном виде
func (s *Scraper) PrintCards(cards []models.BankRate) {
	fmt.Println("\nСобранные данные карточек:")
	fmt.Println("==========================================")
	fmt.Printf("Всего карточек: %d\n\n", len(cards))

	for i, card := range cards {
		fmt.Printf("Карточка %d:\n", i+1)
		fmt.Printf("  Банк: %s\n", card.BankName)
		fmt.Printf("  Вклад: %s\n", card.DepositName)
		fmt.Printf("  Ставка: %.2f%%\n", card.Rate)
		fmt.Println("  ------------------------------------")
	}
}
