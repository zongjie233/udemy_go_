package dbrepo

import (
	"errors"
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
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "me@here.ca" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil
}

// AllNewReservations 返回订单
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}

// UpdateReservation update Reservation
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

// DeleteReservation Delete Reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation update processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}
func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
