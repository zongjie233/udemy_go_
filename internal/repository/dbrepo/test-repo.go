package dbrepo

import (
	"github.com/zongjie233/udemy_lesson/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation 创建上下文，设置超时时间为3秒
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

// InsertRoomRestriction 插入一条房间限制信息到数据库中
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID 查询指定的房间是否可用
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms 根据时间返回可用的房型
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID 通过id获取房间类型
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
