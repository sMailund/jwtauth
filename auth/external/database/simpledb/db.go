package simpledb

import (
	"fmt"
	"jwt-auth/auth/core/domainEntities"
	"sync"
)

type Db struct{
	sync.Mutex
	users map[string]domainEntities.User
}

func NewDb() *Db {
	var db Db
	db.users = make(map[string]domainEntities.User)
	return &db
}

func (d *Db) GetUserByName(name string) (domainEntities.User, error) {
	val, exists := d.users[name]
	if !exists {
		return domainEntities.User{}, fmt.Errorf("no user with username %v\n", name)
	}
	return val, nil
}

func (d *Db) Save(user domainEntities.User) error {
	d.Lock()
	defer d.Unlock()

	_, exists := d.users[user.Name]
	if exists {
		return fmt.Errorf("user with username %v already exists\n", user.Name)
	}
	d.users[user.Name] = user
	return nil
}
