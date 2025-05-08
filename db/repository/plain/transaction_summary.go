package plain

// TransactionSummary represents grouped data by category
type TransactionSummary struct {
	CategoryName string
	TotalAmount  float64
	Count        int64
}
