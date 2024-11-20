package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func openDB() (*sql.DB, func(), error) {
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DNS"))
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
		panic(err)
	}

	db, closeFunc, err := openDB()
	if err != nil {
		panic(err)
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
