package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// Структура для хранения данных карточки вклада
type DepositCard struct {
	BankName    string
	DepositName string
	Rate        string
	Income      string
}

// Scraper - основной компонент для скрапинга вкладов
type Scraper struct {
	wd selenium.WebDriver
}

// NewScraper создает новый экземпляр скрапера
func NewScraper() *Scraper {
	return &Scraper{}
}

// Initialize инициализирует WebDriver
func (s *Scraper) Initialize() error {
	fmt.Println("Проверяем подключение к Selenium...")

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}

	s.wd = wd
	fmt.Println("Успешно подключились!")
	return nil
}

// Close закрывает WebDriver
func (s *Scraper) Close() {
	if s.wd != nil {
		s.wd.Quit()
	}
}

// ScrapeDeposits основной метод для сбора данных о вкладах
func (s *Scraper) ScrapeDeposits() ([]DepositCard, error) {
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
func (s *Scraper) collectCards() ([]DepositCard, error) {
	elements, err := s.wd.FindElements(selenium.ByCSSSelector, "div.DepositCard_wrapper__jKpqw")
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска карточек: %v", err)
	}

	var cards []DepositCard

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
func (s *Scraper) parseCardData(text string) DepositCard {
	lines := strings.Split(text, "\n")
	var deposit DepositCard

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
		if deposit.Rate == "" && strings.Contains(line, "%") {
			deposit.Rate = line
			continue
		}

		// Определяем доход (содержит "₽" и обычно идет после ставки)
		if deposit.Income == "" && strings.Contains(line, "₽") {
			deposit.Income = line
			continue
		}

		// Если нашли все данные, выходим
		if deposit.BankName != "" && deposit.DepositName != "" && deposit.Rate != "" && deposit.Income != "" {
			break
		}
	}

	return deposit
}

// PrintCards выводит карточки в форматированном виде
func (s *Scraper) PrintCards(cards []DepositCard) {
	fmt.Println("\nСобранные данные карточек:")
	fmt.Println("==========================================")
	fmt.Printf("Всего карточек: %d\n\n", len(cards))

	for i, card := range cards {
		fmt.Printf("Карточка %d:\n", i+1)
		fmt.Printf("  Банк: %s\n", card.BankName)
		fmt.Printf("  Вклад: %s\n", card.DepositName)
		fmt.Printf("  Ставка: %s\n", card.Rate)
		fmt.Printf("  Доход: %s\n", card.Income)
		fmt.Println("  ------------------------------------")
	}
}
