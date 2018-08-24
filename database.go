package main

import (
	_ "database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	DATABASE_TIME_LAYOUT = "2006-01 02 15:04:05"
)

type DatabaseConn struct {
	Db *sqlx.DB
}

func (pt ProsperTime) Value() (driver.Value, error) {

	if pt.Time.UnixNano() == nilTime {
		return nil, nil
	}
	return driver.Value(fmt.Sprintf("%s", pt.Time.Format(TIME_LAYOUT))), nil
}

func (pt ProsperTime) Scan(value interface{}) (err error) {
	if value == nil {
		pt.Time = time.Time{}
		return nil
	}
	pt.Time, err = time.Parse(TIME_LAYOUT, value.(string))
	return err
}

func OpenDatabase() (conn DatabaseConn, err error) {

	conn.Db, err = sqlx.Connect("postgres", "CONNECTION_STRING")
	if err != nil {
		return conn, err
	}

	return conn, err
}

func (conn *DatabaseConn) LoadDailySnapshots() (snapshots []DailySnapshot, err error) {

	var QUERY = `SELECT timestamp, available_cash_balance, pending_investments_primary_market,
       total_principal_received, outstanding_principal, total_account_value,
       allocation_aa_total, allocation_aa_pct, allocation_a_total, allocation_a_pct,
       allocation_b_total, allocation_b_pct, allocation_c_total, allocation_c_pct,
       allocation_d_total, allocation_d_pct, allocation_e_total, allocation_e_pct,
       allocation_hr_total, allocation_hr_pct, active_notes, current_notes,
       past_due_under30, past_due_over30, charged_off_notes, paid_in_full_notes,
       sold_notes
       FROM dailysnapshot ORDER BY timestamp ASC`

	rows, err := conn.Db.Queryx(QUERY)
	if err != nil {
		return snapshots, err
	}

	snapshot := DailySnapshot{}
	for rows.Next() {
		err = rows.StructScan(&snapshot)
		if err != nil {
			return snapshots, err
		}
		snapshots = append(snapshots, snapshot)
	}
	rows.Close()

	return snapshots, err
}

func (conn *DatabaseConn) StoreDailySnapshot(snap DailySnapshot) (err error) {

	var INSERT = `INSERT INTO dailysnapshot(timestamp, available_cash_balance,
    pending_investments_primary_market, total_principal_received, outstanding_principal,
    total_account_value, allocation_aa_total, allocation_aa_pct, allocation_a_total, allocation_a_pct,
		allocation_b_total, allocation_b_pct, allocation_c_total, allocation_c_pct,
		allocation_d_total, allocation_d_pct, allocation_e_total, allocation_e_pct,
		allocation_hr_total, allocation_hr_pct, active_notes, current_notes,
		past_due_under30, past_due_over30, charged_off_notes, paid_in_full_notes,
		sold_notes, gain_loss_to_date)
    VALUES(:timestamp, :available_cash_balance, :pending_investments_primary_market,
    :total_principal_received, :outstanding_principal, :total_account_value, :allocation_aa_total,
    :allocation_aa_pct, :allocation_a_total, :allocation_a_pct, :allocation_b_total,
    :allocation_b_pct, :allocation_c_total, :allocation_c_pct, :allocation_d_total,
    :allocation_d_pct, :allocation_e_total, :allocation_e_pct, :allocation_hr_total,
    :allocation_hr_pct, :active_notes, :current_notes, :past_due_under30, :past_due_over30,
    :charged_off_notes, :paid_in_full_notes, :sold_notes, :gain_loss_to_date)`

	_, err = conn.Db.NamedExec(INSERT, snap)

	return err
}

