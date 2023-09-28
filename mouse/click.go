package mouse

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	Count         int
	TimerSecond   int
	X             int
	Y             int
	IsDoubleClick bool
	StartAt       time.Time
}

func New() *Controller {
	return &Controller{
		Count:         1,
		TimerSecond:   60,
		X:             int(1920 / 2),
		Y:             int(1080 / 2),
		IsDoubleClick: false,
		StartAt:       time.Now(),
	}
}

func (p *Controller) Execute() error {
	count := p.Count
	msg := "start!"

	sub := p.StartAt.Sub(time.Now().UTC())
	startTicker := time.NewTimer(sub)
	defer startTicker.Stop()

	if sub > time.Duration(0) {
		log.Info().Msgf("timer start, %f秒後実行します", sub.Seconds())
		<-startTicker.C

		// 実行
		robotgo.Move(p.X, p.Y)
		robotgo.Click()
		log.Info().Msg("finished click for starttimer")

		count -= 1
	}

	if count < 1 {
		msg = "the mouse controller timer was done!"
		log.Info().Msg(msg)

		return fmt.Errorf("%s", msg)
	}

	// Timerで実行を終えてからTickerを開始する
	ticker := time.NewTicker(time.Duration(p.TimerSecond) * time.Second)
	defer ticker.Stop()

	log.Info().Msgf("ticker start every %ds", p.TimerSecond)
	// 実行回数
	counter := 0
	for {
		<-ticker.C

		// 実行
		robotgo.Move(p.X, p.Y)
		robotgo.Click()

		if p.IsDoubleClick {
			// Optional: ウインドウのアクティブ化を含む場合やダブルクリックの場合
			time.Sleep(500 * time.Millisecond)
			robotgo.Click()
		}

		counter++
		log.Info().Msgf("%d:%d Clicked: %d\n", time.Now().Minute(), time.Now().Second(), counter)

		count -= 1
		if count < 1 {
			msg += "all times completed. "
			break
		}
	}
	log.Info().Msgf("finished click for every %ds", p.TimerSecond)

	msg += "the mouse controller loop was done!"
	log.Info().Msg(msg)

	return fmt.Errorf("%s", msg)
}
