package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"time"
)

var (
	CACHE = cache.New(5*time.Minute, 0*time.Second)
)

type WebSnapshot struct {
	Timestamp                       []time.Time `json:"timestamp"`
	AvailableCashBalance            []float32   `json:"available_cash_balance"`
	PendingInvestmentsPrimaryMarket []float32   `json:"pending_investments_primary_market"`
	TotalPrincipalReceived          []float32   `json:"total_principal_received"`
	OutstandingPrincipal            []float32   `json:"outstanding_principal"`
	TotalAccountValue               []float32   `json:"total_account_value"`
	AllocationAATotal               []float32   `json:"allocation_AA_total"`
	AllocationAAPct                 []float32   `json:"allocation_AA_pct"`
	AllocationATotal                []float32   `json:"allocation_A_total"`
	AllocationAPct                  []float32   `json:"allocation_A_pct"`
	AllocationBTotal                []float32   `json:"allocation_B_total"`
	AllocationBPct                  []float32   `json:"allocation_B_pct"`
	AllocationCTotal                []float32   `json:"allocation_C_total"`
	AllocationCPct                  []float32   `json:"allocation_C_pct"`
	AllocationDTotal                []float32   `json:"allocation_D_total"`
	AllocationDPct                  []float32   `json:"allocation_D_pct"`
	AllocationETotal                []float32   `json:"allocation_E_total"`
	AllocationEPct                  []float32   `json:"allocation_E_pct"`
	AllocationHRTotal               []float32   `json:"allocation_HR_total"`
	AllocationHRPct                 []float32   `json:"allocation_HR_pct"`
	ActiveNotes                     []int       `json:"active_notes"`
	CurrentNotes                    []int       `json:"current_notes"`
	PastDueUnder30                  []int       `json:"past_due_under30"`
	PastDueOver30                   []int       `json:"past_due_over30"`
	ChargedOffNotes                 []int       `json:"charged_off_notes"`
	PaidInFullNotes                 []int       `json:"paid_in_full_notes"`
	SoldNotes                       []int       `json:"sold_notes"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func DailySnapshotHandler(w http.ResponseWriter, r *http.Request) {

	ws, found := CACHE.Get("DailySnapshot")
	if !found {

		db, err := OpenDatabase()
		if err != nil {
			log.Fatal(err)
		}

		snapshots, err := db.LoadDailySnapshots()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Loaded %d snapshots", len(snapshots))

		ws := WebSnapshot{}
		for _, snapshot := range snapshots {
			ws.Timestamp = append(ws.Timestamp, snapshot.Timestamp)
			ws.AvailableCashBalance = append(ws.AvailableCashBalance, snapshot.AvailableCashBalance)
			ws.PendingInvestmentsPrimaryMarket = append(ws.PendingInvestmentsPrimaryMarket, snapshot.PendingInvestmentsPrimaryMarket)
			ws.TotalPrincipalReceived = append(ws.TotalPrincipalReceived, snapshot.TotalPrincipalReceived)
			ws.OutstandingPrincipal = append(ws.OutstandingPrincipal, snapshot.OutstandingPrincipal)
			ws.TotalAccountValue = append(ws.TotalAccountValue, snapshot.TotalAccountValue)
			ws.AllocationAATotal = append(ws.AllocationAATotal, snapshot.AllocationAATotal)
			ws.AllocationAAPct = append(ws.AllocationAAPct, snapshot.AllocationAAPct)
			ws.AllocationATotal = append(ws.AllocationATotal, snapshot.AllocationATotal)
			ws.AllocationAPct = append(ws.AllocationAPct, snapshot.AllocationAPct)
			ws.AllocationBTotal = append(ws.AllocationBTotal, snapshot.AllocationBTotal)
			ws.AllocationBPct = append(ws.AllocationBPct, snapshot.AllocationBPct)
			ws.AllocationCTotal = append(ws.AllocationCTotal, snapshot.AllocationCTotal)
			ws.AllocationCPct = append(ws.AllocationCPct, snapshot.AllocationCPct)
			ws.AllocationDTotal = append(ws.AllocationDTotal, snapshot.AllocationDTotal)
			ws.AllocationDPct = append(ws.AllocationDPct, snapshot.AllocationDPct)
			ws.AllocationETotal = append(ws.AllocationETotal, snapshot.AllocationETotal)
			ws.AllocationEPct = append(ws.AllocationEPct, snapshot.AllocationEPct)
			ws.AllocationHRTotal = append(ws.AllocationHRTotal, snapshot.AllocationHRTotal)
			ws.AllocationHRPct = append(ws.AllocationHRPct, snapshot.AllocationHRPct)
			ws.ActiveNotes = append(ws.ActiveNotes, snapshot.ActiveNotes)
			ws.CurrentNotes = append(ws.CurrentNotes, snapshot.CurrentNotes)
			ws.PastDueUnder30 = append(ws.PastDueUnder30, snapshot.PastDueUnder30)
			ws.PastDueOver30 = append(ws.PastDueOver30, snapshot.PastDueOver30)
			ws.ChargedOffNotes = append(ws.ChargedOffNotes, snapshot.ChargedOffNotes)
			ws.PaidInFullNotes = append(ws.PaidInFullNotes, snapshot.PaidInFullNotes)
			ws.SoldNotes = append(ws.SoldNotes, snapshot.SoldNotes)
		}

		CACHE.Set("DailySnapshot", ws, cache.DefaultExpiration)
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err := enc.Encode(&ws)
	if err != nil {
		log.Fatal(err)
	}
}

func RunWebServer() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/prosper", DailySnapshotHandler)

	n := negroni.New()
	n.UseHandler(r)
	n.Run(":9090")
}
