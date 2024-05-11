package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/internal/config"
	"log"
)

type ClickHouse struct {
	Conn driver.Conn
}

func (c *ClickHouse) connect(cfg *config.Config) error {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:9000", cfg.ClickhouseServer)},
			Auth: clickhouse.Auth{
				Database: cfg.ClickhouseDB,
				Username: cfg.ClickhouseUser,
				Password: cfg.ClickhousePassword,
			},
			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return err
	}

	if err := conn.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}

	c.Conn = conn
	return nil
}

func (c *ClickHouse) Save(cfg *config.Config, data *zenapi.Response) error {
	err := c.connect(cfg)
	if err != nil {
		return err
	}
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}()

	saveInstruments(c.Conn, data.Instrument)
	saveCountries(c.Conn, data.Country)
	saveCompanies(c.Conn, data.Company)
	saveUsers(c.Conn, data.User)
	saveAccounts(c.Conn, data.Account)
	saveTags(c.Conn, data.Tag)
	saveMerchants(c.Conn, data.Merchant)
	saveBudgets(c.Conn, data.Budget)
	saveReminders(c.Conn, data.Reminder)
	saveReminderMarkers(c.Conn, data.ReminderMarker)
	saveTransactions(c.Conn, data.Transaction)

	return nil
}

//func saveBatch(conn driver.Conn, entity db.DatabaseEntity) error {
//	ctx := context.Background()
//
//	str := "TRUNCATE TABLE IF EXISTS " + entity.GetTableName()
//	err := conn.Exec(ctx, str)
//	if err != nil {
//		return err
//	}
//
//	batch, err := conn.PrepareBatch(ctx, entity.GetInsertQuery())
//	if err != nil {
//		return err
//	}
//
//	err = batch.Append(entity.GetValues()...)
//	if err != nil {
//		return err
//	}
//
//	if err := batch.Send(); err != nil {
//		return err
//	}
//
//	return nil
//}

