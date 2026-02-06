package repositories

import (
	"database/sql"
	"kasir-api-irfan/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return ReportRepository{db: db}
}

func (repo *ReportRepository) GetSalesSummary(startDate, endDate string) (*models.SalesSummary, error) {
	// If no dates provided, use today
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Get total revenue and transaction count
	var totalRevenue int
	var totalTransaksi int
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get top selling product
	var productName string
	var qtyTerjual int
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity)
		FROM transaction_details td
		JOIN product p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.id, p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &qtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &models.SalesSummary{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: models.ProductInfo{
			Nama:       productName,
			QtyTerjual: qtyTerjual,
		},
	}, nil
}

func (repo *ReportRepository) GetSalesSummaryToday() (*models.SalesSummary, error) {
	today := time.Now().Format("2006-01-02")
	return repo.GetSalesSummary(today, today)
}
