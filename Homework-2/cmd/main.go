package main

import (
	"bufio"
	"fmt"
	"homework-2/internal/service"
	"homework-2/internal/storage"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)
	statusCh := make(chan string)
	outputCh := make(chan string)
	addCh := make(chan string)
	defer close(statusCh)
	defer close(outputCh)
	defer close(addCh)

	stor, err := storage.PvzNew()
	if err != nil {
		fmt.Println(err)
		fmt.Println("не удалось подключиться к хранилищу")
		return
	}
	serv := service.PvzNew(&stor)

	// Горутина для обработки сигналов
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigCh)
	go func() {
		for sig := range sigCh {
			fmt.Printf("Получен сигнал %s. Завершение работы...\n", sig)
			time.After(1 * time.Second)
			close(sigCh)
			close(statusCh)
			os.Exit(0)
		}
	}()

	// Горутина для вывода информации о текущем состоянии
	go func() {
		for {
			select {
			case status := <-statusCh:
				fmt.Println(status)
			case <-time.After(5 * time.Second):
				fmt.Println("Ожидание...")
			}
		}
	}()

	// Горутина для вывода ПВЗ по id
	go func() {
		for {
			select {
			case <-outputCh:
				serv.Output(statusCh)
			default:
				continue
			}
		}
	}()

	// Горутина для добавления ПВЗ
	go func() {
		for {
			select {
			case <-addCh:
				serv.Add(statusCh)
			default:
				continue
			}
		}
	}()

	//Интерактиная часть приложения
	for {
		fmt.Println("Выберите действие\nadd - Добавить ПВЗ\noutput - Найти ПВЗ по id")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при чтении ввода:", err)
			continue
		}
		input = strings.TrimSpace(input)
		switch input {
		case "add":
			statusCh <- "status: запрошена команда add"
			addCh <- "add"
		case "output":
			statusCh <- "status: запрошена команда output"
			outputCh <- "output"
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
