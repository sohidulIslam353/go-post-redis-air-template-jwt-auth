package models

import (
	"time"

	"github.com/uptrace/bun"
)

type JobType struct {
	bun.BaseModel `bun:"table:job_types"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name,notnull"`
	Slug          string    `bun:"slug,notnull"`
	Status        int       `bun:"status,notnull,default:1"`
	CreatedAt     time.Time `bun:"created_at,default:now()"`
	UpdatedAt     time.Time `bun:"updated_at,default:now()"`
}
