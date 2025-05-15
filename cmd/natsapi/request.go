package natsapi

type Request struct {
	ListIp []string `json:"list_ip_addresses"`
	TaskId string   `json:"task_id"`
	Source string   `json:"source"`
}
