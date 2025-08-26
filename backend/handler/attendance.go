package handler

import (
    "time"

    "reonify/clocked/database"
    "reonify/clocked/model"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func uid(c *fiber.Ctx) uint {
    token := c.Locals("user").(*jwt.Token)
    claims := token.Claims.(jwt.MapClaims)
    return uint(int(claims["user_id"].(float64)))
}

type clockInInput struct {
    ClockIn time.Time `json:"clock_in"`
}

type clockOutInput struct {
    ClockOut     time.Time        `json:"clock_out"`
    TaskID       *uint            `json:"task_id"`
    NewTask      *createTaskInput `json:"new_task"`
    MarkComplete bool             `json:"mark_completed"`
}

// POST /api/attendance/clockin
func ClockIn(c *fiber.Ctx) error {
    db := database.DB
    var input clockInInput
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid clock in payload")
    }
    if input.ClockIn.IsZero() {
        input.ClockIn = time.Now()
    }

    rec := model.Attendance{
        UserID:  uid(c),
        ClockIn: input.ClockIn,
    }
    if err := db.Create(&rec).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Couldn't save clock in")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Clocked in", "data": rec})
}

// POST /api/attendance/clockout
func ClockOut(c *fiber.Ctx) error {
    db := database.DB
    var input clockOutInput
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid clock out payload")
    }
    if input.ClockOut.IsZero() {
        input.ClockOut = time.Now()
    }

    // Find the latest open attendance for user
    var att model.Attendance
    if err := db.Where("user_id = ? AND clock_out IS NULL", uid(c)).Order("created_at DESC").First(&att).Error; err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "No open attendance to clock out")
    }

    // Determine / create task
    var taskID *uint = input.TaskID
    if taskID == nil && input.NewTask != nil {
        t := model.Task{
            UserID:      uid(c),
            Name:        input.NewTask.Name,
            Subject:     input.NewTask.Subject,
            Description: input.NewTask.Description,
            HowItsDone:  input.NewTask.HowItsDone,
        }
        if t.Name == "" {
            return fiber.NewError(fiber.StatusBadRequest, "new_task.name is required")
        }
        if err := db.Create(&t).Error; err != nil {
            return fiber.NewError(fiber.StatusInternalServerError, "Couldn't create task")
        }
        taskID = &t.ID

        // Optionally mark complete at clock out time
        if input.MarkComplete {
            now := input.ClockOut
            t.Completed = true
            t.CompletedAt = &now
            _ = db.Save(&t).Error
        }
    } else if taskID != nil && input.MarkComplete {
        var t model.Task
        if err := db.Where("user_id = ?", uid(c)).First(&t, *taskID).Error; err == nil {
            now := input.ClockOut
            t.Completed = true
            t.CompletedAt = &now
            _ = db.Save(&t).Error
        }
    }

    // Update attendance
    att.ClockOut = &input.ClockOut
    att.TaskID = taskID
    if err := db.Save(&att).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Couldn't save clock out")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Clocked out", "data": att})
}

// GET /api/attendance (optional list)
func ListAttendance(c *fiber.Ctx) error {
    db := database.DB
    var atts []model.Attendance
    if err := db.Where("user_id = ?", uid(c)).Order("created_at DESC").Find(&atts).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch attendance")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Attendance", "data": atts})
}

