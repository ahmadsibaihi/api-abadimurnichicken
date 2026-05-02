package expenditure

import (
	"time"
	"gorm.io/gorm"
)

type Repository interface {
	Create(exp *Expenditure) error
	GetAll(filter map[string]interface{}) ([]Expenditure, error)
	GetByID(id uint) (*Expenditure, error)
	Update(id uint, updates map[string]interface{}) error
	SoftDelete(id uint, deletedBy uint) error
	CreateAuditLog(log *ExpenditureAuditLog) error
	GetDailySummary(date time.Time) (total float64, details map[string]float64, err error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(exp *Expenditure) error {
	return r.db.Create(exp).Error
}

func (r *repository) GetAll(filter map[string]interface{}) ([]Expenditure, error) {
	var list []Expenditure
	query := r.db.Where("deleted_at IS NULL")
	if date, ok := filter["date"]; ok {
		query = query.Where("date = ?", date)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *repository) GetByID(id uint) (*Expenditure, error) {
	var exp Expenditure
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&exp).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *repository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&Expenditure{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) SoftDelete(id uint, deletedBy uint) error {
	now := time.Now()
	return r.db.Model(&Expenditure{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": deletedBy,
			"deleted_at": now,
		}).Error
}

func (r *repository) CreateAuditLog(log *ExpenditureAuditLog) error {
	return r.db.Create(log).Error
}

func (r *repository) GetDailySummary(date time.Time) (total float64, details map[string]float64, err error) {
	type Result struct {
		Category string
		Total    float64
	}
	var results []Result
	err = r.db.Model(&Expenditure{}).
		Select("category, SUM(amount) as total").
		Where("date = ? AND status = ? AND deleted_at IS NULL", date, "approved").
		Group("category").Scan(&results).Error
	if err != nil {
		return 0, nil, err
	}
	details = make(map[string]float64)
	for _, r := range results {
		details[r.Category] = r.Total
		total += r.Total
	}
	return total, details, nil
}