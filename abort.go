package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AbortResponse struct {
	Msg string `json:"msg"`
}

func Abort() {
	cookieValue := GetCookie()
	rsvID, _, _ := GetSitInfo(cookieValue)

	base := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?act=del_resv&_=1764240101438&id=" + rsvID

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

	var abortResponse AbortResponse
	err = json.Unmarshal(body, &abortResponse)

	fmt.Println(abortResponse.Msg)
}
