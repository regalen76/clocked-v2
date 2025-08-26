package handler

import (
    "time"

    "reonify/clocked/database"
    "reonify/clocked/model"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

type createTaskInput struct {
    Name        string `json:"name"`
    Subject     string `json:"subject"`
    Description string `json:"description"`
    HowItsDone  string `json:"how_its_done"`
}

type updateTaskInput struct {
    Completed   *bool      `json:"completed"`
    CompletedAt *time.Time `json:"completed_at"`
}

func getUserIDFromCtx(c *fiber.Ctx) (uint, error) {
    token := c.Locals("user").(*jwt.Token)
    claims := token.Claims.(jwt.MapClaims)
    uid := uint(int(claims["user_id"].(float64)))
    return uid, nil
}

// List tasks for the current user
func GetTasks(c *fiber.Ctx) error {
    db := database.DB
    uid, _ := getUserIDFromCtx(c)

    var tasks []model.Task
    if err := db.Where("user_id = ?", uid).Order("created_at DESC").Find(&tasks).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch tasks")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Tasks", "data": tasks})
}

// Create a new task
func CreateTask(c *fiber.Ctx) error {
    db := database.DB
    uid, _ := getUserIDFromCtx(c)

    var input createTaskInput
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid task payload")
    }
    if input.Name == "" {
        return fiber.NewError(fiber.StatusBadRequest, "Task name is required")
    }

    task := model.Task{
        UserID:      uid,
        Name:        input.Name,
        Subject:     input.Subject,
        Description: input.Description,
        HowItsDone:  input.HowItsDone,
    }

    if err := db.Create(&task).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Couldn't create task")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Created task", "data": task})
}

// Update a task (e.g., mark complete)
func UpdateTask(c *fiber.Ctx) error {
    db := database.DB
    uid, _ := getUserIDFromCtx(c)
    id := c.Params("id")

    var task model.Task
    if err := db.Where("user_id = ?", uid).First(&task, id).Error; err != nil {
        return fiber.NewError(fiber.StatusNotFound, "Task not found")
    }

    var input updateTaskInput
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid task payload")
    }

    if input.Completed != nil {
        task.Completed = *input.Completed
        if task.Completed {
            // set CompletedAt if provided, else now
            if input.CompletedAt != nil {
                task.CompletedAt = input.CompletedAt
            } else {
                now := time.Now()
                task.CompletedAt = &now
            }
        } else {
            task.CompletedAt = nil
        }
    }

    if err := db.Save(&task).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Couldn't update task")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Updated task", "data": task})
}

// Tasks by a given day (created or completed on date)
func GetTasksByDay(c *fiber.Ctx) error {
    db := database.DB
    uid, _ := getUserIDFromCtx(c)
    dateStr := c.Query("date") // YYYY-MM-DD
    if dateStr == "" {
        return fiber.NewError(fiber.StatusBadRequest, "date query param required")
    }
    // build start/end of day
    t, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "invalid date format (YYYY-MM-DD)")
    }
    start := t
    end := t.Add(24 * time.Hour)

    var tasks []model.Task
    if err := db.Where("user_id = ? AND ((created_at >= ? AND created_at < ?) OR (completed_at IS NOT NULL AND completed_at >= ? AND completed_at < ?))",
        uid, start, end, start, end).Order("created_at DESC").Find(&tasks).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch tasks by date")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Tasks for day", "data": tasks})
}

