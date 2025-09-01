package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name,notnull"`
	Slug          string    `bun:"slug,notnull"`
	Status        int       `bun:"status,notnull,default:1"`
	CreatedAt     time.Time `bun:"created_at,default:now()"`
	UpdatedAt     time.Time `bun:"updated_at,default:now()"`

	// revers join optional
	Subcategories []*Subcategory `bun:"rel:has-many,join:id=category_id"`
}
