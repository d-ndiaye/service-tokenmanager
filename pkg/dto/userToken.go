package dto

type DeviceToken struct {
	Id                       string                   `json:"id"`
	Platform                 int                      `json:"platform"`
	TenantNotificationConfig TenantNotificationConfig `json:"tenantNotificationConfig"`
}
type TenantNotificationConfig struct {
	TenantId string   `json:"tenantId"`
	Category []string `json:"category"`
}
type UserTokenDto struct {
	Username    string      `json:"username"`
	DeviceToken DeviceToken `json:"deviceToken"`
}
