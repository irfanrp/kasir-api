package handlers

import (
	"encoding/json"
	"kasir-api-irfan/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetSalesReport godoc
// @Summary Get sales report
// @Description Mendapatkan laporan penjualan berdasarkan periode (default: hari ini)
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.SalesSummary
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/report [get]
func (h *ReportHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	summary, err := h.service.GetSalesSummary(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
