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

func AdminSubCategoryList(c *gin.Context) {
	successMsg := c.Query("success")
	// Filters
	search := c.Query("search")
	categoryID := c.Query("category_id")
	createdAt := c.Query("created_at")

	// Cursor pagination
	lastID, limit := utils.GetCursorPagination(c)

	// Base query
	query := config.DB.NewSelect().
		Model((*models.Subcategory)(nil)).
		Relation("Category") // join

	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
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
	var data []models.Subcategory
	if err := query.Scan(c, &data); err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_list.html", gin.H{
			"title": "Subcategory List",
			"error": "Failed to fetch subcategories: " + err.Error(),
			"data":  []models.Subcategory{},
		})
		return
	}

	// Next cursor
	var nextCursor int64
	if len(data) > 0 {
		nextCursor = data[len(data)-1].ID
	}

	// Base query for count
	countQuery := config.DB.NewSelect().Model((*models.Subcategory)(nil))
	// Count total items
	totalCount, err := countQuery.Count(c)
	if err != nil {
		totalCount = 0
	}

	var categories []models.Category

	err = config.DB.NewSelect().
		Model(&categories).
		Where("status = ?", 1).
		Scan(c.Request.Context())

	if err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_list.html", gin.H{
			"title": "Subcategory List",
			"error": "Failed to fetch categories: " + err.Error(),
		})
		return
	}
	// Render
	c.HTML(http.StatusOK, "subcategory_list.html", gin.H{
		"title":      "Subcategory List",
		"PageName":   "subcategory_list",
		"data":       data,
		"categories": categories,
		"nextCursor": nextCursor,
		"limit":      limit,
		"total":      totalCount,
		"filters": gin.H{
			"search":      search,
			"category_id": categoryID,
			"created_at":  createdAt,
		},
		"success": successMsg, // pass to template
	})
}

// Subcategory create
func AdminSubCategoryCreate(c *gin.Context) {
	var categories []models.Category

	err := config.DB.NewSelect().
		Model(&categories).
		Where("status = ?", 1).
		Scan(c.Request.Context())

	if err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_create.html", gin.H{
			"title": "Subcategory List",
			"error": "Failed to fetch categories: " + err.Error(),
		})
		return
	}
	// ✅ success query param ধরুন
	successMsg := c.Query("success")
	c.HTML(http.StatusOK, "subcategory_create.html", gin.H{
		"title":      "Subcategory Create",
		"categories": categories,
		"success":    successMsg,
	})
}

// sub category store
func AdminSubCategoryStore(c *gin.Context) {
	var input dto.SubcategoryStoreDTO

	//  Bind + Validate
	if valid, errs := utils.ValidateStruct(c, &input); !valid {
		// categories আনতে হবে, কারণ error হলে আবার form show করতে হবে
		var categories []models.Category
		err := config.DB.NewSelect().
			Model(&categories).
			Where("status = ?", 1).
			Scan(c.Request.Context())

		if err != nil {
			c.HTML(http.StatusInternalServerError, "subcategory_create.html", gin.H{
				"title": "Create Subcategory",
				"error": "Failed to fetch categories: " + err.Error(),
			})
			return
		}

		c.HTML(http.StatusBadRequest, "subcategory_create.html", gin.H{
			"title":      "Create Subcategory",
			"errors":     errs,
			"data":       input,
			"categories": categories,
		})
		return
	}

	// Slug create
	slug := utils.MakeSlug(input.Name)
	if slug == "" {
		slug = "not-available"
	}

	// Model create
	subcategory := models.Subcategory{
		Name:       input.Name,
		CategoryID: input.CategoryID,
		Slug:       slug,
		Status:     input.Status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Database insert
	if _, err := config.DB.NewInsert().Model(&subcategory).Exec(c); err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_create.html", gin.H{
			"title":  "Create Subcategory",
			"errors": map[string]string{"DB": "Failed to create subcategory: " + err.Error()},
			"data":   input,
		})
		return
	}

	// ✅ Redirect করে আবার Create form call করুন
	c.Redirect(http.StatusSeeOther, "/admin/subcategory-create?success=Subcategory+created+successfully!")
}

