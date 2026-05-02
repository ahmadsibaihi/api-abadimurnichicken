package expenditure

import (
	"time"
)

type Expenditure struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Category    string    `gorm:"size:100;not null" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	ReceiptURL  string    `gorm:"size:255" json:"receipt_url"`
	Status      string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   *uint     `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	RecipientUserID *uint   `gorm:"column:recipient_user_id" json:"recipient_user_id"`
	RecipientName   string  `gorm:"size:100" json:"recipient_name"`
	DeletedBy   *uint     `json:"deleted_by"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type ExpenditureAuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ExpenditureID uint      `json:"expenditure_id"`
	Action        string    `json:"action"`
	OldData       string    `gorm:"type:text" json:"old_data"`
	NewData       string    `gorm:"type:text" json:"new_data"`
	ChangedBy     uint      `json:"changed_by"`
	ChangedAt     time.Time `json:"changed_at"`
}