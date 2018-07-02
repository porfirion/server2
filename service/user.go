package service

type User struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UsersList []uint64
