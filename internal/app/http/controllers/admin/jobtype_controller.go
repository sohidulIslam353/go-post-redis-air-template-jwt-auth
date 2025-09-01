package controllers

import (
	"context"
	"gin-app/config"
	"gin-app/internal/dto"
	"gin-app/internal/models"
	"gin-app/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Job Type function
func AdminJobTypeList(c *gin.Context) {
	successMsg := c.Query("success")
	// Filters
	search := c.Query("search")
	status := c.Query("status")
	createdAt := c.Query("created_at")

	// Cursor pagination
	lastID, limit := utils.GetCursorPagination(c)

	// Base query
	query := config.DB.NewSelect().Model((*models.JobType)(nil))

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
	var jobs []models.JobType
	if err := query.Scan(c, &jobs); err != nil {
		c.HTML(http.StatusInternalServerError, "job_type_list.html", gin.H{
			"title": "Job Type List",
			"error": "Failed to fetch job types: " + err.Error(),
			"data":  []models.JobType{},
		})
		return
	}

	// Next cursor
	var nextCursor int64
	if len(jobs) > 0 {
		nextCursor = jobs[len(jobs)-1].ID
	}

	// Base query for count
	countQuery := config.DB.NewSelect().Model((*models.JobType)(nil))
	// Count total items
	totalCount, err := countQuery.Count(c)
	if err != nil {
		totalCount = 0
	}
	// Render
	c.HTML(http.StatusOK, "job_type_list.html", gin.H{
		"title":      "Job Type List",
		"PageName":   "job_type_list",
		"data":       jobs,
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

// Job type create page
func AdminJobTypeCreate(c *gin.Context) {

	c.HTML(http.StatusOK, "job_type_create.html", gin.H{
		"title": "Job Type Create",
	})
}

// Jobtype store
func AdminJobTypeStore(c *gin.Context) {
	var input dto.JobTypeStoreDTO

	//  Bind + Validate
	if valid, errs := utils.ValidateStruct(c, &input); !valid {
		c.HTML(http.StatusBadRequest, "job_type_create.html", gin.H{
			"title":  "Create Job Type",
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
	job := models.JobType{
		Name:      input.Name,
		Slug:      slug,
		Status:    input.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Database value insert
	if _, err := config.DB.NewInsert().Model(&job).Exec(c); err != nil {
		c.HTML(http.StatusInternalServerError, "job_type_create.html", gin.H{
			"title":  "Create Job Type",
			"errors": map[string]string{"DB": "Failed to create job type: " + err.Error()},
			"data":   input,
		})
		return
	}

	// Success response → empty form + success msg
	c.HTML(http.StatusOK, "job_type_create.html", gin.H{
		"title":   "Create Job Type",
		"success": "Job Type created successfully!",
		"errors":  map[string]string{},
		"data":    dto.JobTypeStoreDTO{}, // clear form
	})
}

// Status job type update
func AdminToggleJobTypeStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	_, err := config.DB.NewUpdate().
		Model(&models.JobType{}).
		Set("status = ?", req.Status).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "status": req.Status})
}

// Delete job type
func AdminDeleteJobType(c *gin.Context) {
	id := c.Param("id")

	// Optional: Protect admin-only
	_, err := config.DB.NewDelete().
		Model(&models.JobType{}).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job type deleted"})
}

// Update Page
func AdminEditJobType(c *gin.Context) {
	id := c.Param("id")

	var job models.JobType
	if err := config.DB.NewSelect().Model(&job).Where("id = ?", id).Scan(c); err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "Not Found"})
		return
	}

	c.HTML(http.StatusOK, "job_type_edit.html", gin.H{
		"title": "Edit Job Type",
		"data":  job,
	})
}

// Update jobtype
func AdminUpdateJobType(c *gin.Context) {
	id := c.Param("id")

	var req dto.JobTypeStoreDTO
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var job models.JobType
	if err := config.DB.NewSelect().Model(&job).Where("id = ?", id).Scan(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job type not found"})
		return
	}

	job.Name = req.Name
	job.Slug = utils.MakeSlug(req.Name)
	job.Status = req.Status
	job.UpdatedAt = time.Now()

	_, err := config.DB.NewUpdate().Model(&job).Where("id = ?", id).Exec(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job type"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/job-type-list?success=Job+Type+updated+successfully!")
}
