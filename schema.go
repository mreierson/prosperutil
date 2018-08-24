package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05 -0700"
)

type ProsperTime struct {
	time.Time
}

func (pt *ProsperTime) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		pt.Time = time.Time{}
		return
	}
	pt.Time, err = time.Parse(TIME_LAYOUT, s)
	return
}

func (pt *ProsperTime) MarshalJSON() ([]byte, error) {
	if pt.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", pt.Time.Format(TIME_LAYOUT))), nil
}

var nilTime = (time.Time{}).UnixNano()

func (pt *ProsperTime) IsSet() bool {
	return pt.UnixNano() != nilTime
}

type AuthRequest struct {
	GrantType    string `url:"grant_type"`
	Username     string `url:"username"`
	Password     string `url:"password"`
	ClientId     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
}

type ProsperConn struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type AccountData struct {
	AvailableCash             float32       `json:"available_cash_balance"`
	PendingPrimary            float32       `json:"pending_investments_primary_market"`
	PendingSecondary          float32       `json:"pending_investments_secondary_market"`
	PendingQuickInvest        float32       `json:"pending_quick_invest_orders"`
	TotalPrincipalReceived    float32       `json:"total_principal_received_on_active_notes"`
	TotalAmountInvested       float32       `json:"total_amount_invested_on_active_notes"`
	TotalOutstandingPrincipal float32       `json:"outstanding_principal_on_active_notes"`
	TotalAccountValue         float32       `json:"total_account_value"`
	InvestedNotes             InvestedNotes `json:"invested_notes"`
}

type InvestedNotes struct {
	AA float32 `json:"AA"`
	A  float32 `json:"A"`
	B  float32 `json:"B"`
	C  float32 `json:"C"`
	D  float32 `json:"D"`
	E  float32 `json:"E"`
	HR float32 `json:"HR"`
	NA float32 `json:"NA"`
}

type NoteData struct {
	LoanNumber                           int     `json:"loan_number" db:"loan_number"`
	AmountBorrowed                       float32 `json:"amount_borrowed" db:"amount_borrowed"`
	BorrowerRate                         float32 `json:"borrower_rate" db:"borrower_rate"`
	ProsperRating                        string  `json:"prosper_rating" db:"prosper_rating"`
	Term                                 int     `json:"term" db:"term"`
	AgeInMonths                          int     `json:"age_in_months" db:"age_in_months"`
	OriginationDate                      string  `json:"origination_date" db:"origination_date"`
	DaysPastDue                          int     `json:"days_past_due" db:"days_past_due"`
	PrincipalBalanceProRataShare         float32 `json:"principal_balance_pro_rata_share" db:"principal_balance_pro_rata_share"`
	InterestPaidProRataShare             float32 `json:"interest_paid_pro_rata_share" db:"interest_paid_pro_rata_share"`
	ProsperFeesPaidProRataShare          float32 `json:"prosper_fees_paid_pro_rata_share" db:"prosper_fees_paid_pro_rata_share"`
	LateFeesPaidProRataShare             float32 `json:"late_fees_paid_pro_rata_share" db:"late_fees_paid_pro_rata_share"`
	DebtSaleProceedsReceivedProRataShare float32 `json:"debt_sale_proceeds_received_pro_rata_share" db:"debt_sale_proceeds_received_pro_rata_share"`
	NextPaymentDueAmountProRataShare     float32 `json:"next_payment_due_amount_pro_rata_share" db:"next_payment_due_amount_pro_rata_share"`
	NextPaymentDueDate                   string  `json:"next_payment_due_date" db:"next_payment_due_date"`
	LoanNoteId                           string  `json:"loan_note_id" db:"loan_note_id"`
	ListingNumber                        int     `json:"listing_number" db:"listing_number"`
	NoteOwnershipAmount                  float32 `json:"note_ownership_amount" db:"note_ownership_amount"`
	NoteSaleGrossAmountReceived          float32 `json:"note_sale_gross_amount_received" db:"note_sale_gross_amount_received"`
	NoteSaleFeesPaid                     float32 `json:"note_sale_fees_paid" db:"note_sale_fees_paid"`
	NoteStatus                           int     `json:"note_status" db:"note_status"`
	NoteStatusDesc                       string  `json:"note_status_description" db:"note_status_desc"`
	NoteDefaultReason                    int     `json:"note_default_reason" db:"note_default_reason"`
	NoteDefaultReasonDesc                string  `json:"note_default_reason_description" db:"note_default_reason_desc"`
	IsSold                               bool    `json:"is_sold" db:"is_sold"`
}

