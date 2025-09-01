package dto

type JobTypeStoreDTO struct {
	Name   string `form:"name" binding:"required,min=2,max=255"`
	Status int    `form:"status" binding:"required,oneof=0 1"`
}

type JobTypeUpdateDTO struct {
	Name   string `form:"name" binding:"required,min=2,max=255"`
	Status int    `form:"status" binding:"required,oneof=0 1"`
}
