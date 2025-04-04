package models

type BillItem struct {
	Amount  string `json:"amount"`
	DueDate string `json:"due_date"`
	Status  string `json:"status"`
}

type BillResponse struct {
	Bills []BillItem `json:"bills"`
}

type Bill struct {
	UserID       string     `json:"user_id"`
	ProviderID   string     `json:"provider_id"`
	ProviderName string     `json:"provider_name"`
	Items        []BillItem `json:"items"`
}
