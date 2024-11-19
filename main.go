package main

import (
	_ "modernc.org/sqlite"
)

func main() {
	/*
	   // настройте подключение к БД

	   store := // создайте объект ParcelStore функцией NewParcelStore
	   service := NewParcelService(store)

	   // регистрация посылки
	   client := 1
	   address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	   p, err := service.Register(client, address)

	   	if err != nil {
	   		fmt.Println(err)
	   		return
	   	}

	   // изменение адреса
	   newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	   err = service.ChangeAddress(p.Number, newAddress)

	   	if err != nil {
	   		fmt.Println(err)
	   		return
	   	}

	   // изменение статуса
	   err = service.NextStatus(p.Number)

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
	   err = service.Delete(p.Number)

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
	   p, err = service.Register(client, address)

	   	if err != nil {
	   		fmt.Println(err)
	   		return
	   	}

	   // удаление новой посылки
	   err = service.Delete(p.Number)

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
	*/
}
