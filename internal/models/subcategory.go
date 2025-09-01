package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Subcategory struct {
	bun.BaseModel `bun:"table:subcategories"`
	ID            int64     `bun:"id,pk,autoincrement"`
	CategoryID    int       `bun:"category_id,notnull"`
	Name          string    `bun:"name,notnull"`
	Slug          string    `bun:"slug,notnull"`
	Status        int       `bun:"status,notnull,default:1"`
	CreatedAt     time.Time `bun:"created_at,default:now()"`
	UpdatedAt     time.Time `bun:"updated_at,default:now()"`

	// Relation with Category
	Category *Category `bun:"rel:belongs-to,join:category_id=id"`
}
