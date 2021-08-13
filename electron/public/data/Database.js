class Database {
  focusID;
  db;

  selMaxScore;
  selMaxPos;
  selWordScore;
  selWordPos;
  selWordAt;
  selFocusID;

  selQueuedWord;

  insertFocus;
  enqueueWord;
  dequeueWord;
  updateScore;
  rotateWords1;
  rotateWords2;

  constructor(fields){
    this.focusID = {$focusID: 0};
    for (const f in fields) {
      this[f] = fields[f];
    }
  }

  static async build(db, focus, words) {
    if (!focus) {
      focus = "";
    }
    if (!words) {
      words = [];
    }
    await db.exec(`
      PRAGMA foreign_keys = ON;

      CREATE TABLE IF NOT EXISTS word_score (
        word  TEXT    PRIMARY KEY,
        score INTEGER
      );

      CREATE TABLE IF NOT EXISTS queue (
        pos  INT  PRIMARY KEY,
        word TEXT UNIQUE
      );

      CREATE TABLE IF NOT EXISTS focus_queue (
        focusID INTEGER REFERENCES focus (focusID),
        pos     INT,
        word    TEXT,

        PRIMARY KEY (focusID, pos),
        UNIQUE      (focusID, word)
      );

      CREATE TABLE IF NOT EXISTS focus (
        focusID INTEGER PRIMARY KEY,
        focus   TEXT    UNIQUE
      );

      INSERT OR IGNORE INTO focus (focus) VALUES ('');

      INSERT OR IGNORE INTO focus_queue (focusID, word, pos)
            SELECT focusID, word, pos
            FROM   queue CROSS JOIN focus
            WHERE  focus.focus = '';

      DELETE FROM queue;
    `);

    const d = new Database({
      db,

      selMaxScore: (await db.prepare(
        `SELECT MAX(score) AS maxScore FROM word_score`)).get,
      selMaxPos: (await db.prepare(
        `SELECT MAX(pos) AS maxPos FROM focus_queue WHERE focusID = $focusID`)).get,
      selWordScore: (await db.prepare(
        `SELECT score FROM word_score WHERE word = $word`)).get,
      selWordPos: (await db.prepare(
        `SELECT pos FROM focus_queue WHERE word = $word AND focusID = $focusID`)).get,
      selWordAt: (await db.prepare(
        `SELECT word FROM focus_queue WHERE pos = $pos AND focusID = $focusID`)).get,
      selFocusID: (await db.prepare(
        `SELECT focusID FROM focus WHERE focus = $focus`)).get,

      selQueuedWords: (await db.prepare(
        `SELECT word FROM focus_queue WHERE focusID = $focusID`)).all,

      insertFocus: (await db.prepare(
        `INSERT OR IGNORE INTO focus (focus) VALUES ($focus)`)).run,
      enqueueWord: (await db.prepare(
        `INSERT INTO focus_queue (focusID, pos, word) VALUES ($focusID, $pos, $word)`)).run,
      dequeueWord: (await db.prepare(
        `DELETE FROM focus_queue WHERE word = $word AND focusID = $focusID`)).run,
      updateScore: (await db.prepare(
        `INSERT OR REPLACE INTO word_score (word, score) VALUES ($word, $score)`)).run,
      rotateWords1: (await db.prepare(`
        UPDATE focus_queue
        SET pos = -1 - ((pos - $first + $count + $offset) % $count + $first)
        WHERE pos BETWEEN $first AND $first + $count - 1
              AND focusID = $focusID
      `)).run,
      rotateWords2: (await db.prepare(`
        UPDATE focus_queue
        SET pos = -1-pos
        WHERE pos < 0 AND focusID = $focusID
      `)).run,
    });

    await db.tx(async () => {
      await d.SetFocus(focus);

      let maxPos = await d.maxPos() || -1;
      const positions = {};

      for (let i = 0; i < words.length; ++i) {
        const $word = words[i];
        let wordPos = await d.wordPos($word);
        if (typeof wordPos == "undefined") {
          wordPos = ++maxPos;
          await d.enqueueWord({...d.focusID, $pos: wordPos, $word});
        }
        positions[$word] = wordPos;
        // console.log({queue: await d.selQueuedWords({...d.focusID})});
      }

      let remove = [];
      for (const {word} of await d.selQueuedWords({...d.focusID})) {
        if (!(word in positions)) {
          remove.push(word);
        }
      }
      await d.removeWords(remove);
    });

    return d;
  }

  async close() {
    await this.db.close();
  }

  get HeadWord() {
    return this.WordAt(0);
  }

  get MaxScore() {
    return this.maxScore();
  }

  get MaxPos() {
    return this.maxPos();
  }

  MoveWord(word, dest) {
    return this.db.tx(() => this.moveWord(word, dest));
  }

  async UpdateScoreAndPos($word, $score, $dest) {
    await this.updateScore({$word, $score});
    if (typeof $dest !== "undefined") {
      await this.moveWord($word, $dest);
    }
  }

  async SetFocus($focus) {
    await this.insertFocus({$focus});
    this.focusID = {$focusID: (await this.selFocusID({$focus}))?.focusID};
  }

  WordScore($word) {
    return this.db.tx(() => this.wordScore($word));
  }

  WordPos($word) {
    return this.db.tx(() => this.wordPos($word));
  }

  async WordAt($pos) {
    return (await this.selWordAt({...this.focusID, $pos}))?.word;
  }

  async maxScore() {
    return (await this.selMaxScore())?.maxScore;
  }

  async maxPos() {
    return (await this.selMaxPos(this.focusID))?.maxPos;
  }

  async moveWord($word, dest) {
    const src = await this.wordPos($word)
    if (src === undefined) {
      throw new RangeError(`${$word} not in queue`);
    }
    const maxPos = await this.maxPos();
    dest = Math.min(Math.max(0, dest), maxPos);
    let $first = Math.min(src, dest), last = Math.max(src, dest);
    let $offset = src < dest ? -1 : 1;

    await this.rotateWords1({...this.focusID, $offset, $first, $count: last+1-$first});
    await this.rotateWords2(this.focusID);
  }

  async removeWords(words) {
    const max = await this.maxPos() || 0;
    for (const $word of words) {
      await this.moveWord($word, max)
      await this.dequeueWord({...this.focusID, $word})
    }
  }

  async selectInt(get, params, col) {
    return ((await get(params)) ?? {})[col];
  }

  wordScore($word) {
    return this.selectInt(this.selWordScore, {$word}, 'score');
  }

  wordPos($word) {
    return this.selectInt(this.selWordPos, {...this.focusID, $word}, 'pos');
  }
}

module.exports = Database;
