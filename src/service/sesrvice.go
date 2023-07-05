package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"test/src/config"
)

func Run() {
	// Запись и чтение конфигураций в окружение
	config.SetEnv()
	config := config.GetEnv()

	// Подключение к базе данных
	sqlReq := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbName, config.DbPassword, config.DbSchema)
	fmt.Println("sss", sqlReq)
	db, err := sqlx.Connect("postgres", sqlReq)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось проверить соединение с базой данных: %v", err)
	}
	// Создание маршрутизатора Gin
	router := gin.Default()

	// Обработчики endpoints
	GetReportHandler(router, db)
	SetReportHandler(router, db)
	GetObservTimeHandler(router, db)

	// Запуск сервера на заданном порту
	err = router.Run(":" + config.Port)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
