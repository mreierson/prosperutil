package main

import (
	_ "fmt"
	"log"
	"math"
	_ "net/http"
	"time"
)

const (
	NOTE_SHARE  = 35
	MIN_BALANCE = 0
)

func AutoInvest() {

	conn, err := OpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	acct, _ := conn.GetAccountData()

	log.Printf("Note share   : %.2f\n", float32(NOTE_SHARE))
	log.Printf("Min. balance : %.2f\n", float32(MIN_BALANCE))
	log.Printf("Cash balance : %.2f\n", acct.AvailableCash)
	log.Printf("Pending      : %.2f\n", acct.PendingPrimary+acct.PendingQuickInvest)

	availableCash := math.Max(float64(acct.AvailableCash-MIN_BALANCE), 0)
	if availableCash < NOTE_SHARE {
		log.Println("Not enough available cash")
	} else {

		listings, err := conn.GetActiveListings()
		if err != nil {
			log.Println(err)
		}
		log.Println("Listings :", len(listings))

		filtered := FilterListings(listings)

		for _, listing := range filtered {

			if availableCash < NOTE_SHARE {
				break
			}

			log.Println("Purchasing ", listing)

			resp, err := conn.Invest(listing.ListingNumber, NOTE_SHARE)
			if err != nil {
				log.Fatal(err)
			} else {
				availableCash -= float64(resp.Bids[0].BidAmount)
			}
		}
	}
}

func TakeDailySnapshot() {

	conn, err := OpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	acct, err := conn.GetAccountData()
	if err != nil {
		log.Fatal(err)
	}

	notes, err := conn.GetNotes()
	if err != nil {
		log.Fatal(err)
	}

	var (
		activeNotes             = 0
		currentNotes            = 0
		pastDueUnder30          = 0
		pastDueOver30           = 0
		chargedOffNotes         = 0
		paidInFullNotes         = 0
		soldNotes               = 0
		totalAlloc      float32 = 0
	)

	for _, note := range notes {

		if note.NoteStatus == 1 {

			activeNotes += 1

			if note.DaysPastDue == 0 {
				currentNotes += 1
			} else if note.DaysPastDue <= 30 {
				pastDueUnder30 += 1
			} else {
				pastDueOver30 += 1
			}

		} else if note.NoteStatus == 2 {

			chargedOffNotes += 1

		} else if note.NoteStatus == 3 {

			if note.IsSold {
				soldNotes += 1
			} else {
				chargedOffNotes += 1
			}

		} else if (note.NoteStatus == 4) || (note.NoteStatus == 6) {

			paidInFullNotes += 1

		}

	}

	totalAlloc += acct.InvestedNotes.AA
	totalAlloc += acct.InvestedNotes.A
	totalAlloc += acct.InvestedNotes.B
	totalAlloc += acct.InvestedNotes.C
	totalAlloc += acct.InvestedNotes.D
	totalAlloc += acct.InvestedNotes.E
	totalAlloc += acct.InvestedNotes.HR

	db, err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = db.StoreDailySnapshot(DailySnapshot{
		Timestamp:                       time.Now(),
		AvailableCashBalance:            acct.AvailableCash,
		PendingInvestmentsPrimaryMarket: acct.PendingPrimary,
		TotalPrincipalReceived:          acct.TotalPrincipalReceived,
		OutstandingPrincipal:            acct.TotalOutstandingPrincipal,
		TotalAccountValue:               acct.TotalAccountValue,
		AllocationAATotal:               acct.InvestedNotes.AA,
		AllocationAAPct:                 acct.InvestedNotes.AA / totalAlloc,
		AllocationATotal:                acct.InvestedNotes.A,
		AllocationAPct:                  acct.InvestedNotes.A / totalAlloc,
		AllocationBTotal:                acct.InvestedNotes.B,
		AllocationBPct:                  acct.InvestedNotes.B / totalAlloc,
		AllocationCTotal:                acct.InvestedNotes.C,
		AllocationCPct:                  acct.InvestedNotes.C / totalAlloc,
		AllocationDTotal:                acct.InvestedNotes.D,
		AllocationDPct:                  acct.InvestedNotes.D / totalAlloc,
		AllocationETotal:                acct.InvestedNotes.E,
		AllocationEPct:                  acct.InvestedNotes.E / totalAlloc,
		AllocationHRTotal:               acct.InvestedNotes.HR,
		AllocationHRPct:                 acct.InvestedNotes.HR / totalAlloc,
		ActiveNotes:                     activeNotes,
		CurrentNotes:                    currentNotes,
		PastDueUnder30:                  pastDueUnder30,
		PastDueOver30:                   pastDueOver30,
		ChargedOffNotes:                 chargedOffNotes,
		PaidInFullNotes:                 paidInFullNotes,
		SoldNotes:                       soldNotes,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func SyncNotes() {

	conn, err := OpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	notes, err := conn.GetNotes()
	if err != nil {
		log.Fatal(err)
	}

	db, err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = db.StoreNotes(notes)
	if err != nil {
		log.Fatal(err)
	}
}

func SyncListings() {

	conn, err := OpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	db, err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	listingNumbers, err := db.QueryMissingListings()
	if err != nil {
		log.Fatal(err)
	}

	var listings []ListingData
	var missing []int

	for _, listingNumber := range listingNumbers {
		listing, err := conn.GetListing(listingNumber)
		if err != nil {
			log.Println(err)
			missing = append(missing, listingNumber)
		}
		listings = append(listings, listing)
	}

	db.SetMissingFlags(missing)
	db.StoreListings(listings)
}

func ListingDetail(listingNumber int) {

	conn, err := OpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	listing, err := conn.GetListing(listingNumber)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(listing)
	}
}

func ListingAll() {

	conn, err := OpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	listings, err := conn.GetActiveListings()
	if err != nil {
		log.Fatal(err)
	} else {
		for _, listing := range listings {
			log.Println(listing)
		}
		log.Printf("%d listing(s)", len(listings))
	}
}

func DumpSnapshots() {

	db, err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	snapshots, err := db.LoadDailySnapshots()
	if err != nil {
		log.Fatal(err)
	} else {
		for _, snapshot := range snapshots {
			log.Println(snapshot)
		}
	}
}
