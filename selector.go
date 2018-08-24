package main

import (
	"log"
	"sort"
)

type ByEffectiveYield []ListingData

func (l ByEffectiveYield) Len() int {
	return len(l)
}

func (l ByEffectiveYield) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ByEffectiveYield) Less(i, j int) bool {
	return l[j].EffectiveYield < l[i].EffectiveYield
}

func FilterListings(listings []ListingData) (filtered []ListingData) {

	for _, listing := range listings {

		if "AA" == listing.ProsperRating || "A" == listing.ProsperRating {

			log.Println("Filtering loan (Rating) : ", listing)

		} else if listing.PublicRecordsLast10Years > 0 || listing.PublicRecordsLast12Months > 0 {

			log.Println("Filtering loan (PublicRecord) : ", listing)

		} else if listing.DelinquenciesLast7Years > 0 {

			log.Println("Filtering loan (Delinquencies) : ", listing)

		} else if listing.PercentFunded <= 0.7 {

			log.Println("Filtering loan (% Funded) : ", listing)

		} else if listing.MonthsEmployed < 40 {

			log.Println("Filtering loan (EmploymentMonths) : ", listing)

		} else if listing.WasDelinquentDerog > 0 || listing.NowDelinquentDerog > 0 {

			log.Println("Filtering loan (DelinquentDerog) : ", listing)

		} else if listing.VerificationStage < 2 {

			log.Println("Filtering loan (VerificationStage) : ", listing)

		} else {

			filtered = append(filtered, listing)

		}

	}

	sort.Sort(ByEffectiveYield(filtered))

	return filtered
}