type NoteResult struct {
	Results     []NoteData `json:"result"`
	ResultCount int        `json:"result_count"`
	TotalCount  int        `json:"total_count"`
}

type ListingData struct {
	AmountDelinquent                          float32     `json:"amount_delinquent" db:"amount_delinquent"`
	AmountFunded                              float32     `json:"amount_funded" db:"amount_funded"`
	AmountParticipation                       float32     `json:"amount_participation" db:"amount_participation"`
	AmountRemaining                           float32     `json:"amount_remaining" db:"amount_remaining"`
	BankcardUtilization                       float32     `json:"bankcard_utilization" db:"bankcard_utilization"`
	BorrowerAPR                               float32     `json:"borrower_apr" db:"borrower_apr"`
	BorrowerCity                              string      `json:"borrower_city" db:"borrower_city"`
	BorrowerListingDesc                       string      `json:"borrower_listing_description" db:"borrower_listing_description"`
	BorrowerMetroArea                         string      `json:"borrower_metropolitan_area" db:"borrower_metropolitan_area"`
	BorrowerRate                              float32     `json:"borrower_rate" db:"borrower_rate"`
	BorrowerState                             string      `json:"borrower_state" db:"borrower_state"`
	CreditLinesLast7Years                     int         `json:"credit_lines_last7_years" db:"credit_lines_last7_years"`
	CreditPullDate                            ProsperTime `json:"credit_pull_date" db:"credit_pull_date"`
	CurrentCreditLines                        int         `json:"current_credit_lines" db:"current_credit_lines"`
	CurrentDelinquencies                      int         `json:"current_delinquencies" db:"current_delinquencies"`
	DelinquenciesLast7Years                   int         `json:"delinquencies_last7_years" db:"delinquencies_last7_years"`
	DelinquenciesOver30Days                   int         `json:"delinquencies_over30_days" db:"delinquencies_over30_days"`
	DelinquenciesOver60Days                   int         `json:"delinquencies_over60_days" db:"delinquencies_over60_days"`
	DelinquenciesOver90Days                   int         `json:"delinquencies_over90_days" db:"delinquencies_over90_days"`
	DtiWProsperLoan                           float32     `json:"dti_wprosper_loan" db:"dti_wprosper_loan"`
	EffectiveYield                            float32     `json:"effective_yield" db:"effective_yield"`
	EmploymentStatusDesc                      string      `json:"employment_status_description" db:"employment_status_description"`
	EstimatedLossRate                         float32     `json:"estimated_loss_rate" db:"estimated_loss_rate"`
	EstimatedReturn                           float32     `json:"estimated_return" db:"estimated_return"`
	FicoScore                                 string      `json:"fico_score" db:"fico_score"`
	FirstRecordedCreditLine                   ProsperTime `json:"first_recorded_credit_line" db:"first_recorded_credit_line"`
	FundingThreshold                          float32     `json:"funding_threshold" db:"funding_threshold"`
	IncomeRange                               int         `json:"income_range" db:"income_range"`
	IncomeRangeDesc                           string      `json:"income_range_description" db:"income_range_description"`
	IncomeVerifable                           bool        `json:"income_verifiable" db:"income_verifiable"`
	InquiriesLast6Months                      int         `json:"inquires_last6_months" db:"inquires_last6_months"`
	InstallmentBalance                        float32     `json:"installment_balance" db:"installment_balance"`
	InvestmentTypeDesc                        string      `json:"investment_type_description" db:"investment_type_description"`
	InvestmentTypeId                          int         `json:"investment_typeid" db:"investment_typeid"`
	IsHomeowner                               bool        `json:"is_homeowner" db:"is_homeowner"`
	LastUpdatedDate                           ProsperTime `json:"last_updated_date" db:"last_updated_date"`
	LenderIndicator                           int         `json:"lender_indicator" db:"lender_indicator"`
	LenderYield                               float32     `json:"lender_yield" db:"lender_yield"`
	ListingAmount                             float32     `json:"listing_amount" db:"listing_amount"`
	ListingCategoryId                         int         `json:"listing_category_id" db:"listing_category_id"`
	ListingCreationgDate                      ProsperTime `json:"listing_creation_date" db:"listing_creation_date"`
	ListingEndDate                            ProsperTime `json:"listing_end_date" db:"listing_end_date"`
	ListingMonthlyPayment                     float32     `json:"listing_monthly_payment" db:"listing_monthly_payment"`
	ListingNumber                             int         `json:"listing_number" db:"listing_number"`
	ListingStartDate                          ProsperTime `json:"listing_start_date" db:"listing_start_date"`
	ListingStatus                             int         `json:"listing_status" db:"listing_status"`
	ListingStatusReason                       string      `json:"listing_status_reason" db:"listing_status_reason"`
	ListingTerm                               int         `json:"listing_term" db:"listing_term"`
	ListingTitle                              string      `json:"listing_title" db:"listing_title"`
	LoanNumber                                int         `json:"loan_number" db:"loan_number"`
	LoanOriginationDate                       ProsperTime `json:"loan_origination_date" db:"loan_origination_date"`
	MaxPriorProsperLoan                       float32     `json:"max_prior_prosper_loan" db:"max_prior_prosper_loan"`
	MemberKey                                 string      `json:"member_key" db:"member_key"`
	MinPriorProsperLoan                       float32     `json:"min_prior_prosper_loan" db:"min_prior_prosper_loan"`
	MonthlyDebt                               float32     `json:"monthly_debt" db:"monthly_debt"`
	MonthsEmployed                            float32     `json:"months_employed" db:"months_employed"`
	NowDelinquentDerog                        int         `json:"now_delinquent_derog" db:"now_delinquent_derog"`
	Occupation                                string      `json:"occupation" db:"occupation"`
	OldestTradeOpenDate                       string      `json:"oldest_trade_open_date" db:"oldest_trade_open_date"`
	OpenCreditLines                           int         `json:"open_credit_lines" db:"open_credit_lines"`
	PartialFundingIndicator                   bool        `json:"partial_funding_indicator" db:"partial_funding_indicator"`
	PercentFunded                             float32     `json:"percent_funded" db:"percent_funded"`
	PriorProsperLoanEarliestPayOff            int         `json:"prior_prosper_loan_earliest_pay_off" db:"prior_prosper_loan_earliest_pay_off"`
	PriorProsperLoans31dpd                    int         `json:"prior_prosper_loans31dpd" db:"prior_prosper_loans31dpd"`
	PriorProsperLoans61dpd                    int         `json:"prior_prosper_loans61dpd" db:"prior_propser_loans61dpd"`
	PriorProsperLoans                         int         `json:"prior_prosper_loans" db:"prior_prosper_loans"`
	PriorProsperLoansActive                   int         `json:"prior_prosper_loans_active" db:"prior_prosper_loans_active"`
	PriorProsperLoansBalanceOutstanding       float32     `json:"prior_prosper_loans_balance_outstanding" db:"prior_prosper_loans_balance_outstanding"`
	PriorProsperLoansCyclesBilled             int         `json:"prior_prosper_loans_cycles_billed" db:"prior_prosper_loans_cycles_billed"`
	PriorProsperLoansLateCycles               int         `json:"prior_prosper_loans_late_cycles" db:"prior_prosper_loans_late_cycles"`
	PriorProsperLoansLatePaymentsOneMonthPlus int         `json:"prior_prosper_loans_late_payments_one_month_plus" db:"prior_prosper_loans_late_payments_one_month_plus"`
	PriorProsperLoansOntimePayments           int         `json:"prior_prosper_loans_ontime_payments" db:"prior_prosper_loans_ontime_payments"`
	PriorProsperLoansPrincipalBorrowed        float32     `json:"prior_prosper_loans_principal_borrowed" db:"prior_prosper_loans_principal_borrowed"`
	PriorProsperLoansPrincipalOutstanding     float32     `json:"prior_prosper_loans_principal_outstanding" db:"prior_prosper_loans_principal_outstanding"`
	ProsperRating                             string      `json:"prosper_rating" db:"prosper_rating"`
	ProsperScore                              int         `json:"prosper_score" db:"prosper_score"`
	PublicRecordsLast10Years                  int         `json:"public_records_last10_years" db:"public_records_last10_years"`
	PublicRecordsLast12Months                 int         `json:"public_records_last12_months" db:"public_records_last12_months"`
	RealEstateBalance                         float32     `json:"real_estate_balance" db:"real_estate_balance"`
	RealEstatePayment                         float32     `json:"real_estate_payment" db:"real_estate_payment"`
	RevolvingBalance                          float32     `json:"revolving_balance" db:"revolving_balance"`
	SatisfactoryAccounts                      int         `json:"satisfactory_accounts" db:"satisfactory_accounts"`
	ScoreX                                    string      `json:"scorex" db:"scorex"`
	ScoreXChange                              string      `json:"scorex_change" db:"scorex_change"`
	SipOfferId                                string      `json:"sip_offer_id" db:"sip_offer_id"`
	StatedMonthlyIncome                       float32     `json:"stated_monthly_income" db:"stated_monthly_income"`
	TotalInquiries                            int         `json:"total_inquiries" db:"total_inquiries"`
	TotalOpenRevolvingAccounts                int         `json:"total_open_revolving_accounts" db:"total_open_revolving_accounts"`
	TotalTradeItems                           int         `json:"total_trade_items" db:"total_trade_items"`
	VerificationStage                         int         `json:"verification_stage" db:"verification_stage"`
	WasDelinquentDerog                        int         `json:"was_delinquent_derog" db:"was_delinquent_derog"`
	WholeLoan                                 bool        `json:"whole_loan" db:"whole_loan"`
	WholeLoanEndDate                          ProsperTime `json:"whole_loan_end_date" db:"whole_loan_end_date"`
	WholeLoanStartDate                        ProsperTime `json:"whole_loan_start_date" db:"whole_loan_start_date"`
}

