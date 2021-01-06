package chatex

import "fmt"

func OrderLink(id uint64) string {
	return fmt.Sprintf("https://t.me/Chatex_bot?start=ad_%v", id)
}

func OrderLinkMd(id uint64) string {
	return fmt.Sprintf("[%v](%v)", id, OrderLink(id))
}