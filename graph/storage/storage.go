package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vikstrous/dataloadgen-example/graph/model"
)

type UserStorage struct {
	lock sync.Mutex
	data map[string]model.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		data: make(map[string]model.User),
	}
}

var (
	notFoundError = errors.New("not found")
	existsError   = errors.New("exists")
)

func (u *UserStorage) Get(userID string) (model.User, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	fmt.Println("UserStorage.Get")
	time.Sleep(time.Second)
	user, ok := u.data[userID]
	if !ok {
		return model.User{}, notFoundError
	}
	return user, nil
}

func (u *UserStorage) GetMulti(userIDs []string) ([]*model.User, []error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	fmt.Printf("UserStorage.GetMulti %d\n", len(userIDs))
	time.Sleep(time.Second)
	users := make([]*model.User, 0, len(userIDs))
	errs := make([]error, 0, len(userIDs))
	for _, userID := range userIDs {
		user, ok := u.data[userID]
		if ok {
			users = append(users, &user)
			errs = append(errs, nil)
		} else {
			users = append(users, nil)
			errs = append(errs, notFoundError)
		}
	}
	return users, errs
}

func (u *UserStorage) Put(user model.User) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	fmt.Println("UserStorage.Put")
	time.Sleep(time.Second)
	_, ok := u.data[user.ID]
	if ok {
		return existsError
	}
	u.data[user.ID] = user
	return nil
}

type TodoStorage struct {
	lock sync.Mutex
	data map[string]model.Todo
}

func NewTodoStorage() *TodoStorage {
	return &TodoStorage{
		data: make(map[string]model.Todo),
	}
}

func (t *TodoStorage) GetAll() ([]*model.Todo, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Println("TodoStorage.GetAll")
	time.Sleep(time.Second)
	todos := []*model.Todo{}
	for _, v := range t.data {
		v := v
		todos = append(todos, &v)
	}
	return todos, nil
}

func (t *TodoStorage) Put(todo model.Todo) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Println("TodoStorage.Put")
	time.Sleep(time.Second)
	_, ok := t.data[todo.ID]
	if ok {
		return existsError
	}
	t.data[todo.ID] = todo
	return nil
}
