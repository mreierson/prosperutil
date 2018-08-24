package main

import (
	_ "bytes"
	_ "encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	_ "gopkg.in/resty.v0"
	"log"
	_ "net/http"
	_ "net/url"
)

const (
	CLIENT_ID     = "CLIENT_ID"
	CLIENT_SECRET = "CLIENT_SECRET"

	BASE_URL           = "https://api.prosper.com/v1/"
	BASE_URL_NEW       = "https://api.prosper.com/"
	OAUTH_URL          = "security/oauth/token"
	ACCOUNT_URL        = "accounts/prosper"
	ORDERS_URL         = "orders/"
	NOTE_URL           = "notes/"
	MY_LISTING_URL     = "listingsvc/v2/listings/"
	ACTIVE_LISTING_URL = "listingsvc/v2/listings/?biddable=true&invested=false&limit=500"
)

func OpenAPI() (conn ProsperConn, err error) {

	form := AuthRequest{
		GrantType:    "password",
		Username:     "USER_NAME",
		Password:     "PASSWORD",
		ClientId:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
	}

	_, err = sling.New().Base(BASE_URL).Post(OAUTH_URL).BodyForm(form).ReceiveSuccess(&conn)
	return conn, err
}

func (conn ProsperConn) GetAccountData() (acctData AccountData, err error) {

	_, err = conn.buildBaseRequest().Get(ACCOUNT_URL).ReceiveSuccess(&acctData)
	return acctData, err
}

func (conn ProsperConn) GetNotes() (notes []NoteData, err error) {

	var noteResult NoteResult

	_, err = conn.buildBaseRequest().Get(NOTE_URL).ReceiveSuccess(&noteResult)
	if err != nil {
		return notes, err
	}

	notes = append(notes, noteResult.Results...)

	currentOffset := noteResult.ResultCount
	countRemaining := noteResult.TotalCount - currentOffset

	for countRemaining > 0 {

		currentUrl := fmt.Sprintf("%s?offset=%d", NOTE_URL, currentOffset)

		_, err := conn.buildBaseRequest().Get(currentUrl).ReceiveSuccess(&noteResult)
		if err != nil {
			return notes, err
		}

		notes = append(notes, noteResult.Results...)
		currentOffset += noteResult.ResultCount
		countRemaining -= noteResult.ResultCount
	}

	return notes, err
}

func (conn ProsperConn) GetListing(listingNumber int) (listing ListingData, err error) {

	var listingResult ListingResult

	url := fmt.Sprintf("%s?listing_number=%d", MY_LISTING_URL, listingNumber)

	_, err = conn.buildBaseRequestNew().Get(url).ReceiveSuccess(&listingResult)
	if err != nil {
		return listing, err
	}

	if listingResult.ResultCount > 0 {

		listing = listingResult.Results[0]

	} else {
		err = fmt.Errorf("Listing NOT found : listing_number = %d", listingNumber)
	}

	return listing, err
}

func (conn ProsperConn) GetActiveListings() (listings []ListingData, err error) {

	var listingResult ListingResult

	_, err = conn.buildBaseRequestNew().Get(ACTIVE_LISTING_URL).ReceiveSuccess(&listingResult)
	if err != nil {
		return listings, err
	}

	log.Printf("total count %d\n", listingResult.TotalCount)

	listings = append(listings, listingResult.Results...)

	currentOffset := listingResult.ResultCount
	countRemaining := listingResult.TotalCount - currentOffset

	for countRemaining > 0 {

		currentUrl := fmt.Sprintf("%s&offset=%d", ACTIVE_LISTING_URL, currentOffset)

		_, err = conn.buildBaseRequest().Get(currentUrl).ReceiveSuccess(&listingResult)
		if err != nil {
			return listings, err
		}

		listings = append(listings, listingResult.Results...)
		currentOffset += listingResult.ResultCount
		countRemaining -= listingResult.ResultCount
	}

	return listings, err
}

func (conn ProsperConn) Invest(listingNumber int, amount float32) (orderResponse OrderResponse, err error) {

	bidRequests := &BidRequests{}
	bidRequests.Bids = append(bidRequests.Bids, BidRequest{
		ListingNumber: listingNumber,
		BidAmount:     amount,
	})

	resp, err := conn.buildBaseRequest().Post(ORDERS_URL).Set("Content-Type", "application/json").BodyJSON(&bidRequests).ReceiveSuccess(&orderResponse)
	if resp.StatusCode >= 300 {
		err = fmt.Errorf("StatusCode = %d : %s", resp.StatusCode, resp.Status)
	}
	if err != nil {
		fmt.Println(err)
	}

	return orderResponse, err
}

func (conn ProsperConn) buildBaseRequest() *sling.Sling {
	return sling.New().Base(BASE_URL).Set("Authorization", "bearer "+conn.AccessToken).Set("Accept", "application/json").Set("timezone", "America/Chicago")
}

func (conn ProsperConn) buildBaseRequestNew() *sling.Sling {
	return sling.New().Base(BASE_URL_NEW).Set("Authorization", "bearer "+conn.AccessToken).Set("Accept", "application/json").Set("timezone", "America/Chicago")
}

func (l ListingData) String() string {

	return fmt.Sprintf("Listing { Id=%d, %s, %s, %s, %.2f, %s, %.2f, %s, %s, %s, %.2f, %.2f }",
		l.ListingNumber,
		l.ListingTitle,
		l.ProsperRating,
		l.FicoScore,
		l.RevolvingBalance,
		l.IncomeRangeDesc,
		l.EffectiveYield,
		l.EmploymentStatusDesc,
		l.BorrowerCity,
		l.BorrowerState,
		l.PercentFunded,
		l.MonthsEmployed)
}
