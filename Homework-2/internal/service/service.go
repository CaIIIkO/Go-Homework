package service

import (
	"bufio"
	"errors"
	"fmt"
	"homework-2/internal/model"
	"homework-2/internal/storage"
	"os"
	"strconv"
	"strings"
	"time"
)

type storageI interface {
	Accept(input model.DataInputAccept) error
	Delete(id int) error
	Return(id int) error
	Give(sliceId []int, idClient int) error
	Refund(id int, idClient int, idPVZ int) error
	List(idClient int, idpvz int, n int) ([]int, error)
	RefundList(idpvz int) ([]int, error)
	ListAll() ([]storage.DataDTO, error)
}

type pvzStorageI interface {
	Add(input model.PvzInputAdd) error
	Output(id int) error
}

type Service struct {
	storage    storageI
	pvzStorage pvzStorageI
}

// PvzNew ...
func PvzNew(s pvzStorageI) Service {
	return Service{pvzStorage: s}
}

func (s Service) Add(statusCh chan<- string) error {
	var input model.PvzInputAdd
	reader := bufio.NewReader(os.Stdin)
	statusCh <- "status add: Ввод id"
	fmt.Println("Введите ID:")
	_, err := fmt.Scan(&input.ID)
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return nil
	}
	reader.ReadString('\n')

	statusCh <- "status add: Ввод name"
	fmt.Println("Введите имя ПВЗ:")
	input.Name, _ = reader.ReadString('\n')
	input.Name = strings.TrimRight(input.Name, "\r\n")

	statusCh <- "status add: Ввод address"
	fmt.Println("Введите адрес ПВЗ:")
	input.Address, _ = reader.ReadString('\n')
	input.Address = strings.TrimRight(input.Address, "\r\n")

	statusCh <- "status add: Ввод contact"
	fmt.Println("Введите контактные данные ПВЗ:")
	input.Contact, _ = reader.ReadString('\n')
	input.Contact = strings.TrimRight(input.Contact, "\r\n")

	statusCh <- "status add: Добавление ПВЗ"
	return s.pvzStorage.Add(input)
}

func (s Service) Output(statusCh chan<- string) error {
	var id int
	statusCh <- "status output: Ввод ID"
	fmt.Println("Введите ID:")
	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return nil
	}
	statusCh <- "status output: Вывод информации о ПВЗ"
	return s.pvzStorage.Output(id)
}

// New...
func New(s storageI) Service {
	return Service{storage: s}
}

// Create ...
func (s Service) Accept(input model.DataInputAccept) error {
	all, err := s.storage.ListAll()
	if err != nil {
		return err
	}
	for indx := range all {
		if input.ID == all[indx].ID {
			return errors.New("ошибка! заказ с таким id уже существует")
		}
	}
	if input.DateStorage.Before(time.Now()) {
		return errors.New("ошибка! дата хранения не может быть указана в прошедшем времени")
	}
	return s.storage.Accept(input)
}

// Return ...
func (s Service) Return(id int) error {
	all, err := s.storage.ListAll()
	if err != nil {
		return err
	}
	for indx, data := range all {
		if data.ID == id {
			if all[indx].IsDelete {
				return errors.New("ошибка! заказ находится в архиве, его нельзя вернуть")
			}
			if all[indx].IsReturn {
				return errors.New("ошибка! заказ уже был возращен курьеру")
			}
			if all[indx].IsIssued {
				return errors.New("ошибка! заказ был выдан клиенту, его нельзя вернуть")
			}
			if !all[indx].DateStorage.Before(time.Now()) {
				return errors.New("ошибка! у заказа ещё не истек срок хранения, его нельзя вернуть")
			}
		}
	}
	return s.storage.Return(id)
}

// Delete ...
func (s Service) Delete(id int) error {
	if id == 0 {
		return errors.New("нулевой id цели")
	}
	return s.storage.Delete(id)
}

// Give ...
func (s Service) Give(strInt string, idClient int) error {
	if strInt == "" {
		return errors.New("не введены номера заказов")
	}
	strValues := strings.Split(strInt, ",")
	var sliceId []int
	for _, strValue := range strValues {
		num, err := strconv.Atoi(strings.TrimSpace(strValue))
		if err != nil {
			return errors.New("ошибка! не верный формат ввода id")
		}
		sliceId = append(sliceId, num)
	}
	all, err := s.storage.ListAll()
	if err != nil {
		return err
	}
	for index := range sliceId {
		isOk := false
		for jndx := range all {
			if sliceId[index] == all[jndx].ID {
				if idClient != all[jndx].IdClient {
					return errors.New("ошибка! не все заказы принадлежат одному клиенту")
				}
				if all[jndx].IsDelete {
					return errors.New("ошибка! такого заказа не существует")
				}
				if all[jndx].IsIssued {
					return errors.New("ошибка! заказ уже выдан клиенту")
				}
				if all[jndx].IsReturn {
					return errors.New("ошибка! заказ был возращён курьеру")
				}
				isOk = true
				break
			}
		}
		if !isOk {
			return errors.New("ошибка! такого заказа не сущестует")
		}
	}
	return s.storage.Give(sliceId, idClient)
}

// Refund ...
func (s Service) Refund(id int, idClient int, idPVZ int) error {
	all, err := s.storage.ListAll()
	if err != nil {
		return err
	}
	for index := range all {
		if id == all[index].ID {
			if all[index].IsDelete {
				return errors.New("ошибка! заказ находится в архиве, его нельзя вернуть")
			}
			diff := time.Since(all[index].DateIssue)
			if diff > 48*time.Hour {
				return errors.New("ошибка! с момента получения товара прошло 48 часов")
			}
			if idPVZ != all[index].IdPVZ {
				return errors.New("ошибка! товар был выдан на другом ПВЗ")
			}
			if idClient != all[index].IdClient {
				return errors.New("ошибка! ID клиента не соответсвует заказу")
			}
			if all[index].IsIssuedBack {
				return errors.New("ошибка! товар уже возращен на ПВЗ")
			}
		}
	}
	return s.storage.Refund(id, idClient, idPVZ)
}

// List ...
func (s Service) List(idClient int, idpvz int, n int) ([]int, error) {
	return s.storage.List(idClient, idpvz, n)
}

// RefundList ...
func (s Service) RefundList(idpvz int) ([]int, error) {
	return s.storage.RefundList(idpvz)
}
