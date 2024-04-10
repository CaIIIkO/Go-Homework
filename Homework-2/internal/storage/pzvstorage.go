package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"homework-2/internal/model"
	"io"
	"os"
	"sync"
)

const pvzStorageName = "PVZstorage"

type PvzStorage struct {
	storage *os.File
	mu      sync.Mutex
}

func PvzNew() (PvzStorage, error) {
	file, err := os.OpenFile(pvzStorageName, os.O_CREATE, 0777)
	if err != nil {
		return PvzStorage{}, err
	}
	return PvzStorage{storage: file}, nil
}

// Add добавляет новый ПВЗ в хранилище
func (ps *PvzStorage) Add(input model.PvzInputAdd) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	all, err := ps.ListAll()
	if err != nil {
		return err
	}
	for index := len(all) - 1; index >= 0; index-- {
		if input.ID == all[index].ID {
			fmt.Println("Ошибка! ПВЗ с таким id уже существует")
			return nil
		}
	}

	newData := PvzDTO{
		ID:      input.ID,
		Name:    input.Name,
		Address: input.Address,
		Contact: input.Contact,
	}
	all = append(all, newData)
	err = ps.pvzWriteBytes(all)
	if err != nil {
		return err
	}
	fmt.Println("ПВЗ добавлен")
	return nil
}

func (ps *PvzStorage) Output(id int) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	all, err := ps.ListAll()
	if err != nil {
		return err
	}
	for index, data := range all {
		fmt.Println(index)
		if id == data.ID {
			fmt.Println("ID        :", all[index].ID)
			fmt.Println("Name      :", all[index].Name)
			fmt.Println("Address   :", all[index].Address)
			fmt.Println("Contact   :", all[index].Contact)
			return nil
		}
	}
	fmt.Println("Ошибка! ПВЗ с таким id нет")
	return nil
}

func (ps *PvzStorage) pvzWriteBytes(data []PvzDTO) error {
	rawBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(pvzStorageName, rawBytes, 0777)
	if err != nil {
		return err
	}
	ps.storage.Close()
	ps.storage, err = os.OpenFile(pvzStorageName, os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	return nil
}

// ListAll return all pvz from storage
func (ps *PvzStorage) ListAll() ([]PvzDTO, error) {
	reader := bufio.NewReader(ps.storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var data []PvzDTO
	if len(rawBytes) == 0 {
		return data, nil
	}
	err = json.Unmarshal(rawBytes, &data)
	if err != nil {
		return nil, err
	}
	ps.storage.Close()
	ps.storage, err = os.OpenFile(pvzStorageName, os.O_CREATE, 0777)
	if err != nil {
		return data, err
	}
	return data, nil
}
