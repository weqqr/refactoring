package store

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"strconv"
	"sync"

	"refactoring/internal/model"
)

var ErrUserNotFound = errors.New("user_not_found")

type storeData struct {
	Increment int                   `json:"increment"`
	List      map[string]model.User `json:"list"`
}

type Store struct {
	path string
	lock *sync.RWMutex
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
		lock: &sync.RWMutex{},
		data: data,
	}, nil
}

func (s *Store) save() error {
	b, err := json.Marshal(&s.data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, fs.ModePerm)
}

func (s *Store) CreateUser(user model.User) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data.Increment++
	id := strconv.Itoa(s.data.Increment)
	s.data.List[id] = user

	return id, s.save()
}

func (s *Store) GetUser(id string) (model.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if user, ok := s.data.List[id]; ok {
		return user, nil
	}
	return model.User{}, ErrUserNotFound
}

func (s *Store) UpdateUser(id string, user model.User) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.data.List[id]; !ok {
		return ErrUserNotFound
	}

	s.data.List[id] = user

	return s.save()
}

func (s *Store) DeleteUser(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.data.List[id]; !ok {
		return ErrUserNotFound
	}

	delete(s.data.List, id)
	s.save()

	return s.save()
}

func (s *Store) ListUsers() map[string]model.User {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data.List
}
