package services

import (
	"kasir-api-irfan/models"
	"kasir-api-irfan/repositories"
)

type ReportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) *ReportService {
	return &ReportService{
		repo: repo,
	}
}

func (s *ReportService) GetSalesSummary(startDate, endDate string) (*models.SalesSummary, error) {
	return s.repo.GetSalesSummary(startDate, endDate)
}

func (s *ReportService) GetSalesSummaryToday() (*models.SalesSummary, error) {
	return s.repo.GetSalesSummaryToday()
}
