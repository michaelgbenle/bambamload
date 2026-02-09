package models

type DashboardStats struct {
	TotalSuppliers    int64
	VerifiedSuppliers int64
	ActiveSuppliers   int64

	TotalBuyers  int64
	ActiveBuyers int64

	TotalProducts    int64
	ApprovedProducts int64
	PendingProducts  int64

	TotalOrders     int64
	CompletedOrders int64
	InTransitOrders int64
}
