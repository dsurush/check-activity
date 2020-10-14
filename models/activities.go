package models

type Activity struct {
	ID int64 `json:"id"`
	UserId int64 `json:"user_id"`
	Token string `json:"token"`
	UnixTime int64 `json:"unix_time"`
	Status bool `json:"status"`
	WorkTime int64 `json:"work_time"`
	Exited bool `json:"exited"`
}