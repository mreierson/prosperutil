package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var (
	autoInvest    = kingpin.Flag("invest", "Auto-Invest").Short('i').Bool()
	dailySnapshot = kingpin.Flag("daily-snapshot", "Daily Snapshot").Short('d').Bool()
	syncNotes     = kingpin.Flag("note-sync", "Sync. notes w/ database").Short('n').Bool()
	syncListings  = kingpin.Flag("listing-sync", "Sync. listings w/ database").Short('l').Bool()
	webServer     = kingpin.Flag("server", "Run webapp").Short('s').Bool()
	listingDetail = kingpin.Flag("listing-detail", "Show listing detail by id").Int()
	listingAll    = kingpin.Flag("listing-all", "Show all available listings").Bool()
	dumpSnapshot  = kingpin.Flag("dump-snapshots", "Print all snapshots in database").Bool()
)

func main() {

	kingpin.Parse()

	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "/opt/pinvest/logs/pinvest.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     10,
	}))

	if *autoInvest {
		log.Println("Starting AutoInvest")
		AutoInvest()
		log.Println("Finished AutoInvest")
	}

	if *dailySnapshot {
		log.Println("Starting DailySnapshot")
		TakeDailySnapshot()
		log.Println("Finished DailySnapshot")
	}

	if *syncNotes {
		log.Println("Starting SyncNotes")
		SyncNotes()
		log.Println("Finished SyncNotes")
	}

	if *syncListings {
		log.Println("Starting SyncListings")
		SyncListings()
		log.Println("Finished SyncListings")
	}

	if *webServer {
		log.Println("Starting WebServer")
		RunWebServer()
		log.Println("Finished WebServer")
	}

	if *listingDetail > 0 {
		log.Println("Fetching listing detail for", *listingDetail)
		ListingDetail(*listingDetail)
		log.Println("Finished fetching listing detail")
	}

	if *listingAll {
		log.Println("Fetching all available listings")
		ListingAll()
		log.Println("Finished fetching all available listings")
	}

	if *dumpSnapshot {
		log.Println("Dumping database snapshots")
		DumpSnapshots()
	}
}