// edit
func AdminEditSubCategory(c *gin.Context) {
	id := c.Param("id")
	// Fetch subcategory by ID
	var subcategory models.Subcategory
	if err := config.DB.NewSelect().Model(&subcategory).Where("id = ?", id).Scan(c); err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Subcategory Not Found",
		})
		return
	}

	// fetch category also
	var categories []models.Category

	err := config.DB.NewSelect().
		Model(&categories).
		Where("status = ?", 1).
		Scan(c.Request.Context())

	if err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_create.html", gin.H{
			"title": "Subcategory List",
			"error": "Failed to fetch categories: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "subcategory_edit.html", gin.H{
		"title":      "Edit Subcategory",
		"PageName":   "subcategory_edit",
		"data":       subcategory,
		"categories": categories,
	})
}

// update subcategory
func AdminUpdateSubCategory(c *gin.Context) {
	var input dto.SubcategoryUpdateDTO

	//  Bind + Validate
	if valid, errs := utils.ValidateStruct(c, &input); !valid {
		// categories আনতে হবে, কারণ error হলে আবার form show করতে হবে
		var categories []models.Category
		err := config.DB.NewSelect().
			Model(&categories).
			Where("status = ?", 1).
			Scan(c.Request.Context())

		if err != nil {
			c.HTML(http.StatusInternalServerError, "subcategory_edit.html", gin.H{
				"title": "Edit Subcategory",
				"error": "Failed to fetch categories: " + err.Error(),
			})
			return
		}

		c.HTML(http.StatusBadRequest, "subcategory_edit.html", gin.H{
			"title":      "Edit Subcategory",
			"errors":     errs,
			"data":       input,
			"categories": categories,
		})
		return
	}

	// Fetch existing subcategory
	var subcategory models.Subcategory
	if err := config.DB.NewSelect().Model(&subcategory).Where("id = ?", c.Param("id")).Scan(c.Request.Context()); err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Subcategory Not Found",
		})
		return
	}

	// Update fields
	subcategory.Name = input.Name
	subcategory.Slug = utils.MakeSlug(input.Name)
	subcategory.CategoryID = input.CategoryID
	subcategory.Status = input.Status

	// Database update
	if _, err := config.DB.NewUpdate().Model(&subcategory).Where("id = ?", subcategory.ID).Exec(c); err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_edit.html", gin.H{
			"title":  "Edit Subcategory",
			"errors": map[string]string{"DB": "Failed to update subcategory: " + err.Error()},
			"data":   input,
		})
		return
	}

	// ✅ Redirect to list
	c.Redirect(http.StatusSeeOther, "/admin/subcategory-list?success=Subcategory+updated+successfully!")
}

// delete
func AdminDeleteSubCategory(c *gin.Context) {
	id := c.Param("id")

	// Database delete
	if _, err := config.DB.NewDelete().Model(&models.Subcategory{}).Where("id = ?", id).Exec(c); err != nil {
		c.HTML(http.StatusInternalServerError, "subcategory_list.html", gin.H{
			"title":  "Subcategory List",
			"errors": map[string]string{"DB": "Failed to delete subcategory: " + err.Error()},
		})
		return
	}

	// ✅ Redirect to list
	c.JSON(http.StatusOK, gin.H{"message": "Subcategory deleted successfully"})
}

// toggle status changed
func AdminToggleSubCategoryStatus(c *gin.Context) {
	id := c.Param("id")

	var subcategory models.Subcategory
	if err := config.DB.NewSelect().Model(&subcategory).Where("id = ?", id).Scan(c.Request.Context()); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
		return
	}

	subcategory.Status = 1 - subcategory.Status
	_, err := config.DB.NewUpdate().Model(&subcategory).Where("id = ?", id).Exec(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subcategory status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subcategory status updated successfully"})
}
