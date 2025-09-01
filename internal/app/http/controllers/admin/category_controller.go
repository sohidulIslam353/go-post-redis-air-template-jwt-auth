package controllers

import (
	"gin-app/config"
	"gin-app/internal/dto"
	"gin-app/internal/models"
	"gin-app/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminCategoryList(c *gin.Context) {
	successMsg := c.Query("success")
	// Filters
	search := c.Query("search")
	status := c.Query("status")
	createdAt := c.Query("created_at")

	// Cursor pagination
	lastID, limit := utils.GetCursorPagination(c)

	// Base query
	query := config.DB.NewSelect().Model((*models.Category)(nil))

	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if createdAt != "" {
		query = query.Where("DATE(created_at) = ?", createdAt)
	}

	// Cursor condition
	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	query = query.Order("id ASC").Limit(limit)

	// Fetch
	var categories []models.Category
	if err := query.Scan(c, &categories); err != nil {
		c.HTML(http.StatusInternalServerError, "category_list.html", gin.H{
			"title": "Category List",
			"error": "Failed to fetch categories: " + err.Error(),
			"data":  []models.Category{},
		})
		return
	}

	// Next cursor
	var nextCursor int64
	if len(categories) > 0 {
		nextCursor = categories[len(categories)-1].ID
	}

	// Base query for count
	countQuery := config.DB.NewSelect().Model((*models.JobType)(nil))
	// Count total items
	totalCount, err := countQuery.Count(c)
	if err != nil {
		totalCount = 0
	}
	// Render
	c.HTML(http.StatusOK, "category_list.html", gin.H{
		"title":      "Category List",
		"PageName":   "category_list",
		"data":       categories,
		"nextCursor": nextCursor,
		"limit":      limit,
		"total":      totalCount,
		"filters": gin.H{
			"search":     search,
			"status":     status,
			"created_at": createdAt,
		},
		"success": successMsg, // pass to template
	})
}

// Create page
func AdminCategoryCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "category_create.html", gin.H{
		"title":    "Create Category",
		"PageName": "category_create",
	})
}

// Category Store
func AdminCategoryStore(c *gin.Context) {
	var input dto.CategoryStoreDTO

	//  Bind + Validate
	if valid, errs := utils.ValidateStruct(c, &input); !valid {
		c.HTML(http.StatusBadRequest, "category_create.html", gin.H{
			"title":  "Create Category",
			"errors": errs,
			"data":   input, // old input re-fill
		})
		return
	}

	//  Slug Creating
	slug := utils.MakeSlug(input.Name)
	if slug == "" {
		slug = "not-available"
	}

	//  Model Creating
	category := models.Category{
		Name:      input.Name,
		Slug:      slug,
		Status:    input.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Database value insert
	if _, err := config.DB.NewInsert().Model(&category).Exec(c); err != nil {
		c.HTML(http.StatusInternalServerError, "category_create.html", gin.H{
			"title":  "Create Category",
			"errors": map[string]string{"DB": "Failed to create category: " + err.Error()},
			"data":   input,
		})
		return
	}

	// Success response → empty form + success msg
	c.HTML(http.StatusOK, "category_create.html", gin.H{
		"title":   "Create Category",
		"success": "Category created successfully!",
		"errors":  map[string]string{},
		"data":    dto.CategoryStoreDTO{}, // clear form
	})
}

// Category Edit
func AdminEditCategory(c *gin.Context) {
	id := c.Param("id")

	// Fetch category by ID
	var category models.Category
	if err := config.DB.NewSelect().Model(&category).Where("id = ?", id).Scan(c); err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Category Not Found",
		})
		return
	}

	c.HTML(http.StatusOK, "category_edit.html", gin.H{
		"title":    "Edit Category",
		"PageName": "category_edit",
		"data":     category,
	})
}

// Category update
func AdminUpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req dto.CategoryUpdateDTO
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := config.DB.NewSelect().Model(&category).Where("id = ?", id).Scan(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.Name = req.Name
	category.Slug = utils.MakeSlug(req.Name)
	category.Status = req.Status
	category.UpdatedAt = time.Now()

	_, err := config.DB.NewUpdate().Model(&category).Where("id = ?", id).Exec(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/category-list?success=Category+updated+successfully!")
}

// Category Deleted
func AdminDeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Check if category exists
	var category models.Category
	if err := config.DB.NewSelect().Model(&category).Where("id = ?", id).Scan(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Delete category
	_, err := config.DB.NewDelete().Model(&category).Where("id = ?", id).Exec(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func AdminToggleCategoryStatus(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := config.DB.NewSelect().Model(&category).Where("id = ?", id).Scan(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.Status = 1 - category.Status
	_, err := config.DB.NewUpdate().Model(&category).Where("id = ?", id).Exec(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category status updated successfully"})
}
