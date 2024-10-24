package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	db "astra-test-task/database"

	"github.com/joho/godotenv"
)

// Структура для парсинга необходимых данных из json файла
type RunsResults struct {
	Runs []struct {
		Results []struct {
			RuleID    string `json:"ruleId"`
			Locations []struct {
				PhysicalLocation struct {
					ArtifactLocation struct {
						Uri string `json:"uri"`
					} `json:"artifactLocation"`
					Region struct {
						StartLine int `json:"startLine"`
						EndLine   int `json:"endLine"`
					} `json:"region"`
				} `json:"physicalLocation"`
				Properties struct {
					Severity int `json:"x-severity"`
				} `json:"properties"`
			} `json:"locations"`
		} `json:"results"`
	} `json:"runs"`
}

func main() {
	// Поулчаем путь к json файлу через первый аргумент приложения
	filePath := os.Args[1]

	// Открываем файл
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Файл report.sarif успешно открыт")

	// Откладываем закрытие файла, чтобы произвести парсинг данных
	defer jsonFile.Close()

	// Переменная для хранения декодированных данных
	var data RunsResults

	// Словарь для хранения пар вида: тип уязвимости - количество
	severityMap := map[int]int{
		0: 0,
		1: 0,
		2: 0,
	}

	// Читаем json файл и сохраняем декодированные данные
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&data); err != nil {
		log.Fatalf("Ошибка декодирования: %v", err)
	}
	log.Println("Успешное декодирование данных json")

	// Читаем .env файл
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	connUrl := os.Getenv("POSTGRES_CONN_URL")

	// Подключаемся к базе данных
	dbpool, err := db.ConnectDB(connUrl)
	if err != nil {
		log.Fatalf("Ошибка создания пула к базе данных: %v\n", err)
	}
	defer dbpool.Close()
	log.Println("Успешное подключение к базе данных")

	// Создаем таблицу warnings
	pg := db.NewPGXPool(dbpool)
	err = pg.CreateWarningsTable()
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v\n", err)
	}
	// Вычленяем данные, которые необходимо занести в бд
	for _, runs_results := range data.Runs {
		for _, res := range runs_results.Results {
			for _, loc := range res.Locations {
				// Добавляем запись в бд
				tmpWarning := db.Warning{
					RuleId:    res.RuleID,
					Uri:       loc.PhysicalLocation.ArtifactLocation.Uri,
					StartLine: loc.PhysicalLocation.Region.StartLine,
					XSeverity: loc.Properties.Severity,
				}
				err = pg.InsertWarning(&tmpWarning)
				if err != nil {
					log.Printf("Ошибка добавления строки в таблицу: %v\n", err)
				}
				// Увеличиваем счетчик количества уязвимостей по критичности
				severityMap[loc.Properties.Severity]++
			}
		}
	}

	// вывод в консоль классификации количества уязвимостей
	fmt.Println("\nКлассификация количества уязвимостей по критичности")
	fmt.Println(severityMap[2], "- высокой критичности")
	fmt.Println(severityMap[1], "- средней критичности")
	fmt.Println(severityMap[0], "- информационных")
	fmt.Println(severityMap[0]+severityMap[1]+severityMap[2], "- суммарно по всем срабатываниям")
}
