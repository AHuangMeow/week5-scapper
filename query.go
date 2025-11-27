package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type QueryResponse struct {
	Msg string `json:"msg"`
}

func GetSitInfo(cookieValue string) (string, string, string) {
	base := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/center.aspx?act=get_History_resv&strat=90&StatFlag=New&_=1764240101437"

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

	var queryResponse QueryResponse
	err = json.Unmarshal(body, &queryResponse)

	if strings.Contains(queryResponse.Msg, "没有数据") {
		return "", "", ""
	} else {
		rsvIdRe := regexp.MustCompile(`rsvId='(\d+)'`)
		rsvIdMatches := rsvIdRe.FindStringSubmatch(queryResponse.Msg)

		timeRe := regexp.MustCompile(`text-primary'>(\d{2}-\d{2} \d{2}:\d{2})</span>`)
		timeMatches := timeRe.FindAllStringSubmatch(queryResponse.Msg, -1)

		return rsvIdMatches[1], timeMatches[0][1], timeMatches[1][1]
	}
}

func Query() {
	rsvID, start, end := GetSitInfo(GetCookie())
	if rsvID == "" {
		fmt.Println("No reserved sit")
	} else {
		fmt.Println("Reserved sit:")
		fmt.Println("RsvID: ", rsvID)
		fmt.Println("Start time: ", start)
		fmt.Println("End time: ", end)
	}
}
