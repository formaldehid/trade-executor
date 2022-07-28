package trader

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type db struct {
	DB *sql.DB
}

type DB interface {
	CloseDB()
	InsertTicker(params *Ticker) error
	InsertOrder(params *Order) error
	InsertOrderFill(params *OrderFill) error
	FinalizeOrder(o *Order) error
}

func NewDB(cfg *Config) (DB, error) {
	sqlDb, err := sql.Open("sqlite3", cfg.DBDataSourceName)
	if err != nil {
		return nil, err
	}

	err = createTables(sqlDb)
	if err != nil {
		return nil, err
	}

	db := &db{
		DB: sqlDb,
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	const create string = `
	  	CREATE TABLE IF NOT EXISTS tickers (
			id INTEGER NOT NULL PRIMARY KEY,
			symbol TEXT NOT NULL,
			best_bid_price REAL NOT NULL,
			best_bid_qty REAL NOT NULL,
			best_ask_price REAL NOT NULL,
			best_ask_qty REAL NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders (
			id TEXT NOT NULL PRIMARY KEY,
			type TEXT NOT NULL,
			symbol TEXT NOT NULL,
			size REAL NOT NULL,
			price REAL NOT NULL,
			is_finalized BOOLEAN DEFAULT false NOT NULL,
			final_price REAL DEFAULT 0 NOT NULL
		);

		CREATE TABLE IF NOT EXISTS order_fills (
			id TEXT NOT NULL PRIMARY KEY,
			order_id TEXT NOT NULL,
			ticker_id INTEGER NOT NULL,
			created_at TEXT NOT NULL,
			size REAL NOT NULL,
			price REAL NOT NULL,
			FOREIGN KEY(order_id) REFERENCES orders(id),
			FOREIGN KEY(ticker_id) REFERENCES tickers(id)
		);

		DELETE FROM tickers;
		DELETE FROM orders;
		DELETE FROM order_fills;
	`

	_, err := db.Exec(create)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) CloseDB() {
	err := db.DB.Close()
	if err != nil {
		log.Error(err)
	}
}

func (db *db) InsertTicker(t *Ticker) error {
	const i string = "INSERT INTO tickers VALUES(?,?,?,?,?,?);"

	_, err := db.DB.Exec(i, t.Id, t.Symbol, t.BestBidPrice, t.BestBidQty, t.BestAskPrice, t.BestAskQty)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) InsertOrder(o *Order) error {
	const i string = "INSERT INTO orders VALUES(?,?,?,?,?,false,0);"

	_, err := db.DB.Exec(i, o.Id, o.Type, o.Symbol, o.Size, o.Price)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) InsertOrderFill(of *OrderFill) error {
	const i string = "INSERT INTO order_fills VALUES(?,?,?,?,?,?);"

	_, err := db.DB.Exec(i, of.Id, of.Order.Id, of.Ticker.Id, of.CreatedAt, of.Size, of.Price)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) FinalizeOrder(o *Order) error {
	const u string = `
		UPDATE orders
			SET is_finalized = true,
			    final_price = (SELECT SUM(price) FROM order_fills WHERE order_id = ?)
			WHERE id = ?
	`

	_, err := db.DB.Exec(u, o.Id, o.Id)
	if err != nil {
		return err
	}
	return nil
}
