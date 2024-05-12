package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/internal/config"
	"log"
)

func (s *Store) Save(cfg *config.Config, data *zenapi.Response) error {
	err := s.connect(cfg)
	if err != nil {
		return err
	}
	defer func() {
		err := s.Conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}()

	saveInstruments(s.Conn, data.Instrument)
	saveCountries(s.Conn, data.Country)
	saveCompanies(s.Conn, data.Company)
	saveUsers(s.Conn, data.User)
	saveAccounts(s.Conn, data.Account)
	saveTags(s.Conn, data.Tag)
	saveMerchants(s.Conn, data.Merchant)
	saveBudgets(s.Conn, data.Budget)
	saveReminders(s.Conn, data.Reminder)
	saveReminderMarkers(s.Conn, data.ReminderMarker)
	saveTransactions(s.Conn, data.Transaction)

	return nil
}

func saveTransactions(conn driver.Conn, transactions []zenapi.Transaction) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS transaction"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, transaction := range transactions {
		data = append(data, []interface{}{
			transaction.ID, transaction.Changed, transaction.Created, transaction.User, transaction.Deleted,
			transaction.Hold, transaction.IncomeInstrument, transaction.IncomeAccount, transaction.Income,
			transaction.OutcomeInstrument, transaction.OutcomeAccount, transaction.Outcome, transaction.Tag,
			transaction.Merchant, transaction.Payee, transaction.OriginalPayee, transaction.Comment,
			transaction.Date, transaction.Mcc, transaction.ReminderMarker, transaction.OpIncome,
			transaction.OpIncomeInstrument, transaction.OpOutcome, transaction.OpOutcomeInstrument,
			transaction.Latitude, transaction.Longitude,
		})
	}

	query := "INSERT INTO transaction (id, changed, created, user, deleted, hold, income_instrument, income_account, income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, original_payee, comment, date, mcc, reminder_marker, op_income, op_income_instrument, op_outcome, op_outcome_instrument, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveReminderMarkers(conn driver.Conn, markers []zenapi.ReminderMarker) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS reminder_marker"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, marker := range markers {
		data = append(data, []interface{}{
			marker.ID, marker.Changed, marker.User, marker.IncomeInstrument, marker.IncomeAccount,
			marker.Income, marker.OutcomeInstrument, marker.OutcomeAccount, marker.Outcome, marker.Tag,
			marker.Merchant, marker.Payee, marker.Comment, marker.Date, marker.Reminder,
			marker.State, marker.Notify,
		})
	}

	query := "INSERT INTO reminder_marker (id, changed, user, income_instrument, income_account, income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, comment, date, reminder, state, notify) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveReminders(conn driver.Conn, reminders []zenapi.Reminder) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS reminder"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, reminder := range reminders {
		data = append(data, []interface{}{
			reminder.ID, reminder.Changed, reminder.User, reminder.IncomeInstrument, reminder.IncomeAccount,
			reminder.Income, reminder.OutcomeInstrument, reminder.OutcomeAccount, reminder.Outcome,
			reminder.Tag, reminder.Merchant, reminder.Payee, reminder.Comment, reminder.Interval, reminder.Step,
			reminder.Points, reminder.StartDate, reminder.EndDate, reminder.Notify,
		})
	}

	query := "INSERT INTO reminder (id, changed, user, income_instrument, income_account, income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, comment, interval, step, points, start_date, end_date, notify) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveBudgets(conn driver.Conn, budgets []zenapi.Budget) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS budget"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, budget := range budgets {
		data = append(data, []interface{}{
			budget.Changed, budget.User, budget.Tag, budget.Date,
			budget.Income, budget.IncomeLock, budget.Outcome, budget.OutcomeLock,
		})
	}

	query := "INSERT INTO budget (changed, user, tag, date, income, income_lock, outcome, outcome_lock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveMerchants(conn driver.Conn, merchants []zenapi.Merchant) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS merchant"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, merchant := range merchants {
		data = append(data, []interface{}{
			merchant.ID, merchant.Changed, merchant.User, merchant.Title,
		})
	}

	query := "INSERT INTO merchant (id, changed, user, title) VALUES (?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveTags(conn driver.Conn, tags []zenapi.Tag) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS tag"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, tag := range tags {
		data = append(data, []interface{}{
			tag.ID, tag.Changed, tag.User, tag.Title, tag.Parent, tag.Icon,
			tag.Picture, tag.Color, tag.ShowIncome, tag.ShowOutcome,
			tag.BudgetIncome, tag.BudgetOutcome, tag.Required,
		})
	}

	query := "INSERT INTO tag (id, changed, user, title, parent, icon, picture, color, show_income, show_outcome, budget_income, budget_outcome, required) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveAccounts(conn driver.Conn, accounts []zenapi.Account) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS account"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, account := range accounts {
		data = append(data, []interface{}{
			account.ID, account.Changed, account.User, account.Role, account.Instrument, account.Company,
			account.Type, account.Title, account.SyncID, account.Balance, account.StartBalance, account.CreditLimit,
			account.InBalance, account.Savings, account.EnableCorrection, account.EnableSMS, account.Archive,
			account.Capitalization, account.Percent, account.StartDate, account.EndDateOffset,
			account.EndDateOffsetInterval, account.PayoffStep, account.PayoffInterval,
		})
	}

	query := "INSERT INTO account (id, changed, user, role, instrument, company, type, title, sync_id, balance, start_balance, credit_limit, in_balance, savings, enable_correction, enable_sms, archive, capitalization, percent, start_date, end_date_offset, end_date_offset_interval, payoff_step, payoff_interval) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveUsers(conn driver.Conn, users []zenapi.User) {
	ctx := context.Background()

	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS user"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, user := range users {
		data = append(data, []interface{}{
			user.ID, user.Changed, user.Login, user.Currency, user.Parent,
		})
	}

	if err := executeBatch(conn, ctx, "INSERT INTO user (id, changed, login, currency, parent) VALUES (?, ?, ?, ?, ?)", data); err != nil {
		log.Fatal(err)
	}
}

func saveCompanies(conn driver.Conn, companies []zenapi.Company) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS company"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, company := range companies {
		data = append(data, []interface{}{
			company.ID, company.Changed, company.Title, company.FullTitle, company.Www, company.Country,
		})
	}

	query := "INSERT INTO company (id, changed, title, full_title, www, country) VALUES (?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveCountries(conn driver.Conn, countries []zenapi.Country) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS country"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, country := range countries {
		data = append(data, []interface{}{
			country.ID, country.Title, country.Currency, country.Domain,
		})
	}

	query := "INSERT INTO country (id, title, currency, domain) VALUES (?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}

func saveInstruments(conn driver.Conn, instruments []zenapi.Instrument) {
	ctx := context.Background()
	if err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS instrument"); err != nil {
		log.Fatal(err)
	}

	var data [][]interface{}
	for _, instrument := range instruments {
		data = append(data, []interface{}{
			instrument.ID, instrument.Changed, instrument.Title, instrument.ShortTitle, instrument.Symbol, instrument.Rate,
		})
	}

	query := "INSERT INTO instrument (id, changed, title, short_title, symbol, rate) VALUES (?, ?, ?, ?, ?, ?)"
	if err := executeBatch(conn, ctx, query, data); err != nil {
		log.Fatal(err)
	}
}