package models

type ProductInfo struct {
	Nama        string `json:"nama"`
	QtyTerjual  int    `json:"qty_terjual"`
}

type SalesSummary struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalTransaksi int            `json:"total_transaksi"`
	ProdukTerlaris ProductInfo    `json:"produk_terlaris"`
}
