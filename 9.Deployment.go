package main

import (
	_ "database/sql" // add this
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // add this
)

type newStudent struct {
	Student_id       uint64 `json:"student_id" binding:"required"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint64 `json:"student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
}

func postHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent
	c.Bind(&newStudent)
	db.Create(&newStudent)
	c.JSON(http.StatusOK, gin.H{"message": "success create", "data": newStudent})
}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var newStudent []newStudent

	db.Find(&newStudent)
	c.JSON(http.StatusOK, gin.H{"message": "success find all", "data": newStudent})

}

func getHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent
	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success find by id", "data": newStudent})

}

func putHandler(c *gin.Context, db *gorm.DB) {
	var newStudent = newStudent{}

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}

	var reqStudent = newStudent

	c.Bind(&reqStudent)

	db.Model(&newStudent).Where("student_id=?", studentId).Update(reqStudent)

	c.JSON(http.StatusOK, gin.H{
		"message": "success update",
		"data":    reqStudent,
	})

}

func delHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent

	studentId := c.Param("student_id")

	db.Delete(&newStudent, "student_id=?", studentId)

	c.JSON(http.StatusOK, gin.H{
		"message": "success delete",
	})

}

func setupRouter() *gin.Engine {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error load env")
	}

	conn := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})

	r.GET("/student", func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})

	r.GET("/student/:student_id", func(ctx *gin.Context) {
		getHandler(ctx, db)
	})

	r.PUT("/student/:student_id", func(ctx *gin.Context) {
		putHandler(ctx, db)
	})

	r.DELETE("/student/:student_id", func(ctx *gin.Context) {
		delHandler(ctx, db)
	})

	return r

}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&newStudent{})

	data := newStudent{}
	if db.Find(&data).RecordNotFound() {
		fmt.Println("=================== run seeder user ======================")
		seederUser(db)
	}
}

func seederUser(db *gorm.DB) {
	data := newStudent{
		Student_id:       1,
		Student_name:     "Dono",
		Student_age:      20,
		Student_address:  "Jakarta",
		Student_phone_no: "0123456789",
	}

	db.Create(&data)
}

func main() {
	r := setupRouter()

	r.Run(":8080")

}
