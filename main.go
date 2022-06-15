package main

import (
	"log"
	"time"

	"github.com/cenkalti/rain/torrent"
	tb "gopkg.in/telebot.v3"
)

var btn = &tb.ReplyMarkup{}
var (
	HelpMain    = btn.Data("help", "help")
	HelpTorrent = btn.Data("help_torrent", "help_torrent")
)

const (
	help_main_msg = `
	Help for core commands:
	/start - Start the bot
	/help - Show this message
	`
	help_torrent_msg = `
	Help for torrent commands:
	/mirror - mirror a torrent from magnet link
    /status - show the status of the torrent
    /cancel - cancel the torrent
    /list - list all the torrents
    /search - search for a torrent
    /kill - kill all downloads
    /pause - pause the torrent
    /resume - resume the torrent
    /remove - remove the torrent
    /info - show the info of the torrent`
)

func Start(c tb.Context) error {
	const start_msg = `
	I am a bot that can help you download torrents.
	You can use me to download torrents from torrent sites.`
	btn.Inline(btn.Row(
		btn.URL("Support", "t.me/rosekcd"),
	))
	return c.Reply(start_msg, &tb.SendOptions{
		ReplyMarkup: btn,
	})
}

func Help(c tb.Context) error {
	if c.Message().Payload == "torrent" {
		return c.Reply(help_torrent_msg)
	}
	const help_msg = `Here are some commands you can use:`
	HelpMain.Text = "Core commands"
	HelpTorrent.Text = "Torrent commands"
	btn.Inline(btn.Row(
		HelpMain,
		HelpTorrent,
	))
	return c.Reply(help_msg, &tb.SendOptions{
		ReplyMarkup: btn,
	})
}

func HelpTorrentCb(c tb.Context) error {
	return c.Edit(help_torrent_msg)
}

func HelpMainCb(c tb.Context) error {
	return c.Edit(help_main_msg)
}

func AddMirror(c tb.Context) error {
	magnet := c.Message().Payload
	if magnet == "" {
		return c.Reply("Please send a magnet link")
	}
	torr, err := client.AddURI(magnet, &torrent.AddTorrentOptions{
		StopAfterDownload: false,
	})
	if err != nil {
		return c.Reply(err.Error())
	}
	msg, _ := c.Bot().Send(c.Chat(), "Torrent added!", &tb.SendOptions{
		ReplyTo: c.Message(),
	})
	for range time.Tick(time.Second * 6) {
		t := client.GetTorrent(torr.ID())
		progress := GenProgressMessage(t)
		if t == nil {
			break
		}
		if t.Stats().Status == torrent.Downloading || t.Stats().Status == torrent.DownloadingMetadata {
			if progress != msg.Text {
				msg, _ = c.Bot().Edit(msg, progress)
			}
		} else if t.Stats().Status == torrent.Stopped && GetDownloadPercentage(t) == "100%" {
			c.Bot().Edit(msg, "Torrent finished!")
			break
		} else if t.Stats().Status == torrent.Stopped {
			c.Bot().Edit(msg, "Torrent stopped!")
			break
		}
	}
	return nil
}

func main() {
	AddEventHandlers()
	log.Println("Starting...")
	bot.Start()
}
