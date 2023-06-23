package dbrepo

import (
	"context"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation 创建上下文，设置超时时间为3秒
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// 创建上下文，设置超时时间为3秒
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newID int
	stmt := `insert 
into
    reservations
    (first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at)    
values
    ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	// 使用m.DB.QueryRowContext()方法执行SQL语句，并通过Scan()方法将查询结果扫描到newID变量中。
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction 插入一条房间限制信息到数据库中
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert 
into
    room_restrictions
    (start_date, end_date,room_id,reservation_id,created_at,updated_at,restriction_id)  
values
    ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDates 查询指定的房间是否可用
func (m *postgresDBRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `select
    count(id) 
from
    room_restrictions 
where
    room_id = $1
    $2 < end_date 
    and $3 > start_date`

	// 执行查询，并将结果存储在numRows变量中
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	// 如果numRows为0，表示日期范围内没有房间限制记录，说明可用性为true
	if numRows == 0 {
		return true, nil
	}
	// 否则，有房间限制记录，说明不可用性为false
	return false, nil
}
