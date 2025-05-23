package repository

import "time"

type Task struct {
	ID        int
	Title     string
	Content   string
	Created   time.Time
	Completed bool
}
