package dto

type NewsDto struct {
	NewsId       string `json:"id"`
	Notification string `json:"notification"`
	TenantId     string `json:"tenantId"`
}
