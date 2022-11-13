package store

import (
	"encoding/json"
	"io/fs"
	"os"
	"strconv"

	"refactoring/internal/model"
)

type storeData struct {
	Increment int                   `json:"increment"`
	List      map[string]model.User `json:"list"`
}

type Store struct {
	path string
	data storeData
}

func Open(path string) (*Store, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := storeData{}
	err = json.Unmarshal(f, &data)
	if err != nil {
		return nil, err
	}

	return &Store{
		path: path,
		data: data,
	}, nil
}

func (s *Store) save() {
	b, _ := json.Marshal(&s.data)
	_ = os.WriteFile(s.path, b, fs.ModePerm)
}

func (s *Store) CreateUser(user model.User) string {
	s.data.Increment++
	id := strconv.Itoa(s.data.Increment)
	s.data.List[id] = user
	s.save()
	return id
}

func (s *Store) GetUser(id string) (model.User, error) {
	return s.data.List[id], nil
}

func (s *Store) UpdateUser(id string, user model.User) error {
	s.data.List[id] = user
	s.save()
	return nil
}

func (s *Store) DeleteUser(id string) error {
	delete(s.data.List, id)
	s.save()

	return nil
}

func (s *Store) ListUsers() map[string]model.User {
	return s.data.List
}
