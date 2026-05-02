package expenditure

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// "go-serve-pos/database"
    // "go-serve-pos/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	os.MkdirAll("./uploads/expenditures", 0755)
	return &Handler{service: s}
}

// Helper ambil user ID dari context
func getUserID(c *fiber.Ctx) uint {
	val := c.Locals("user_id")
	if val == nil {
		return 1
	}
	id, ok := val.(uint)
	if !ok {
		return 1
	}
	return id
}

// Helper ambil user role dari context
func getUserRole(c *fiber.Ctx) string {
	val := c.Locals("user_role")
	if val == nil {
		return "karyawan" // default
	}
	role, ok := val.(string)
	if !ok {
		return "karyawan"
	}
	return role
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ==================== CREATE with FILE UPLOAD & ROLE-BASED ====================
func (h *Handler) CreateExpenditure(c *fiber.Ctx) error {
	// Ambil form-data
	dateStr := c.FormValue("date")
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)
	category := c.FormValue("category")
	description := c.FormValue("description")
	recipientUserIDStr := c.FormValue("recipient_user_id")
	recipientName := c.FormValue("recipient_name")

	userID := getUserID(c)
	userRole := getUserRole(c)

	// Parse tanggal
	var date time.Time
	if dateStr == "" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah, gunakan YYYY-MM-DD"})
		}
	}

	// Upload file bukti
	var receiptURL string
	file, err := c.FormFile("receipt_image")
	if err == nil {
		os.MkdirAll("./uploads/expenditures", 0755)
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomString(8), ext)
		savePath := filepath.Join("./uploads/expenditures", filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file bukti: " + err.Error()})
		}
		receiptURL = "/uploads/expenditures/" + filename
	} else if userRole == "karyawan" {
		// Karyawan wajib upload bukti
		return c.Status(400).JSON(fiber.Map{"error": "Bukti struk wajib diupload"})
	} else if err != nil && err.Error() != "there is no such file" {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal membaca file: " + err.Error()})
	}

	// Tentukan status dan recipient berdasarkan role
	var status string
	var recipientUserID *uint
	var finalRecipientName string

	if userRole == "admin" || userRole == "superadmin" {
		status = "approved" // langsung approved
		// Admin bisa memilih penerima
		if recipientUserIDStr != "" {
			if id, err := strconv.ParseUint(recipientUserIDStr, 10, 32); err == nil {
				uid := uint(id)
				recipientUserID = &uid
			}
		} else if recipientName != "" {
			finalRecipientName = recipientName
		} else {
			// Default penerima adalah admin itu sendiri
			recipientUserID = &userID
		}
	} else { // karyawan
		status = "pending"
		// Karyawan hanya bisa input untuk dirinya sendiri
		recipientUserID = &userID
		finalRecipientName = ""
	}

	exp := &Expenditure{
		Date:            date,
		Amount:          amount,
		Category:        category,
		Description:     description,
		ReceiptURL:      receiptURL,
		Status:          status,
		CreatedBy:       userID,
		RecipientUserID: recipientUserID,
		RecipientName:   finalRecipientName,
	}

	if err := h.service.CreateExpenditure(exp, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(exp)
}

// ==================== GET ALL ====================
func (h *Handler) GetExpenditures(c *fiber.Ctx) error {
	date := c.Query("date")
	status := c.Query("status")
	list, err := h.service.GetExpenditures(date, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// ==================== GET BY ID ====================
func (h *Handler) GetExpenditureByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	exp, err := h.service.GetExpenditureByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(exp)
}

// ==================== UPDATE with FILE UPLOAD ====================
func (h *Handler) UpdateExpenditure(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	oldExp, err := h.service.GetExpenditureByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	if oldExp.Status != "pending" {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot edit approved or rejected expenditure"})
	}

	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)
	category := c.FormValue("category")
	description := c.FormValue("description")

	receiptURL := oldExp.ReceiptURL
	file, err := c.FormFile("receipt_image")
	if err == nil {
		if oldExp.ReceiptURL != "" {
			oldPath := strings.TrimPrefix(oldExp.ReceiptURL, "/")
			os.Remove(filepath.Join(".", oldPath))
		}
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomString(8), ext)
		savePath := filepath.Join("./uploads/expenditures", filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file bukti: " + err.Error()})
		}
		receiptURL = "/uploads/expenditures/" + filename
	} else if err != nil && err.Error() != "there is no such file" {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal membaca file: " + err.Error()})
	}

	req := UpdateRequest{
		Amount:      amount,
		Category:    category,
		Description: description,
		ReceiptURL:  receiptURL,
	}
	userID := getUserID(c)
	if err := h.service.UpdateExpenditure(uint(id), req, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Updated"})
}

// ==================== SOFT DELETE ====================
func (h *Handler) SoftDeleteExpenditure(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	userID := getUserID(c)
	if err := h.service.DeleteExpenditure(uint(id), userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}

// ==================== APPROVE ====================
func (h *Handler) ApproveExpenditure(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	userID := getUserID(c)
	if err := h.service.ApproveExpenditure(uint(id), userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Approved"})
}

// ==================== DAILY SUMMARY ====================
func (h *Handler) DailySummary(c *fiber.Ctx) error {
	date := c.Query("date")
	total, details, transactions, err := h.service.GetDailySummary(date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"total_expense": total,
		"details":       details,
		"transactions":  transactions,
	})
}