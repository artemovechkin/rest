package delivery

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(dataBase *sql.DB) *Service {
	return &Service{
		db: dataBase,
	}
}

var tasks = make(map[int]string)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	Priority    string `json:"priority"`
}

func (s *Service) CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var newTask Task
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	taskID := 0
	err := s.db.QueryRowContext(ctx, "INSERT INTO taskmanager (Title, Description, Status, Priority) VALUES (?, ?, ?, ?) RETURNING ID",
		newTask.Title, newTask.Description, newTask.Status, newTask.Priority).Scan(&taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": taskID})
}

func (s *Service) GetTasks(c *gin.Context) {
	stasusStr := c.Query("status")
	status, err := strconv.ParseBool(stasusStr)
	var taskAll []Task
	var rows *sql.Rows
	if err != nil {
		rows, _ = s.db.Query("SELECT * FROM taskmanager")
	} else {
		rows, _ = s.db.Query("SELECT * FROM taskmanager WHERE Status = (?)", status)
	}
	var task Task
	var flag bool

	for rows.Next() {
		flag = true
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
		if err != nil {
			fmt.Println(err)
			continue
		}
		taskAll = append(taskAll, task)
	}
	if !flag {
		c.JSON(http.StatusOK, gin.H{"error": "Tasks not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, taskAll)
}

func (s *Service) GetTaskByID(c *gin.Context) {
	taskID := getTaskID(c)
	if taskID == -1 {
		return
	}
	var task Task
	s.db.QueryRow("SELECT * from taskmanager WHERE id=?", taskID).Scan(&task.ID,
		&task.Title, &task.Description, &task.Status, &task.Priority)

	c.JSON(http.StatusOK, task)
}

func (s *Service) UpdateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	taskID := getTaskID(c)
	if taskID == -1 {
		return
	}
	var updatedTask Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := s.db.ExecContext(ctx, "UPDATE taskmanager SET Title = ?, Description = ?, Status = ?, Priority = ? WHERE ID = ?",
		updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Priority, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (s *Service) DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	taskID := getTaskID(c)
	res, err := s.db.ExecContext(ctx, "DELETE FROM taskmanager WHERE ID = ?", taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	countRows, _ := res.RowsAffected()
	if countRows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
	cancel()
}

func getTaskID(c *gin.Context) int {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return -1
	}

	return taskID
}

func (s *Service) GetTaskReport(c *gin.Context) {
	query := "SELECT COUNT(*) as total_count, " +
		"SUM(CASE WHEN Status = true THEN 1 ELSE 0 END) as completed_count, " +
		"SUM(CASE WHEN Status = false THEN 1 ELSE 0 END) as pending_count " +
		"FROM taskmanager"

	var summaryReport struct {
		TotalCount          int     `json:"total_count"`
		CompletedCount      int     `json:"completed_count"`
		PendingCount        int     `json:"pending_count"`
		PercentageСompleted float32 `json:"percentage_сompleted"`
	}

	err := s.db.QueryRow(query).Scan(&summaryReport.TotalCount, &summaryReport.CompletedCount, &summaryReport.PendingCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	summaryReport.PercentageСompleted = float32(summaryReport.CompletedCount) / float32(summaryReport.TotalCount) * 100

	c.JSON(http.StatusOK, summaryReport)
}
