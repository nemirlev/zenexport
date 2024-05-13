package clickhouse

import (
	"context"
	"fmt"
	"github.com/nemirlev/zenapi"
)

func (s *Store) Save(data *zenapi.Response) error {
	ctx := context.Background()

	defer func() {
		if err := s.Conn.Close(); err != nil {
			s.Log.WithError(err, "failed to close connection")
		}
	}()

	saveFuncs := []func(ctx context.Context, data *zenapi.Response) error{
		func(ctx context.Context, data *zenapi.Response) error { return s.saveInstruments(ctx, data.Instrument) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveCountries(ctx, data.Country) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveCompanies(ctx, data.Company) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveUsers(ctx, data.User) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveAccounts(ctx, data.Account) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveTags(ctx, data.Tag) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveMerchants(ctx, data.Merchant) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveBudgets(ctx, data.Budget) },
		func(ctx context.Context, data *zenapi.Response) error { return s.saveReminders(ctx, data.Reminder) },
		func(ctx context.Context, data *zenapi.Response) error {
			return s.saveReminderMarkers(ctx, data.ReminderMarker)
		},
		func(ctx context.Context, data *zenapi.Response) error {
			return s.saveTransactions(ctx, data.Transaction)
		},
	}

	for _, saveFunc := range saveFuncs {
		if err := saveFunc(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) saveBatch(ctx context.Context, tableName string, query string, data [][]interface{}) error {
	fmt.Printf("Starting to save data into %s...\n", tableName)
	if err := s.truncateTable(ctx, tableName); err != nil {
		s.Log.WithError(err, "failed to truncate table %s", tableName)
		return err
	}

	if err := s.executeBatch(ctx, query, data); err != nil {
		s.Log.WithError(err, "failed to execute batch for table %s", tableName)
		return err
	}
	fmt.Printf("Finished saving data into %s.\n", tableName)
	return nil
}

func (s *Store) saveTransactions(ctx context.Context, transactions []zenapi.Transaction) error {
	query := `
		INSERT INTO transaction (
			id, changed, created, user, deleted, hold, income_instrument, income_account, 
			income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, 
			original_payee, comment, date, mcc, reminder_marker, op_income, op_income_instrument, 
			op_outcome, op_outcome_instrument, latitude, longitude
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

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

	return s.saveBatch(ctx, "transaction", query, data)
}

func (s *Store) saveReminderMarkers(ctx context.Context, markers []zenapi.ReminderMarker) error {
	query := `
		INSERT INTO reminder_marker (
			id, changed, user, income_instrument, income_account, income, outcome_instrument, 
			outcome_account, outcome, tag, merchant, payee, comment, date, reminder, state, notify
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, marker := range markers {
		data = append(data, []interface{}{
			marker.ID, marker.Changed, marker.User, marker.IncomeInstrument, marker.IncomeAccount,
			marker.Income, marker.OutcomeInstrument, marker.OutcomeAccount, marker.Outcome, marker.Tag,
			marker.Merchant, marker.Payee, marker.Comment, marker.Date, marker.Reminder,
			marker.State, marker.Notify,
		})
	}

	return s.saveBatch(ctx, "reminder_marker", query, data)
}

func (s *Store) saveReminders(ctx context.Context, reminders []zenapi.Reminder) error {
	query := `
		INSERT INTO reminder (
			id, changed, user, income_instrument, income_account, income, outcome_instrument, 
			outcome_account, outcome, tag, merchant, payee, comment, interval, step, points, 
			start_date, end_date, notify
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, reminder := range reminders {
		data = append(data, []interface{}{
			reminder.ID, reminder.Changed, reminder.User, reminder.IncomeInstrument, reminder.IncomeAccount,
			reminder.Income, reminder.OutcomeInstrument, reminder.OutcomeAccount, reminder.Outcome,
			reminder.Tag, reminder.Merchant, reminder.Payee, reminder.Comment, reminder.Interval, reminder.Step,
			reminder.Points, reminder.StartDate, reminder.EndDate, reminder.Notify,
		})
	}

	return s.saveBatch(ctx, "reminder", query, data)
}

func (s *Store) saveBudgets(ctx context.Context, budgets []zenapi.Budget) error {
	query := `
		INSERT INTO budget (
			changed, user, tag, date, income, income_lock, outcome, outcome_lock
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, budget := range budgets {
		data = append(data, []interface{}{
			budget.Changed, budget.User, budget.Tag, budget.Date,
			budget.Income, budget.IncomeLock, budget.Outcome, budget.OutcomeLock,
		})
	}

	return s.saveBatch(ctx, "budget", query, data)
}

func (s *Store) saveMerchants(ctx context.Context, merchants []zenapi.Merchant) error {
	query := `
		INSERT INTO merchant (
			id, changed, user, title
		) VALUES (
			?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, merchant := range merchants {
		data = append(data, []interface{}{
			merchant.ID, merchant.Changed, merchant.User, merchant.Title,
		})
	}

	return s.saveBatch(ctx, "merchant", query, data)
}

func (s *Store) saveTags(ctx context.Context, tags []zenapi.Tag) error {
	query := `
		INSERT INTO tag (
			id, changed, user, title, parent, icon, picture, color, show_income, 
			show_outcome, budget_income, budget_outcome, required
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, tag := range tags {
		data = append(data, []interface{}{
			tag.ID, tag.Changed, tag.User, tag.Title, tag.Parent, tag.Icon,
			tag.Picture, tag.Color, tag.ShowIncome, tag.ShowOutcome,
			tag.BudgetIncome, tag.BudgetOutcome, tag.Required,
		})
	}

	return s.saveBatch(ctx, "tag", query, data)
}

func (s *Store) saveAccounts(ctx context.Context, accounts []zenapi.Account) error {
	query := `
		INSERT INTO account (
			id, changed, user, role, instrument, company, type, title, sync_id, balance, 
			start_balance, credit_limit, in_balance, savings, enable_correction, enable_sms, 
			archive, capitalization, percent, start_date, end_date_offset, 
			end_date_offset_interval, payoff_step, payoff_interval
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

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

	return s.saveBatch(ctx, "account", query, data)
}

func (s *Store) saveUsers(ctx context.Context, users []zenapi.User) error {
	query := `
		INSERT INTO user (
			id, changed, login, currency, parent
		) VALUES (
			?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, user := range users {
		data = append(data, []interface{}{
			user.ID, user.Changed, user.Login, user.Currency, user.Parent,
		})
	}

	return s.saveBatch(ctx, "user", query, data)
}

func (s *Store) saveCompanies(ctx context.Context, companies []zenapi.Company) error {
	query := `
		INSERT INTO company (
			id, changed, title, full_title, www, country
		) VALUES (
			?, ?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, company := range companies {
		data = append(data, []interface{}{
			company.ID, company.Changed, company.Title, company.FullTitle, company.Www, company.Country,
		})
	}

	return s.saveBatch(ctx, "company", query, data)
}

func (s *Store) saveCountries(ctx context.Context, countries []zenapi.Country) error {
	query := `
		INSERT INTO country (
			id, title, currency, domain
		) VALUES (
			?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, country := range countries {
		data = append(data, []interface{}{
			country.ID, country.Title, country.Currency, country.Domain,
		})
	}

	return s.saveBatch(ctx, "country", query, data)
}

func (s *Store) saveInstruments(ctx context.Context, instruments []zenapi.Instrument) error {
	query := `
		INSERT INTO instrument (
			id, changed, title, short_title, symbol, rate
		) VALUES (
			?, ?, ?, ?, ?, ?
		)
	`

	var data [][]interface{}
	for _, instrument := range instruments {
		data = append(data, []interface{}{
			instrument.ID, instrument.Changed, instrument.Title, instrument.ShortTitle, instrument.Symbol, instrument.Rate,
		})
	}

	return s.saveBatch(ctx, "instrument", query, data)
}