func saveTransactions(conn driver.Conn, transactions []zenapi.Transaction) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS transaction")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO transaction (id, changed, created, user, deleted, hold, income_instrument, income_account, income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, original_payee, comment, date, mcc, reminder_marker, op_income, op_income_instrument, op_outcome, op_outcome_instrument, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, transaction := range transactions {
		err := batch.Append(
			transaction.ID,
			transaction.Changed,
			transaction.Created,
			transaction.User,
			transaction.Deleted,
			transaction.Hold,
			transaction.IncomeInstrument,
			transaction.IncomeAccount,
			transaction.Income,
			transaction.OutcomeInstrument,
			transaction.OutcomeAccount,
			transaction.Outcome,
			transaction.Tag,
			transaction.Merchant,
			transaction.Payee,
			transaction.OriginalPayee,
			transaction.Comment,
			transaction.Date,
			transaction.Mcc,
			transaction.ReminderMarker,
			transaction.OpIncome,
			transaction.OpIncomeInstrument,
			transaction.OpOutcome,
			transaction.OpOutcomeInstrument,
			transaction.Latitude,
			transaction.Longitude,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveReminderMarkers(conn driver.Conn, markers []zenapi.ReminderMarker) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS reminder_marker")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO reminder_marker (id, changed, user, income_instrument, income_account, income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, comment, date, reminder, state, notify) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, marker := range markers {
		err := batch.Append(
			marker.ID,
			marker.Changed,
			marker.User,
			marker.IncomeInstrument,
			marker.IncomeAccount,
			marker.Income,
			marker.OutcomeInstrument,
			marker.OutcomeAccount,
			marker.Outcome,
			marker.Tag,
			marker.Merchant,
			marker.Payee,
			marker.Comment,
			marker.Date,
			marker.Reminder,
			marker.State,
			marker.Notify,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveReminders(conn driver.Conn, reminders []zenapi.Reminder) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS reminder")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO reminder (id, changed, user, income_instrument, income_account, income, outcome_instrument, outcome_account, outcome, tag, merchant, payee, comment, interval, step, points, start_date, end_date, notify) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, reminder := range reminders {
		err := batch.Append(
			reminder.ID,
			reminder.Changed,
			reminder.User,
			reminder.IncomeInstrument,
			reminder.IncomeAccount,
			reminder.Income,
			reminder.OutcomeInstrument,
			reminder.OutcomeAccount,
			reminder.Outcome,
			reminder.Tag,
			reminder.Merchant,
			reminder.Payee,
			reminder.Comment,
			reminder.Interval,
			reminder.Step,
			reminder.Points,
			reminder.StartDate,
			reminder.EndDate,
			reminder.Notify,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveBudgets(conn driver.Conn, budgets []zenapi.Budget) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS budget")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO budget (changed, user, tag, date, income, income_lock, outcome, outcome_lock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, budget := range budgets {
		err := batch.Append(
			budget.Changed,
			budget.User,
			budget.Tag,
			budget.Date,
			budget.Income,
			budget.IncomeLock,
			budget.Outcome,
			budget.OutcomeLock,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveMerchants(conn driver.Conn, merchants []zenapi.Merchant) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS merchant")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO merchant (id, changed, user, title) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, merchant := range merchants {
		err := batch.Append(
			merchant.ID,
			merchant.Changed,
			merchant.User,
			merchant.Title,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveTags(conn driver.Conn, tags []zenapi.Tag) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS tag")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO tag (id, changed, user, title, parent, icon, picture, color, show_income, show_outcome, budget_income, budget_outcome, required) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, tag := range tags {
		err := batch.Append(
			tag.ID,
			tag.Changed,
			tag.User,
			tag.Title,
			tag.Parent,
			tag.Icon,
			tag.Picture,
			tag.Color,
			tag.ShowIncome,
			tag.ShowOutcome,
			tag.BudgetIncome,
			tag.BudgetOutcome,
			tag.Required,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveAccounts(conn driver.Conn, accounts []zenapi.Account) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS account")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO account (id, changed, user, role, instrument, company, type, title, sync_id, balance, start_balance, credit_limit, in_balance, savings, enable_correction, enable_sms, archive, capitalization, percent, start_date, end_date_offset, end_date_offset_interval, payoff_step, payoff_interval) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range accounts {
		err := batch.Append(
			account.ID,
			account.Changed,
			account.User,
			account.Role,
			account.Instrument,
			account.Company,
			account.Type,
			account.Title,
			account.SyncID,
			account.Balance,
			account.StartBalance,
			account.CreditLimit,
			account.InBalance,
			account.Savings,
			account.EnableCorrection,
			account.EnableSMS,
			account.Archive,
			account.Capitalization,
			account.Percent,
			account.StartDate,
			account.EndDateOffset,
			account.EndDateOffsetInterval,
			account.PayoffStep,
			account.PayoffInterval,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveUsers(conn driver.Conn, users []zenapi.User) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS user")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO user (id, changed, login, currency, parent) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		err := batch.Append(
			user.ID,
			user.Changed,
			user.Login,
			user.Currency,
			user.Parent,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveCompanies(conn driver.Conn, companies []zenapi.Company) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS company")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO company (id, changed, title, full_title, www, country) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, company := range companies {
		err := batch.Append(
			company.ID,
			company.Changed,
			company.Title,
			company.FullTitle,
			company.Www,
			company.Country,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveCountries(conn driver.Conn, countries []zenapi.Country) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS country")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO country (id, title, currency, domain) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, country := range countries {
		err := batch.Append(
			country.ID,
			country.Title,
			country.Currency,
			country.Domain,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func saveInstruments(conn driver.Conn, instruments []zenapi.Instrument) {
	ctx := context.Background()

	err := conn.Exec(ctx, "TRUNCATE TABLE IF EXISTS instrument")
	if err != nil {
		log.Fatal(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO instrument (id, changed, title, short_title, symbol, rate) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, instrument := range instruments {
		err := batch.Append(
			instrument.ID,
			instrument.Changed,
			instrument.Title,
			instrument.ShortTitle,
			instrument.Symbol,
			instrument.Rate,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}
