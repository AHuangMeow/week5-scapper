package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ReserveResponse struct {
	Msg string `json:"msg"`
}

func ParseTime(str string) string {
	tail := str[len(str)-7:]
	hhmm := strings.Split(tail, "%3A")
	hh, _ := strconv.Atoi(hhmm[0])
	mm, _ := strconv.Atoi(hhmm[1])

	return strconv.Itoa(hh*100 + mm)
}

func Reserve() {
	cookieValue := GetCookie()

	var devID, start, end string
	fmt.Print("Device ID: ")
	fmt.Scanln(&devID)
	fmt.Print("Start time(1970-1-1,00:00): ")
	fmt.Scanln(&start)
	start = strings.ReplaceAll(start, ":", "%3A")
	start = strings.ReplaceAll(start, ",", "+")
	fmt.Print("End time(1970-1-1,00:00): ")
	fmt.Scanln(&end)
	end = strings.ReplaceAll(end, ":", "%3A")
	end = strings.ReplaceAll(end, ",", "+")

	base := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dialogid=&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&up_file=&memo=&act=set_resv&_=1764240101435"
	base += "&dev_id=" + devID
	base += "&start=" + start
	base += "&end=" + end
	base += "&start_time=" + ParseTime(start)
	base += "&end_time=" + ParseTime(end)

	req, err := http.NewRequest("GET", base, nil)
	if err != nil {
		log.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: cookieValue,
	}
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var reserveResponse ReserveResponse
	err = json.Unmarshal(body, &reserveResponse)

	fmt.Println(reserveResponse.Msg)
}
