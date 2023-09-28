package main

import (
	"click-supporter/mouse"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msgf("commandline has arguments length: %d", len(os.Args))

	for i := 0; i < len(os.Args); i++ {
		fmt.Printf("%d: %s\n", i, os.Args[i])
	}

	if len(os.Args) < 5 {
		x, y := robotgo.GetMousePos()
		log.Fatal().Msgf("引数が足りません。引数は4つ、[回数, タイマー(秒数), X座標, Y座標, op:ダブルクリック0/1, op:開始時刻'MM/DD hh:mm']です。現在のマウス座標は[X: %d, Y: %d]", x, y)
	}

	// マウスコントロールに必要な情報を渡す
	c := mouse.New()

	c.Count, _ = strconv.Atoi(os.Args[1])
	c.TimerSecond, _ = strconv.Atoi(os.Args[2])
	c.X, _ = strconv.Atoi(os.Args[3])
	c.Y, _ = strconv.Atoi(os.Args[4])

	if len(os.Args) > 5 { // [オプション]開始時間の引数
		// ダブルクリック
		isDoubleClick, _ := strconv.Atoi(os.Args[5])
		if isDoubleClick == 1 {
			c.IsDoubleClick = true
			fmt.Println("double click!!")
		}

		if len(os.Args) > 6 { // [オプション]開始時間の引数
			// time.Parse()ではUTC指定となるため、+9hの誤差をLocationで修正する
			layout := "01/02 15:04"
			year := time.Now().Year()
			c.StartAt, _ = time.Parse(layout, os.Args[6])
			c.StartAt = c.StartAt.AddDate(year, 0, 0)
			jst, _ := time.LoadLocation("Asia/Tokyo")
			c.StartAt, _ = time.ParseInLocation("2006-01-02 15:04:05 +0000 UTC", c.StartAt.String(), jst)
		}
	}

	log.Info().Msgf(
		"クリック回数:%d回, タイマー: %d秒, mouse X/Y: %d/%d, 開始時刻: %d:%d:%d",
		c.Count,
		c.TimerSecond,
		c.X,
		c.Y,
		c.StartAt.Hour(),
		c.StartAt.Minute(),
		c.StartAt.Second())

	// 指定回数実行
	done := c.Execute()
	if done == nil { // success
		log.Fatal().Msgf("the click supporter ended for no reason, %s", done.Error())
	} else {
		log.Fatal().Msgf("the click supporter is done, %s", done.Error())
	}

}
