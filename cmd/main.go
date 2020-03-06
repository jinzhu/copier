package main

import (
	"copier"
	"database/sql"
	"fmt"
	"time"
)

type SubModel struct {
	SubName string
}

type Model struct {
	Name      sql.NullString
	Num       sql.NullInt32
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	T         *time.Time
	T2        *time.Time
	T3        sql.NullTime
	Sub       []*SubModel
}

func (m Model) String() string {

	return fmt.Sprintf("\nModel:\n\nName: %v\nNum: %v\nCreatedAt: %v\nUpdatedAt: %v\nT: %v\nT2: %v\nT3: %v\nSub: %v",
		m.Name,
		m.Num,
		m.CreatedAt,
		m.UpdatedAt,
		m.T,
		m.T2,
		m.T3,
		m.Sub,
	)
}

type SubRequest struct {
	SubName string
}

type Request struct {
	Name      *string
	Num       int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	T         sql.NullTime
	T2        time.Time
	T3        *time.Time
	Sub       []SubRequest
}

func (m Request) String() string {
	return fmt.Sprintf("\nRequest:\n\nName: %v\nNum: %v\nCreatedAt: %v\nUpdatedAt: %v\nT: %v\nT2: %v\nT3: %v\nSub: %v",
		m.Name,
		m.Num,
		m.CreatedAt,
		m.UpdatedAt,
		m.T,
		m.T2,
		m.T3,
		m.Sub,
	)
}

func main() {
	var (
		from Request
		to   Model
	)
	now := time.Now()
	from.T.Scan(now.Add(time.Hour * 1))
	from.T2 = now.Add(time.Hour * 2)
	now2 := new(time.Time)
	*now2 = now2.Add(time.Hour)
	from.CreatedAt = &now
	from.Sub = append(from.Sub, SubRequest{SubName: "name1"}, SubRequest{SubName: "name2"})
	if err := copier.Copy(&to, from); err != nil {
		fmt.Println(err)
	}

	fmt.Println(from)
	fmt.Println(to)
}
