package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/rain/torrent"
	tb "gopkg.in/telebot.v3"
)

var (
	client  *torrent.Session
	TOKEN   = os.Getenv("TOKEN")
	bot     *tb.Bot
	workDir string
)

func SetupRainTorrent() *torrent.Session {
	SetupWorkDir()
	config := torrent.DefaultConfig
	config.DataDir = workDir + "/torrents"
	config.Database = workDir + "torrents.db"
	client, err := torrent.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func SetupBot() *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token: TOKEN,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func SetupWorkDir() {
	if _, err := os.Stat(workDir + "/torrents"); os.IsNotExist(err) {
		os.MkdirAll(workDir+"/torrents", 0755)
	}
}

func init() {
	wd, _ := os.Getwd()
	workDir = wd
	log.Println("Work dir: " + workDir)
	client = SetupRainTorrent()
	bot = SetupBot()
}

func GetDownloadPercentage(torr *torrent.Torrent) string {
	if torr != nil {
		if torr.Stats().Pieces.Total != 0 {
			return fmt.Sprintf("%.2f", float64(torr.Stats().Pieces.Have)/float64(torr.Stats().Pieces.Total)*100) + "%"
		}
	}
	return "0%"
}

func GetTorrentSize(torr *torrent.Torrent) int64 {
	if torr != nil {
		if torr.Stats().Bytes.Total != 0 {
			return torr.Stats().Bytes.Total
		}
	}
	return 0
}

func GetDownloadSpeed(t *torrent.Torrent) string {
	if t.Stats().Speed.Download != 0 {
		return ByteCountSI(int64(t.Stats().Speed.Download)) + "/s"
	} else {
		return "-/-"
	}
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func GenProgressMessage(t *torrent.Torrent) string {
	PROGRESS_MSG := ""
	PROGRESS_MSG += "Name: <code>" + t.Name() + "</code>\n"
	PROGRESS_MSG += "Size: <code>" + ByteCountSI(GetTorrentSize(t)) + "</code>\n"
	PROGRESS_MSG += "Status: <code>" + t.Stats().Status.String() + "</code>\n"
	PROGRESS_MSG += "Download speed: <code>" + GetDownloadSpeed(t) + "</code>\n"
	PROGRESS_MSG += "ETA: <code>" + fmt.Sprint(t.Stats().ETA) + "</code>\n"
	PROGRESS_MSG += "Progress: <code>" + GetDownloadPercentage(t) + "</code>\n"
	return PROGRESS_MSG
}

func AddEventHandlers() {
	bot.Handle("/start", Start)
	bot.Handle("/help", Help)
	bot.Handle("/mirror", AddMirror)
	bot.Handle(&HelpMain, HelpMainCb)
	bot.Handle(&HelpTorrent, HelpTorrentCb)
}