func (conn *DatabaseConn) StoreNotes(notes []NoteData) (err error) {

	var QUERY = `SELECT listing_number, loan_number, loan_note_id, amount_borrowed,
	  borrower_rate, prosper_rating, term, age_in_months, to_char(origination_date, 'YYYY-MM-DD') origination_date,
	  days_past_due, principal_balance_pro_rata_share, interest_paid_pro_rata_share,
	  prosper_fees_paid_pro_rata_share, late_fees_paid_pro_rata_share,
	  debt_sale_proceeds_received_pro_rata_share, next_payment_due_amount_pro_rata_share,
	  to_char(next_payment_due_date, 'YYYY-MM-DD') next_payment_due_date, note_ownership_amount, note_sale_gross_amount_received,
	  note_sale_fees_paid, note_status, note_status_desc, note_default_reason,
	  note_default_reason_desc, is_sold
	  FROM notes WHERE active='Y' AND listing_number=:listing_number`

	var INSERT = `INSERT INTO notes(
    listing_number, loan_number, loan_note_id, amount_borrowed,
    borrower_rate, prosper_rating, term, age_in_months, origination_date,
    days_past_due, principal_balance_pro_rata_share, interest_paid_pro_rata_share,
    prosper_fees_paid_pro_rata_share, late_fees_paid_pro_rata_share,
    debt_sale_proceeds_received_pro_rata_share, next_payment_due_amount_pro_rata_share,
    next_payment_due_date, note_ownership_amount, note_sale_gross_amount_received,
    note_sale_fees_paid, note_status, note_status_desc, note_default_reason,
    note_default_reason_desc, is_sold, active)
    VALUES(:listing_number, :loan_number, :loan_note_id, :amount_borrowed,
    :borrower_rate, :prosper_rating, :term, :age_in_months, :origination_date,
    :days_past_due, :principal_balance_pro_rata_share, :interest_paid_pro_rata_share,
    :prosper_fees_paid_pro_rata_share, :late_fees_paid_pro_rata_share,
    :debt_sale_proceeds_received_pro_rata_share, :next_payment_due_amount_pro_rata_share,
    :next_payment_due_date, :note_ownership_amount, :note_sale_gross_amount_received,
    :note_sale_fees_paid, :note_status, :note_status_desc, :note_default_reason,
    :note_default_reason_desc, :is_sold, 'Y')`

	var UPDATE = `UPDATE notes SET active='N' WHERE active='Y' AND listing_number=:listing_number`

	tx := conn.Db.MustBegin()
	for _, note := range notes {

		rows, err := tx.NamedQuery(QUERY, note)
		if err != nil {
			return err
		}

		dbNote := NoteData{}
		for rows.Next() {
			err := rows.StructScan(&dbNote)
			if err != nil {
				return err
			}

			rows.Close()
			break
		}

		if dbNote != note {
			_, err = tx.NamedExec(UPDATE, note)
			if err != nil {
				return err
			}
			_, err = tx.NamedExec(INSERT, note)
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()

	return err
}

func (conn *DatabaseConn) QueryMissingListings() (listingNumbers []int, err error) {

	var QUERY_MISSING = `SELECT n.listing_number FROM notes n LEFT OUTER JOIN listings l ON n.listing_number = l.listing_number WHERE n.active='Y' AND n.missing <> 'Y' AND l.listing_number IS NULL`

	rows, err := conn.Db.Query(QUERY_MISSING)
	if err != nil {
		return listingNumbers, err
	}

	var listingNumber int
	for rows.Next() {
		rows.Scan(&listingNumber)
		listingNumbers = append(listingNumbers, listingNumber)
	}
	rows.Close()

	return listingNumbers, err
}

func (conn *DatabaseConn) SetMissingFlags(listingNumbers []int) (err error) {

	var UPDATE = `UPDATE notes SET missing='Y' WHERE listing_number = $1`
	for _, listingNumber := range listingNumbers {
		conn.Db.MustExec(UPDATE, listingNumber)
	}
	return err
}

func (conn *DatabaseConn) StoreListings(listings []ListingData) (err error) {

	var QUERY = `SELECT amount_delinquent, amount_funded, amount_participation,
       amount_remaining, bankcard_utilization, borrower_apr, borrower_city,
       borrower_listing_description, borrower_metropolitan_area, borrower_rate,
       borrower_state, credit_lines_last7_years, credit_pull_date, current_credit_lines,
       current_delinquencies, delinquencies_last7_years, delinquencies_over30_days,
       delinquencies_over60_days, delinquencies_over90_days, dti_wprosper_loan,
       effective_yield, employment_status_description, estimated_loss_rate,
       estimated_return, fico_score, first_recorded_credit_line, funding_threshold,
       income_range, income_range_description, income_verifiable, inquires_last6_months,
       installment_balance, investment_type_description, investment_typeid,
       is_homeowner, last_updated_date, lender_indicator, lender_yield,
       listing_amount, listing_category_id, listing_creation_date, listing_end_date,
       listing_monthly_payment, listing_number, listing_start_date,
       listing_status, listing_status_reason, listing_term, listing_title,
       loan_number, loan_origination_date, max_prior_prosper_loan, member_key,
       min_prior_prosper_loan, monthly_debt, months_employed, now_delinquent_derog,
       occupation, oldest_trade_open_date, open_credit_lines, partial_funding_indicator,
       percent_funded, prior_prosper_loan_earliest_pay_off, prior_prosper_loans31dpd,
       prior_propser_loans61dpd, prior_prosper_loans, prior_prosper_loans_active,
       prior_prosper_loans_balance_outstanding, prior_prosper_loans_cycles_billed,
       prior_prosper_loans_late_cycles, prior_prosper_loans_late_payments_one_month_plus,
       prior_prosper_loans_ontime_payments, prior_prosper_loans_principal_borrowed,
       prior_prosper_loans_principal_outstanding, prosper_rating, prosper_score,
       public_records_last10_years, public_records_last12_months, real_estate_balance,
       real_estate_payment, revolving_balance, satisfactory_accounts,
       scorex, scorex_change, sip_offer_id, stated_monthly_income, total_inquiries,
       total_open_revolving_accounts, total_trade_items, verification_stage,
       was_delinquent_derog, whole_loan, whole_loan_end_date, whole_loan_start_date
       FROM listings
       WHERE active='Y' and listing_number=:listing_number`

	var INSERT = `INSERT INTO listings(
            amount_delinquent, amount_funded, amount_participation,
            amount_remaining, bankcard_utilization, borrower_apr, borrower_city,
            borrower_listing_description, borrower_metropolitan_area, borrower_rate,
            borrower_state, credit_lines_last7_years, credit_pull_date, current_credit_lines,
            current_delinquencies, delinquencies_last7_years, delinquencies_over30_days,
            delinquencies_over60_days, delinquencies_over90_days, dti_wprosper_loan,
            effective_yield, employment_status_description, estimated_loss_rate,
            estimated_return, fico_score, first_recorded_credit_line, funding_threshold,
            income_range, income_range_description, income_verifiable, inquires_last6_months,
            installment_balance, investment_type_description, investment_typeid,
            is_homeowner, last_updated_date, lender_indicator, lender_yield,
            listing_amount, listing_category_id, listing_creation_date, listing_end_date,
            listing_monthly_payment, listing_number, listing_start_date,
            listing_status, listing_status_reason, listing_term, listing_title,
            loan_number, loan_origination_date, max_prior_prosper_loan, member_key,
            min_prior_prosper_loan, monthly_debt, months_employed, now_delinquent_derog,
            occupation, oldest_trade_open_date, open_credit_lines, partial_funding_indicator,
            percent_funded, prior_prosper_loan_earliest_pay_off, prior_prosper_loans31dpd,
            prior_propser_loans61dpd, prior_prosper_loans, prior_prosper_loans_active,
            prior_prosper_loans_balance_outstanding, prior_prosper_loans_cycles_billed,
            prior_prosper_loans_late_cycles, prior_prosper_loans_late_payments_one_month_plus,
            prior_prosper_loans_ontime_payments, prior_prosper_loans_principal_borrowed,
            prior_prosper_loans_principal_outstanding, prosper_rating, prosper_score,
            public_records_last10_years, public_records_last12_months, real_estate_balance,
            real_estate_payment, revolving_balance, satisfactory_accounts,
            scorex, scorex_change, sip_offer_id, stated_monthly_income, total_inquiries,
            total_open_revolving_accounts, total_trade_items, verification_stage,
            was_delinquent_derog, whole_loan, whole_loan_end_date, whole_loan_start_date,
            active)
            VALUES(:amount_delinquent, :amount_funded, :amount_participation, :amount_remaining,
            :bankcard_utilization, :borrower_apr, :borrower_city, :borrower_listing_description,
            :borrower_metropolitan_area, :borrower_rate, :borrower_state, :credit_lines_last7_years,
            :credit_pull_date, :current_credit_lines, :current_delinquencies, :delinquencies_last7_years,
            :delinquencies_over30_days, :delinquencies_over60_days, :delinquencies_over90_days,
            :dti_wprosper_loan, :effective_yield, :employment_status_description, :estimated_loss_rate,
            :estimated_return, :fico_score, :first_recorded_credit_line, :funding_threshold, :income_range,
            :income_range_description, :income_verifiable, :inquires_last6_months, :installment_balance,
            :investment_type_description, :investment_typeid, :is_homeowner, :last_updated_date,
            :lender_indicator, :lender_yield, :listing_amount, :listing_category_id, :listing_creation_date,
            :listing_end_date, :listing_monthly_payment, :listing_number, :listing_start_date, :listing_status,
            :listing_status_reason, :listing_term, :listing_title, :loan_number, :loan_origination_date,
            :max_prior_prosper_loan, :member_key, :min_prior_prosper_loan, :monthly_debt, :months_employed,
            :now_delinquent_derog, :occupation, :oldest_trade_open_date, :open_credit_lines,
            :partial_funding_indicator, :percent_funded, :prior_prosper_loan_earliest_pay_off,
            :prior_prosper_loans31dpd, :prior_propser_loans61dpd, :prior_prosper_loans,
            :prior_prosper_loans_active, :prior_prosper_loans_balance_outstanding,
            :prior_prosper_loans_cycles_billed, :prior_prosper_loans_late_cycles,
            :prior_prosper_loans_late_payments_one_month_plus, :prior_prosper_loans_ontime_payments,
            :prior_prosper_loans_principal_borrowed, :prior_prosper_loans_principal_outstanding,
            :prosper_rating, :prosper_score, :public_records_last10_years, :public_records_last12_months,
            :real_estate_balance, :real_estate_payment, :revolving_balance, :satisfactory_accounts, :scorex,
            :scorex_change, :sip_offer_id, :stated_monthly_income, :total_inquiries,
            :total_open_revolving_accounts, :total_trade_items, :verification_stage, :was_delinquent_derog,
            :whole_loan, :whole_loan_end_date, :whole_loan_start_date, 'Y')`

	var UPDATE = `UPDATE listings SET active='N' WHERE active='Y' AND listing_number=:listing_number`

	tx := conn.Db.MustBegin()

	for _, listing := range listings {

		rows, err := tx.NamedQuery(QUERY, listing)
		if err != nil {
			return err
		}

		dbListing := ListingData{}
		for rows.Next() {
			err := rows.StructScan(&dbListing)
			if err != nil {
				return err
			}

			rows.Close()
			break
		}

		if dbListing != listing {
			_, err = tx.NamedExec(UPDATE, listing)
			if err != nil {
				return err
			}
			_, err = tx.NamedExec(INSERT, listing)
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()

	return err
}
