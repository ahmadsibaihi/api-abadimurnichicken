package expenditure

import (
	"encoding/json"
	"errors"
	"time"
)

type Service interface {
	CreateExpenditure(exp *Expenditure, userID uint) error
	GetExpenditures(date, status string) ([]Expenditure, error)
	GetExpenditureByID(id uint) (*Expenditure, error)
	UpdateExpenditure(id uint, req UpdateRequest, userID uint) error
	DeleteExpenditure(id uint, userID uint) error
	ApproveExpenditure(id uint, userID uint) error
	GetDailySummary(dateStr string) (Total float64, Details map[string]float64, Transactions []Expenditure, err error)
}

type UpdateRequest struct {
	Amount      float64
	Category    string
	Description string
	ReceiptURL  string
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateExpenditure(exp *Expenditure, userID uint) error {
	exp.CreatedBy = userID
	exp.Status = "pending"
	if exp.Date.IsZero() {
		exp.Date = time.Now()
	}
	err := s.repo.Create(exp)
	if err != nil {
		return err
	}
	jsonData, _ := json.Marshal(exp)
	audit := &ExpenditureAuditLog{
		ExpenditureID: exp.ID,
		Action:        "INSERT",
		NewData:       string(jsonData),
		ChangedBy:     userID,
		ChangedAt:     time.Now(),
	}
	_ = s.repo.CreateAuditLog(audit)
	return nil
}

func (s *service) GetExpenditures(date, status string) ([]Expenditure, error) {
	filter := map[string]interface{}{}
	if date != "" {
		filter["date"] = date
	}
	if status != "" {
		filter["status"] = status
	}
	return s.repo.GetAll(filter)
}

func (s *service) GetExpenditureByID(id uint) (*Expenditure, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateExpenditure(id uint, req UpdateRequest, userID uint) error {
	exp, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if exp.Status != "pending" {
		return errors.New("cannot edit approved or rejected expenditure")
	}
	oldJSON, _ := json.Marshal(exp)
	updates := map[string]interface{}{
		"amount":      req.Amount,
		"category":    req.Category,
		"description": req.Description,
		"receipt_url": req.ReceiptURL,
		"updated_by":  userID,
		"updated_at":  time.Now(),
	}
	err = s.repo.Update(id, updates)
	if err != nil {
		return err
	}
	newExp, _ := s.repo.GetByID(id)
	newJSON, _ := json.Marshal(newExp)
	audit := &ExpenditureAuditLog{
		ExpenditureID: id,
		Action:        "UPDATE",
		OldData:       string(oldJSON),
		NewData:       string(newJSON),
		ChangedBy:     userID,
		ChangedAt:     time.Now(),
	}
	_ = s.repo.CreateAuditLog(audit)
	return nil
}

func (s *service) DeleteExpenditure(id uint, userID uint) error {
	exp, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if exp.Status == "approved" {
		return errors.New("cannot delete approved expenditure")
	}
	oldJSON, _ := json.Marshal(exp)
	err = s.repo.SoftDelete(id, userID)
	if err != nil {
		return err
	}
	audit := &ExpenditureAuditLog{
		ExpenditureID: id,
		Action:        "DELETE",
		OldData:       string(oldJSON),
		ChangedBy:     userID,
		ChangedAt:     time.Now(),
	}
	_ = s.repo.CreateAuditLog(audit)
	return nil
}

func (s *service) ApproveExpenditure(id uint, userID uint) error {
	exp, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if exp.Status != "pending" {
		return errors.New("expenditure already processed")
	}
	oldJSON, _ := json.Marshal(exp)
	err = s.repo.Update(id, map[string]interface{}{
		"status":     "approved",
		"updated_by": userID,
		"updated_at": time.Now(),
	})
	if err != nil {
		return err
	}
	newExp, _ := s.repo.GetByID(id)
	newJSON, _ := json.Marshal(newExp)
	audit := &ExpenditureAuditLog{
		ExpenditureID: id,
		Action:        "APPROVE",
		OldData:       string(oldJSON),
		NewData:       string(newJSON),
		ChangedBy:     userID,
		ChangedAt:     time.Now(),
	}
	_ = s.repo.CreateAuditLog(audit)
	return nil
}

func (s *service) GetDailySummary(dateStr string) (total float64, details map[string]float64, transactions []Expenditure, err error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now()
	}
	total, details, err = s.repo.GetDailySummary(date)
	if err != nil {
		return 0, nil, nil, err
	}
	list, _ := s.repo.GetAll(map[string]interface{}{
		"date":   date.Format("2006-01-02"),
		"status": "approved",
	})
	transactions = list
	return total, details, transactions, nil
}