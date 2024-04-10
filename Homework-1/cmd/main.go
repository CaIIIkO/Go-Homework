package main

import (
	"flag"
	"fmt"
	"homework-1/internal/model"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"time"
)

func main() {
	flagCommand := flag.String("command", "help", "command")
	flagIdPVZ := flag.Int("idpvz", 0, "idpzv for data")
	flagID := flag.Int("id", 0, "id for Data")
	flagIdClient := flag.Int("idclient", -1, "id for client")
	flagN := flag.Int("n", -1, "number of entries")
	flagSliceId := flag.String("idslice", "", "slice id for order")
	defaultDate := time.Now().Add(336 * time.Hour).Format("2006-01-02")
	flagDateStorage := flag.String("datestorage", defaultDate, "date of storage for the PVZ")

	flag.Parse()

	stor, err := storage.New()
	if err != nil {
		fmt.Println(err)
		fmt.Println("не удалось подключиться к хранилищу")
		return
	}

	serv := service.New(&stor)

	switch *flagCommand {
	case "help":
		fmt.Println("Доступные команды:")
		fmt.Println("  -commnad=help                                                          - Показать список команд")
		fmt.Println("  -command=accept    -id=0 -idpvz=0 -datestorage=00-00-0000 -idclient=0  - Принять заказ от курьера")   //-command=accept -id=0 -idpvz=0 -datestorage=00-00-0000 -idclient=0
		fmt.Println("  -command=return    -id=0                                               - Вернуть заказ курьеру")      //-command=return -id=0
		fmt.Println("  -command=give      -idslice=\"1, 2, 3\" -idclient=0                      - Выдать заказ клиенту")     //-command=give -idslice="1, 2, 3" -idclient=0
		fmt.Println("  -command=refund    -id=0 -idclient=0 -idpvz=0                          - Принять возврат от клиента") //-command=refund -id=0 -idclient=0 -idpvz = 0
		fmt.Println("  -command=list      -idclient=0 -idpvz=0 -n=0                           - Получить список заказов")    //-command=list -idclient=0 -idpvz=0 -n=0
		fmt.Println("  -command=refunlist -idpvz=0                                            - Получить список возвратов")  //-command=refunlist -idpvz=0
	case "accept":
		parsedDate, err := time.Parse("01-02-2006", *flagDateStorage)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = serv.Accept(model.DataInputAccept{
			ID:          *flagID,
			IdPVZ:       *flagIdPVZ,
			DateStorage: parsedDate,
			IdClient:    *flagIdClient,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Товар добавлен на ПВЗ")
	case "return":
		err = serv.Return(*flagID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Товар возращён курьеру")
	case "give":
		err = serv.Give(*flagSliceId, *flagIdClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Товары выданы клиенту")
	case "refund":
		err = serv.Refund(*flagID, *flagIdClient, *flagIdPVZ)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Товар возращен на ПВЗ")
	case "list":
		list, err := serv.List(*flagIdClient, *flagIdPVZ, *flagN)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказы:", list)
	case "refundlist":
		_, err := serv.RefundList(*flagIdPVZ)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	default:
		fmt.Println("неизвестная команда")
	}

}
