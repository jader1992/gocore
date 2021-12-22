package demo

import (
	"database/sql"
	"time"
)

// UserModel 定义的用户数据结构
type UserModel struct {
	UserId int
	Name   string
	Age    int
}

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdateAt     time.Time
}
