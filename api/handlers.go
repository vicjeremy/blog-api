package api

// The application business logic
import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// migrate the schema
	if err := DB.AutoMigrate(&Blog{}); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func CreateBlog(c *gin.Context) {
	var blog Blog

	// bind the request body
	if err := c.ShouldBindJSON(&blog); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if err := DB.Create(&blog).Error; err != nil {
		log.Printf("Failed to create blog: %v", err)
		ResponseJSON(c, http.StatusInternalServerError, "Failed to create blog", nil)
		return
	}

	ResponseJSON(c, http.StatusCreated, "Blog created successfully", blog)
}

func GetBlogs(c *gin.Context) {
	var blogs []Blog
	term := c.Query("term")
	result := DB
	if term != "" {
		searchTerm := "%" + strings.ToLower(term) + "%"
		result = DB.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(category) LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	if err := result.Find(&blogs).Error; err != nil {
		log.Printf("Database error in GetBlogs: %v", err)
		if strings.Contains(strings.ToLower(err.Error()), "column") {
			ResponseJSON(c, http.StatusInternalServerError, "Invalid search field", nil)
		} else if strings.Contains(strings.ToLower(err.Error()), "syntax") {
			ResponseJSON(c, http.StatusInternalServerError, "Invalid search syntax", nil)
		} else {
			ResponseJSON(c, http.StatusInternalServerError, "Failed to retrieve blogs", nil)
		}
		return
	}

	if len(blogs) == 0 {
		ResponseJSON(c, http.StatusNotFound, "No blogs found matching your criteria", []Blog{})
		return
	}
	ResponseJSON(c, http.StatusOK, "Blogs retrieved successfully", blogs)
}

func GetBlog(c *gin.Context) {
	var blog Blog
	if err := DB.First(&blog, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Blog not found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Blog retrieved successfully", blog)
}

func UpdateBlog(c *gin.Context) {
	var blog Blog
	if err := DB.First(&blog, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Blog not found", nil)
		return
	}

	var updateData Blog
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	// Only update allowed fields
	updates := map[string]interface{}{
		"title":    updateData.Title,
		"content":  updateData.Content,
		"category": updateData.Category,
		"tags":     updateData.Tags,
	}

	if err := DB.Model(&blog).Updates(updates).Error; err != nil {
		log.Printf("Failed to update blog: %v", err)
		ResponseJSON(c, http.StatusInternalServerError, "Failed to update blog", nil)
		return
	}

	// Fetch the updated blog to return
	DB.First(&blog, c.Param("id"))
	ResponseJSON(c, http.StatusOK, "Blog updated successfully", blog)
}

func DeleteBlog(c *gin.Context) {
	var blog Blog
	result := DB.First(&blog, c.Param("id"))
	if result.Error != nil {
		ResponseJSON(c, http.StatusNotFound, "Blog not found", nil)
		return
	}

	if err := DB.Delete(&blog).Error; err != nil {
		log.Printf("Failed to delete blog: %v", err)
		ResponseJSON(c, http.StatusInternalServerError, "Failed to delete blog", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Blog deleted successfully", nil)
}
