package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"homework-2/internal/model"
	"io"
	"os"
	"time"
)

const orderStorageName = "OrderStorage"

type OrderStorage struct {
	storage *os.File
}

func New() (OrderStorage, error) {
	file, err := os.OpenFile(orderStorageName, os.O_CREATE, 0777)
	if err != nil {
		return OrderStorage{}, err
	}
	return OrderStorage{storage: file}, nil
}

// Accept...
func (s *OrderStorage) Accept(input model.DataInputAccept) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	newData := DataDTO{
		ID:           input.ID,
		IdPVZ:        input.IdPVZ,
		IdClient:     input.IdClient,
		DateStorage:  input.DateStorage,
		DateIssue:    time.Now(),
		IsReturn:     false,
		IsIssued:     false,
		IsIssuedBack: false,
		IsDelete:     false,
	}
	all = append(all, newData)
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// Return...
func (s *OrderStorage) Return(id int) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for indx, data := range all {
		if data.ID == id {
			all[indx].IsReturn = true
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// Give...
func (s OrderStorage) Give(sliceId []int, idClient int) error { // strInt has the format "1, 2, 3, 5, 7"

	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for index := range sliceId {
		for jndx := range all {
			if sliceId[index] == all[jndx].ID {
				all[jndx].IsIssued = true
				all[jndx].IsIssuedBack = false
				all[jndx].DateIssue = time.Now()
			}
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// Refund...
func (s OrderStorage) Refund(id int, idClient int, idPVZ int) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for index := range all {
		if id == all[index].ID {
			all[index].IsIssued = false
			all[index].IsIssuedBack = true
			err = writeBytes(all)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("ошибка! заказ не найден")
}

// List ...
func (s OrderStorage) List(idClient int, idpvz int, n int) ([]int, error) {
	all, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	sliceId := []int{}
	count := 0
	for index := len(all) - 1; index >= 0; index-- {
		if idClient == all[index].IdClient && idpvz == all[index].IdPVZ && !all[index].IsIssued && !all[index].IsReturn && !all[index].IsDelete {
			sliceId = append(sliceId, all[index].ID)
			count += 1
			if count == n {
				return sliceId, nil
			}
		}
	}
	return sliceId, nil
}

// RefundList
func (s OrderStorage) RefundList(idpvz int) ([]int, error) {
	sliceId := []int{}
	all, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	count := 1
	for index := len(all) - 1; index >= 0; index-- {
		if idpvz == all[index].IdPVZ && all[index].IsIssuedBack {
			sliceId = append(sliceId, all[index].ID)
			fmt.Println(count, ".", "id:", all[index].ID, " dateRefund:", all[index].DateIssue)
			count++
		}
	}
	return sliceId, nil
}

func writeBytes(data []DataDTO) error {
	rawBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(orderStorageName, rawBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (s *OrderStorage) Delete(id int) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for indx, data := range all {
		if data.ID == id {
			all[indx].IsDelete = true
		}
	}

	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// ListAll return all data from storage
func (s *OrderStorage) ListAll() ([]DataDTO, error) {
	reader := bufio.NewReader(s.storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var data []DataDTO
	if len(rawBytes) == 0 {
		return data, nil
	}
	err = json.Unmarshal(rawBytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
