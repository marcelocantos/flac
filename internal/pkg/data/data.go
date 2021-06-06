package data

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Database struct {
	db *sql.DB

	// Read
	maxScoreStmt  *sql.Stmt
	maxPosStmt    *sql.Stmt
	wordScoreStmt *sql.Stmt
	wordPosStmt   *sql.Stmt
	wordAtStmt    *sql.Stmt

	// Write
	enqueueWordStmt  *sql.Stmt
	updateScoreStmt  *sql.Stmt
	rotateWords1Stmt *sql.Stmt
	rotateWords2Stmt *sql.Stmt
}

func NewDatabase(path string) (*Database, error) {
	var d Database
	var err error
	d.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	for _, def := range []string{
		`word_score (
			word TEXT PRIMARY KEY,
			score INTEGER
		)`,
		`queue (
			pos INT PRIMARY KEY,
			word TEXT UNIQUE
		)`,
	} {
		_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS " + def)
		if err != nil {
			return nil, err
		}
	}

	for stmt, query := range map[**sql.Stmt]string{
		&d.maxScoreStmt:    `SELECT COALESCE(MAX(score), -1) FROM word_score`,
		&d.maxPosStmt:      `SELECT COALESCE(MAX(pos), -1) FROM queue`,
		&d.wordScoreStmt:   `SELECT score FROM word_score WHERE word = ?`,
		&d.wordPosStmt:     `SELECT pos FROM queue WHERE word = ?`,
		&d.wordAtStmt:      `SELECT word FROM queue WHERE pos = ?`,
		&d.enqueueWordStmt: `INSERT INTO queue (pos, word) VALUES (?, ?)`,
		&d.updateScoreStmt: `INSERT OR REPLACE INTO word_score (word, score) VALUES (?, ?)`,
		&d.rotateWords1Stmt: `
			UPDATE queue
			SET pos = -1-((pos - $first + $count + $offset) % $count + $first)
			WHERE pos BETWEEN $first AND $first + $count - 1
		`,
		&d.rotateWords2Stmt: `
			UPDATE queue
			SET pos = -1-pos
			WHERE pos < 0
		`,
	} {
		*stmt, err = d.db.Prepare(query)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) Populate(words []string) error {
	elideRE := regexp.MustCompile(`\P{Han}`)

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	index, err := d.maxPos(tx)
	if err != nil {
		return err
	}

	getWordPos := d.wordPos(tx)
	enqueueWord := tx.Stmt(d.enqueueWordStmt)
	for _, word := range words {
		if elideRE.MatchString(word) {
			continue
		}
		_, err := getWordPos(word)
		if err != nil {
			if _, is := err.(ErrNotFound); !is {
				return err
			}
			index++
			enqueueWord.Exec(index, word)
		}
	}
	return nil
}

func (d *Database) MaxScore() (int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Commit()

	return d.maxScore(tx)
}

func (d *Database) maxScore(tx *sql.Tx) (int, error) {
	var score sql.NullInt64
	if err := tx.Stmt(d.maxScoreStmt).QueryRow().Scan(&score); err != nil {
		return 0, err
	}
	return int(score.Int64), nil
}

func (d *Database) MaxPos() (int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Commit()

	return d.maxPos(tx)
}

func (d *Database) maxPos(tx *sql.Tx) (int, error) {
	var pos sql.NullInt64
	if err := tx.Stmt(d.maxPosStmt).QueryRow().Scan(&pos); err != nil {
		return 0, err
	}
	return int(pos.Int64), nil
}

func (d *Database) WordScore(word string) (int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Commit()

	return d.wordScore(tx)(word)
}

func (d *Database) wordScore(tx *sql.Tx) func(word string) (int, error) {
	getWordScoreStmt := tx.Stmt(d.wordScoreStmt)
	return func(word string) (int, error) {
		return d.selectInt(getWordScoreStmt, "%s: not found in word_score", word)
	}
}

func (d *Database) WordPos(word string) (int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Commit()

	return d.wordPos(tx)(word)
}

func (d *Database) wordPos(tx *sql.Tx) func(word string) (int, error) {
	getWordPosStmt := tx.Stmt(d.wordPosStmt)
	return func(word string) (int, error) {
		return d.selectInt(getWordPosStmt, "%s: not found in queue", word)
	}
}

func (d *Database) selectInt(stmt *sql.Stmt, format string, args ...interface{}) (int, error) {
	var pos sql.NullInt64
	err := stmt.QueryRow(args...).Scan(&pos)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound(errors.Errorf(format, args...))
		}
		return 0, err
	}
	return int(pos.Int64), err
}

func (d *Database) WordAt(pos int) (string, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Commit()

	return d.wordAt(tx)(pos)
}

func (d *Database) wordAt(tx *sql.Tx) func(pos int) (string, error) {
	getWordAtStmt := tx.Stmt(d.wordAtStmt)
	return func(pos int) (string, error) {
		var word string
		err := getWordAtStmt.QueryRow(pos).Scan(&word)
		if err != nil {
			if err == sql.ErrNoRows {
				err = ErrNotFound(fmt.Errorf("no word at index %d", pos))
			}
			return "", err
		}
		return word, err
	}
}

func (d *Database) HeadWord() (string, error) {
	return d.WordAt(0)
}

func (d *Database) UpdateScoreAndPos(word string, score, dest int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	if _, err := tx.Stmt(d.updateScoreStmt).Exec(word, score); err != nil {
		return err
	}
	if dest >= 0 {
		return d.moveWord(tx, word, dest)
	}
	return nil
}

func (d *Database) MoveWord(word string, dest int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	return d.moveWord(tx, word, dest)
}

func (d *Database) moveWord(tx *sql.Tx, word string, dest int) error {
	max, err := d.maxPos(tx)
	if err != nil {
		return err
	}
	if dest < 0 {
		dest = 0
	} else if dest > max {
		dest = max
	}

	src, err := d.wordPos(tx)(word)
	if err != nil {
		return err
	}

	first, last := src, dest
	if first > last {
		first, last = last, first
	}

	offset := 1
	if src < dest {
		offset = -1
	}

	_, err = tx.Stmt(d.rotateWords1Stmt).Exec(
		sql.Named("offset", offset),
		sql.Named("first", first),
		sql.Named("count", last+1-first),
	)
	if err != nil {
		return err
	}
	_, err = tx.Stmt(d.rotateWords2Stmt).Exec()
	return err
}

type ErrNotFound error
