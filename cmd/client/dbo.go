package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDBO() *sql.DB {
	db, err := sql.Open("sqlite3", "./deep_lairs.db")
	if err != nil {
		panic(err.Error())
	}
	// create users table if it doesn't exist
	_, err = initUsers(db)
	if err != nil {
		panic(err.Error())
	}
	// create characters table if it doesn't exist
	_, err = initCharacters(db)
	if err != nil {
		panic(err.Error())
	}
	// create entity_fightables table if it doesn't exist
	_, err = initEntityFightables(db)
	if err != nil {
		panic(err.Error())
	}
	// create fightables table if it doesn't exist
	_, err = initFightables(db)
	if err != nil {
		panic(err.Error())
	}
	// create lu_bonus_types table if it doesn't exist
	_, err = initBonusTypes(db)
	if err != nil {
		panic(err.Error())
	}
	// create lu_slots table if it doesn't exist
	_, err = initSlots(db)
	if err != nil {
		panic(err.Error())
	}
	// create items table if it doesn't exist
	_, err = initItems(db)
	if err != nil {
		panic(err.Error())
	}
	// create item_states table if it doesn't exist
	_, err = initItemStates(db)
	if err != nil {
		panic(err.Error())
	}
	// return db connection
	return db
}

func initUsers(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY NOT NULL DEFAULT (
				substr(hex(randomblob(16)),1,8) || '-' ||
		substr(hex(randomblob(16)),9,4) || '-' ||
		substr(hex(randomblob(16)),13,4) || '-' ||
		substr(hex(randomblob(16)),17,4) || '-' ||
		substr(hex(randomblob(16)),21,12)
		),
		username TEXT UNIQUE,
		password TEXT,
		email TEXT UNIQUE,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	return db.Exec(sql)
}

func initCharacters(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS characters (
		id TEXT PRIMARY KEY NOT NULL DEFAULT (
				substr(hex(randomblob(16)),1,8) || '-' ||
		substr(hex(randomblob(16)),9,4) || '-' ||
		substr(hex(randomblob(16)),13,4) || '-' ||
		substr(hex(randomblob(16)),17,4) || '-' ||
		substr(hex(randomblob(16)),21,12)
		),
		user_id TEXT NOT NULL,
		name TEXT,
		last_name TEXT,
		location_id TEXT,
		xp INTEGER,
		max_xp INTEGER,
		level INTEGER,
		class TEXT,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) NOT NULL
	);
	`
	return db.Exec(sql)
}

func initEntityFightables(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS entity_fightables (
		id TEXT PRIMARY KEY DEFAULT (
				substr(hex(randomblob(16)),1,8) || '-' ||
		substr(hex(randomblob(16)),9,4) || '-' ||
		substr(hex(randomblob(16)),13,4) || '-' ||
		substr(hex(randomblob(16)),17,4) || '-' ||
		substr(hex(randomblob(16)),21,12)
		),
		entity_id TEXT NOT NULL,
		fightable_id TEXT NOT NULL,
		entity_type_id INTEGER NOT NULL,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	return db.Exec(sql)
}

func initFightables(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS fightables (
		id TEXT PRIMARY KEY DEFAULT (
				substr(hex(randomblob(16)),1,8) || '-' ||
		substr(hex(randomblob(16)),9,4) || '-' ||
		substr(hex(randomblob(16)),13,4) || '-' ||
		substr(hex(randomblob(16)),17,4) || '-' ||
		substr(hex(randomblob(16)),21,12)
		),
		entity_fightable_id TEXT NOT NULL,
		health INTEGER,
		base_max_health INTEGER,
		attack INTEGER,
		base_attack INTEGER,
		defense INTEGER,
		base_defense INTEGER,
		mana INTEGER,
		base_max_mana INTEGER,
		stamina INTEGER,
		base_max_stamina INTEGER,
		speed INTEGER,
		base_speed INTEGER,
		intelligence INTEGER,
		base_intelligence INTEGER,
		image TEXT,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(entity_fightable_id) REFERENCES entity_fightables(id) NOT NULL
	);
	`
	return db.Exec(sql)
}

func initBonusTypes(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS lu_bonus_types (
		id INTEGER PRIMARY KEY,
		name TEXT,
		description TEXT,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	return db.Exec(sql)
}

func initSlots(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS lu_slots (
		id INTEGER PRIMARY KEY,
		name TEXT,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	return db.Exec(sql)
}

func initItems(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS items (
		id TEXT PRIMARY KEY DEFAULT (
				substr(hex(randomblob(16)),1,8) || '-' ||
		substr(hex(randomblob(16)),9,4) || '-' ||
		substr(hex(randomblob(16)),13,4) || '-' ||
		substr(hex(randomblob(16)),17,4) || '-' ||
		substr(hex(randomblob(16)),21,12)
		),
		name TEXT,
		description TEXT,
		bonus_type INTEGER NOT NULL,
		bonus_amount INTEGER,
		value INTEGER,
		slot_id INTEGER NOT NULL,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(bonus_type) REFERENCES lu_bonus_types(id) NOT NULL
	);
	`
	return db.Exec(sql)
}

func initItemStates(db *sql.DB) (sql.Result, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS item_states (
		id TEXT PRIMARY KEY DEFAULT (
				substr(hex(randomblob(16)),1,8) || '-' ||
		substr(hex(randomblob(16)),9,4) || '-' ||
		substr(hex(randomblob(16)),13,4) || '-' ||
		substr(hex(randomblob(16)),17,4) || '-' ||
		substr(hex(randomblob(16)),21,12)
		),
		fightable_id TEXT NOT NULL,
		item_id TEXT NOT NULL,
		equipped BOOLEAN,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_modified DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(fightable_id) REFERENCES fightables(id) NOT NULL,
		FOREIGN KEY(item_id) REFERENCES items(id) NOT NULL
	);
	`
	return db.Exec(sql)
}
