package enum

type SessionOwner = string

const (
	Buyer      SessionOwner = "buyer"
	Supplier   SessionOwner = "supplier"
	Admin      SessionOwner = "admin"
	SuperAdmin SessionOwner = "super_admin"
)
