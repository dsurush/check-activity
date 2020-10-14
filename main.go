package main

import (
	"ActivityCheck/models"
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"log"
	"net/http"
	"time"
)

func DoRequest(time int64, status bool, token string) bool {
	client := &http.Client{}
	body := fmt.Sprintf(`"time" : %d
"status" : %t`, time, status)
	fmt.Println(body)
	AuthToken := fmt.Sprintf(`Bearer %s`, token)
	//	s := `{
//  "time" : 960,
//  "status" : false
//}`
	req, err := http.NewRequest(
		"POST", "http://127.0.0.1:8888/api/exit", bytes.NewBuffer([]byte(body)),
	)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", AuthToken)
//	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwiZXhwIjoxNjAyNjM5ODA3LCJsb2dpbiI6InRlc3QiLCJyb2xlIjoidXNlciJ9.HQtfZ1bqtEw-JR4YmAlJRGoHyTxRUCeNrAMIbSqTvfg")
	if err != nil {
		log.Printf("Can't build request\n")
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Can't send Request\n")
		return false
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println(resp.StatusCode)
		return true
	}
	fmt.Println(resp.StatusCode)
	return false
}

func GetActivitiesFromDB(pool *pgxpool.Pool) (Activities []models.Activity, err error){
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `select *from activities where exited = false`)
	if err != nil {
		fmt.Printf("can't read user rows %e", err)
		return
	}
	defer rows.Close()

	for rows.Next(){
		Activity := models.Activity{}
		err := rows.Scan(
			&Activity.ID,
			&Activity.UserId,
			&Activity.Token,
			&Activity.UnixTime,
			&Activity.Status,
			&Activity.WorkTime,
			&Activity.Exited,
				)
		if err != nil {
			fmt.Println("can't scan err is = ", err)
		}
		Activities = append(Activities, Activity)
	}
	if rows.Err() != nil {
		log.Printf("rows err %s", err)
		return nil, rows.Err()
	}
	return
}

func SendRequests(Activities []models.Activity) {
	for _, value := range Activities {
		ok := DoRequest(value.WorkTime, value.Status, value.Token)
		if !ok {
			fmt.Println(`Cant send request to user with user_id =`, value.UserId)
		}
	}
}

func main() {
	pool, err := pgxpool.Connect(context.Background(), `postgres://dsurush:dsurush@localhost:5432/ccd?sslmode=disable`)
	if err != nil {
		log.Printf("Owibka - %e", err)
		log.Fatal("Can't Connection to DB")
	} else {
		fmt.Println("CONNECTION TO DB IS SUCCESS")
	}

	for true {
		activities, err := GetActivitiesFromDB(pool)
		if err != nil {
			fmt.Println("Can't get from DB")
			continue
		}
		SendRequests(activities)
		time.Sleep(time.Minute * 5)
		//ok := DoRequest(960, false, `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwiZXhwIjoxNjAyNjM5ODA3LCJsb2dpbiI6InRlc3QiLCJyb2xlIjoidXNlciJ9.HQtfZ1bqtEw-JR4YmAlJRGoHyTxRUCeNrAMIbSqTvfg`)
		//fmt.Println(ok)
	}
}
