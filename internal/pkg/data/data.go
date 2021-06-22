package data

import (
	"database/sql"
	"fmt"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB

	// Read
	focusIDStmt     *sql.Stmt
	maxScoreStmt    *sql.Stmt
	maxPosStmt      *sql.Stmt
	wordScoreStmt   *sql.Stmt
	wordPosStmt     *sql.Stmt
	wordAtStmt      *sql.Stmt
	queuedWordsStmt *sql.Stmt

	// Write
	insertFocusStmt  *sql.Stmt
	enqueueWordStmt  *sql.Stmt
	dequeueWordStmt  *sql.Stmt
	updateScoreStmt  *sql.Stmt
	rotateWords1Stmt *sql.Stmt
	rotateWords2Stmt *sql.Stmt

	focusID sql.NamedArg
}

func NewDatabase(path string) (*Database, error) {
	var d Database
	var err error
	d.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	createTable := func(def string) string {
		return "CREATE TABLE IF NOT EXISTS " + def
	}
	for _, def := range []string{
		createTable(`word_score (
            word  TEXT    PRIMARY KEY,
            score INTEGER
        )`),
		createTable(`queue (
            pos  INT  PRIMARY KEY,
            word TEXT UNIQUE
        )`),
		createTable(`focus_queue (
            focusID INTEGER REFERENCES focus (focusID),
            pos     INT,
            word    TEXT,

            PRIMARY KEY (focusID, pos),
            UNIQUE      (focusID, word)
        )`),
		createTable(`focus (
            focusID INTEGER PRIMARY KEY,
            focus   TEXT    UNIQUE
        )`),
		`PRAGMA foreign_keys = ON`,
		`INSERT OR IGNORE INTO focus (focus) VALUES ('')`,
		`INSERT OR IGNORE INTO focus_queue (focusID, word, pos)
            SELECT focusID, word, pos
            FROM   queue CROSS JOIN focus
			WHERE  focus.focus = ''
        `,
		`DELETE FROM queue`,
	} {
		_, err := d.db.Exec(def)
		if err != nil {
			tx.Rollback()
			return nil, errors.WrapPrefix(err, def, 0)
		}
	}
	tx.Commit()

	for stmt, query := range map[**sql.Stmt]string{
		&d.focusIDStmt:     `SELECT focusID FROM focus WHERE focus = ?`,
		&d.maxScoreStmt:    `SELECT COALESCE(MAX(score), -1) FROM word_score`,
		&d.maxPosStmt:      `SELECT COALESCE(MAX(pos), -1) FROM focus_queue WHERE focusID = $focusID`,
		&d.wordScoreStmt:   `SELECT score FROM word_score WHERE word = ?`,
		&d.wordPosStmt:     `SELECT pos FROM focus_queue WHERE word = $word AND focusID = $focusID`,
		&d.wordAtStmt:      `SELECT word FROM focus_queue WHERE pos = $pos AND focusID = $focusID`,
		&d.queuedWordsStmt: `SELECT word FROM focus_queue WHERE focusID = $focusID`,

		&d.insertFocusStmt: `INSERT OR IGNORE INTO focus (focus) VALUES (?)`,
		&d.enqueueWordStmt: `INSERT INTO focus_queue (focusID, pos, word) VALUES ($focusID, $pos, $word)`,
		&d.dequeueWordStmt: `DELETE FROM focus_queue WHERE word = $word AND focusID = $focusID`,
		&d.updateScoreStmt: `INSERT OR REPLACE INTO word_score (word, score) VALUES (?, ?)`,
		&d.rotateWords1Stmt: `
			UPDATE focus_queue
			SET pos = -1-((pos - $first + $count + $offset) % $count + $first)
			WHERE pos BETWEEN $first AND $first + $count - 1
			      AND focusID = $focusID
		`,
		&d.rotateWords2Stmt: `
			UPDATE focus_queue
			SET pos = -1-pos
			WHERE pos < 0 AND focusID = $focusID
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

func (d *Database) Populate(focus string, words []string) (e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer commit(tx, &e)

	if err := d.SetFocus(tx, focus); err != nil {
		return err
	}

	pos, err := d.maxPos(tx)
	if err != nil {
		return err
	}

	positions := make(map[string]int, len(words))

	getWordPos := d.wordPos(tx)
	enqueueWord := tx.Stmt(d.enqueueWordStmt)
	for i, word := range words {
		positions[word] = i
		_, err := getWordPos(word)
		if err != nil {
			if _, is := err.(ErrNotFound); !is {
				return err
			}
			pos++
			if _, err := enqueueWord.Exec(
				d.focusID,
				sql.Named("pos", pos),
				sql.Named("word", word),
			); err != nil {
				return err
			}
		}
	}

	queuedWords := tx.Stmt(d.queuedWordsStmt)
	rows, err := queuedWords.Query(d.focusID)
	if err != nil {
		return err
	}
	var word string
	var remove []string
	for rows.Next() {
		if err := rows.Scan(&word); err != nil {
			return errors.Wrap(err, 0)
		}
		if _, has := positions[word]; !has {
			remove = append(remove, word)
		}
	}
	return d.removeWords(tx, remove)
}

func (d *Database) SetFocus(tx *sql.Tx, focus string) (e error) {
	_, err := tx.Stmt(d.insertFocusStmt).Exec(focus)
	if err != nil {
		return err
	}

	var focusID uint64
	if err := tx.Stmt(d.focusIDStmt).QueryRow(focus).Scan(&focusID); err != nil {
		return err
	}
	d.focusID = sql.Named("focusID", focusID)

	return nil
}

func (d *Database) MaxScore() (_ int, e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer commit(tx, &e)

	return d.maxScore(tx)
}

func (d *Database) maxScore(tx *sql.Tx) (int, error) {
	var score sql.NullInt64
	if err := tx.Stmt(d.maxScoreStmt).QueryRow().Scan(&score); err != nil {
		return 0, err
	}
	return int(score.Int64), nil
}

func (d *Database) MaxPos() (_ int, e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer commit(tx, &e)

	return d.maxPos(tx)
}

func (d *Database) maxPos(tx *sql.Tx) (int, error) {
	var pos sql.NullInt64
	if err := tx.Stmt(d.maxPosStmt).QueryRow(d.focusID).Scan(&pos); err != nil {
		return 0, err
	}
	return int(pos.Int64), nil
}

func (d *Database) WordScore(word string) (_ int, e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer commit(tx, &e)

	return d.wordScore(tx)(word)
}

func (d *Database) wordScore(tx *sql.Tx) func(word string) (int, error) {
	getWordScoreStmt := tx.Stmt(d.wordScoreStmt)
	return func(word string) (int, error) {
		return d.selectInt(getWordScoreStmt, "%s: not found in word_score", word)
	}
}

func (d *Database) WordPos(word string) (_ int, e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer commit(tx, &e)

	return d.wordPos(tx)(word)
}

func (d *Database) wordPos(tx *sql.Tx) func(word string) (int, error) {
	getWordPosStmt := tx.Stmt(d.wordPosStmt)
	return func(word string) (int, error) {
		return d.selectInt(getWordPosStmt, "%s: not found in queue", sql.Named("word", word), d.focusID)
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

func (d *Database) WordAt(pos int) (_ string, e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return "", err
	}
	defer commit(tx, &e)

	return d.wordAt(tx)(pos)
}

func (d *Database) wordAt(tx *sql.Tx) func(pos int) (string, error) {
	getWordAtStmt := tx.Stmt(d.wordAtStmt)
	return func(pos int) (string, error) {
		var word string
		err := getWordAtStmt.QueryRow(sql.Named("pos", pos), d.focusID).Scan(&word)
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

func (d *Database) UpdateScoreAndPos(word string, score, dest int) (e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer commit(tx, &e)

	if _, err := tx.Stmt(d.updateScoreStmt).Exec(word, score); err != nil {
		return err
	}
	if dest >= 0 {
		return d.moveWord(tx, word, dest)
	}
	return nil
}

func (d *Database) MoveWord(word string, dest int) (e error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer commit(tx, &e)

	return d.moveWord(tx, word, dest)
}

func (d *Database) moveWord(tx *sql.Tx, word string, dest int) error {
	if dest < 0 {
		dest = 0
	} else {
		max, err := d.maxPos(tx)
		if err != nil {
			return err
		}
		if dest > max {
			dest = max
		}
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
		d.focusID,
	)
	if err != nil {
		return err
	}
	_, err = tx.Stmt(d.rotateWords2Stmt).Exec(d.focusID)
	return err
}

func (d *Database) removeWords(tx *sql.Tx, words []string) error {
	max, err := d.maxPos(tx)
	if err != nil {
		return err
	}

	for _, word := range words {
		if err := d.moveWord(tx, word, max); err != nil {
			return err
		}

		if _, err := tx.Stmt(d.dequeueWordStmt).Exec(sql.Named("word", word), d.focusID); err != nil {
			return err
		}
	}
	return err
}

type ErrNotFound error

func commit(tx *sql.Tx, err *error) {
	if *err != nil {
		tx.Rollback()
	} else {
		*err = tx.Commit()
	}
}