type ListingResult struct {
	Results     []ListingData `json:"result"`
	ResultCount int           `json:"result_count"`
	TotalCount  int           `json:"total_count"`
}

type BidRequest struct {
	ListingNumber int     `json:"listing_id"`
	BidAmount     float32 `json:"bid_amount"`
}

type BidRequests struct {
	Bids []BidRequest `json:"bid_requests"`
}

type OrderResponse struct {
	OrderId         string       `json:"order_id"`
	Bids            []BidRequest `json:"bid_requests"`
	EffectiveYield  float32      `json:"effective_yield"`
	EstimatedLoss   float32      `json:"estimated_loss"`
	EstimatedReturn float32      `json:"estimated_return"`
	OrderStatus     string       `json:"order_status"`
	OrderDate       string       `json:"order_date"`
}

type DailySnapshot struct {
	Timestamp                       time.Time `db:"timestamp"`
	AvailableCashBalance            float32   `db:"available_cash_balance"`
	PendingInvestmentsPrimaryMarket float32   `db:"pending_investments_primary_market"`
	TotalPrincipalReceived          float32   `db:"total_principal_received"`
	OutstandingPrincipal            float32   `db:"outstanding_principal"`
	TotalAccountValue               float32   `db:"total_account_value"`
	AllocationAATotal               float32   `db:"allocation_aa_total"`
	AllocationAAPct                 float32   `db:"allocation_aa_pct"`
	AllocationATotal                float32   `db:"allocation_a_total"`
	AllocationAPct                  float32   `db:"allocation_a_pct"`
	AllocationBTotal                float32   `db:"allocation_b_total"`
	AllocationBPct                  float32   `db:"allocation_b_pct"`
	AllocationCTotal                float32   `db:"allocation_c_total"`
	AllocationCPct                  float32   `db:"allocation_c_pct"`
	AllocationDTotal                float32   `db:"allocation_d_total"`
	AllocationDPct                  float32   `db:"allocation_d_pct"`
	AllocationETotal                float32   `db:"allocation_e_total"`
	AllocationEPct                  float32   `db:"allocation_e_pct"`
	AllocationHRTotal               float32   `db:"allocation_hr_total"`
	AllocationHRPct                 float32   `db:"allocation_hr_pct"`
	ActiveNotes                     int       `db:"active_notes"`
	CurrentNotes                    int       `db:"current_notes"`
	PastDueUnder30                  int       `db:"past_due_under30"`
	PastDueOver30                   int       `db:"past_due_over30"`
	ChargedOffNotes                 int       `db:"charged_off_notes"`
	PaidInFullNotes                 int       `db:"paid_in_full_notes"`
	SoldNotes                       int       `db:"sold_notes"`
	GainLossToDate                  float32   `db:"gain_loss_to_date"`
}
