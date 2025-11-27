package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/howeyc/gopass"
)

func GetCookie() string {
	var username string
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	fmt.Print("Password: ")
	passwordBytes, _ := gopass.GetPasswdMasked()
	password := string(passwordBytes)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("http://kjyy.ccnu.edu.cn/clientweb/xcus/ic2/Default.aspx"),
		chromedp.Sleep(2*time.Second),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Click(`input[type="submit"]`, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
	)

	if err != nil {
		log.Fatal(err)
	}

	var cookies []*network.Cookie
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		allCookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}
		cookies = allCookies
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cookies {
		return c.Value
	}

	return ""
}
