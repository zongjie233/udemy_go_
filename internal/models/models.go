package models

import (
	"time"
)

// User 是用户模型，包含了用户的基本信息和权限级别
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room 是房屋模型，包含了房间的基本信息
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction 是限制模型，包含了不同类型的限制信息
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation 是预订模型，包含了预订的详细信息，包括预订人信息、房型信息和预订时间段
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomRestriction 是房间限制模型，包含了房间的限制信息，包括限制时间段和关联的预订和限制类型
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Restriction   Restriction
	Reservation   Reservation
}

// MailData email 数据结构
type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
