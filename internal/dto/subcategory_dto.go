package dto

type SubcategoryStoreDTO struct {
	Name       string `form:"name" binding:"required,min=2,max=255"`
	CategoryID int    `form:"category_id" binding:"required"`
	Status     int    `form:"status" binding:"required,oneof=0 1"`
}

type SubcategoryUpdateDTO struct {
	Name       string `form:"name" binding:"required,min=2,max=255"`
	CategoryID int    `form:"category_id" binding:"required"`
	Status     int    `form:"status" binding:"required,oneof=0 1"`
}
