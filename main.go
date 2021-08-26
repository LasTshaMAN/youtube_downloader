package main

import "youtube_downloader/dwnld"

const dir = "/home/lastshaman/youtube"

var urls = []string{
	"https://www.youtube.com/watch?v=e8nSEZLGoCA&ab_channel=StephanLivera",
	"https://www.youtube.com/watch?v=TCZm-J3nD9E&ab_channel=AltcoinDaily",
	"https://www.youtube.com/watch?v=VLZRriYx8qM&ab_channel=RAGIN%27AvC",
	"https://www.youtube.com/watch?v=dyDsQ-XCoNo&ab_channel=Bankless",
	"https://www.youtube.com/watch?v=UguvomJFqiE&ab_channel=PartTimeLarry",
	"https://www.youtube.com/watch?v=N4LlTiKYD98&list=PLmkdAgtxf3aiJo_1IkMVZr948id4W1p18&index=4&ab_channel=Bankless",
	"https://www.youtube.com/watch?v=vGjqqVN4qDs&ab_channel=Decred",
	"https://www.youtube.com/watch?v=Cz14dpWxgKk&ab_channel=Bankless",
	"https://www.youtube.com/watch?v=0ss-uSsPpRw&ab_channel=GrandAmphiTh%C3%A9atre",
	"https://www.youtube.com/watch?v=m-NGxJfS0mw&ab_channel=GrandAmphiTh%C3%A9atre",
	"https://www.youtube.com/watch?v=88GyLoZbDNw&ab_channel=BlackHat",
	"https://www.youtube.com/watch?v=T7_NTDYQq4k&ab_channel=JupiterBroadcasting",
	"https://www.youtube.com/watch?v=oCQou6xuXbk&ab_channel=CoinBureau",
	"https://www.youtube.com/watch?v=7NdjivxrDoc&ab_channel=CoinBureau",
	"https://www.youtube.com/watch?v=lANUSdHg2oc&ab_channel=CoinBureau",
	"https://www.youtube.com/watch?v=7tFL0NFBwUc&ab_channel=linux.conf.au",
	"https://www.youtube.com/watch?v=A3G-3hp88mo&ab_channel=NetworkChuck",
	"https://www.youtube.com/watch?v=l7Urkkr1y0s&ab_channel=DJWare",
	"https://www.youtube.com/watch?v=X65hV4mkulM&ab_channel=CyberInitiative",
	"https://www.youtube.com/watch?v=BO8ZSAfh0Jo&ab_channel=Web3Foundation",
	"https://www.youtube.com/watch?v=xqopwqXyURw&ab_channel=MITBitcoinClub",
	"https://www.youtube.com/watch?v=9nXop2lLDa4&ab_channel=TalksatGoogle",
	"https://www.youtube.com/watch?v=a51OpyZYiA8&ab_channel=TheDefiant",
	"https://www.youtube.com/watch?v=KZmLmt9f_uk&ab_channel=EthereumFoundation",
}

// The library used for downloading is https://github.com/kkdai/youtube
// For an example on how to download a playlist see - https://github.com/kkdai/youtube/blob/master/example_test.go.

func main() {
	dwnld.DownloadBatch(dir, urls)
}
