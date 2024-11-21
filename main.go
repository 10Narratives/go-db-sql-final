package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func setDefault(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

func openDB(driver, dns string) (*sql.DB, func(), error) {
	db, err := sql.Open(driver, dns)
	if err != nil {
		return nil, nil, err
	}

	closeFunc := func() {
		_ = db.Close()
	}

	return db, closeFunc, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	setDefault("DB_DRIVER", "postgres")
	setDefault("DB_DNS", "example.db")

	databaseDriver := os.Getenv("DB_DRIVER")
	databaseDNS := os.Getenv("DB_DNS")

	db, closeFunc, err := openDB(databaseDriver, databaseDNS)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closeFunc()

	store := NewParcelStore(db)
	service := NewParcelService(store)

	// регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(int64(client), address)

	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(int(p.Number), newAddress)

	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение статуса
	err = service.NextStatus(int(p.Number))

	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	err = service.PrintClientParcels(client)

	if err != nil {
		fmt.Println(err)
		return
	}

	// попытка удаления отправленной посылки
	err = service.Delete(int(p.Number))

	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// предыдущая посылка не должна удалиться, т.к. её статус НЕ «зарегистрирована»
	err = service.PrintClientParcels(client)

	if err != nil {
		fmt.Println(err)
		return
	}

	// регистрация новой посылки
	p, err = service.Register(int64(client), address)

	if err != nil {
		fmt.Println(err)
		return
	}

	// удаление новой посылки
	err = service.Delete(int(p.Number))

	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// здесь не должно быть последней посылки, т.к. она должна была успешно удалиться
	err = service.PrintClientParcels(client)

	if err != nil {
		fmt.Println(err)
		return
	}

}
