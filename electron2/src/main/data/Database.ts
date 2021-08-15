import * as AsyncDB from './AsyncDB';
import * as Interface from '../../common/data/Interface';

interface Stmts {
  selMaxScore:    AsyncDB.Get;
  selMaxPos:      AsyncDB.Get;
  selWordScore:   AsyncDB.Get;
  selWordPos:     AsyncDB.Get;
  selWordAt:      AsyncDB.Get;
  selFocusID:     AsyncDB.Get;

  selQueuedWords: AsyncDB.All;

  insertFocus:    AsyncDB.Run;
  enqueueWord:    AsyncDB.Run;
  dequeueWord:    AsyncDB.Run;
  updateScore:    AsyncDB.Run;
  rotateWords1:   AsyncDB.Run;
  rotateWords2:   AsyncDB.Run;
}

export default class Database implements Interface.Database {
  constructor(
    public focusID: {$focusID: number},
    public db: AsyncDB.Database,
    public s: Stmts,
  ){}

  static async build(
    db: AsyncDB.Database,
    focus: string,
    words: string[],
  ): Promise<Database> {
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

    const d = new Database(
      {$focusID: 0},
      db,
      {
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
      },
    );

    await db.tx(async () => {
      await d.setFocus(focus);

      let maxPos = await d.maxPos() || -1;
      const positions: {[word: string]: number} = {};

      for (let i = 0; i < words.length; ++i) {
        const $word = words[i];
        let wordPos = await d.wordPos($word);
        if (typeof wordPos == "undefined") {
          wordPos = ++maxPos;
          await d.s.enqueueWord({...d.focusID, $pos: wordPos, $word});
        }
        positions[$word] = wordPos;
        // console.log({queue: await d.selQueuedWords({...d.focusID})});
      }

      const remove: string[] = [];
      for (const row of await d.s.selQueuedWords({...d.focusID})) {
        const word = row.word as string;
        if (!(word in positions)) {
          remove.push(word);
        }
      }
      await d.removeWords(remove);
    });

    return d;
  }

  async close(): Promise<void> {
    await this.db.close();
  }

  get HeadWord(): Promise<string> {
    return this.WordAt(0);
  }

  get MaxScore(): Promise<number> {
    return this.maxScore();
  }

  get MaxPos(): Promise<number> {
    return this.maxPos();
  }

  MoveWord(word: string, dest: number): Promise<void> {
    return this.db.tx(() => this.moveWord(word, dest));
  }

  async UpdateScore(word: string, score: number): Promise<void> {
    return this.db.tx(async () => {
      await this.s.updateScore({$word: word, $score: score});
    });
  }

  async UpdateScoreAndPos(word: string, score: number, dest: number): Promise<void> {
    return this.db.tx(async () => {
      await this.s.updateScore({$word: word, $score: score});
      await this.moveWord(word, dest);
    });
  }

  async SetFocus(focus: string): Promise<void> {
    return this.db.tx(async () => {
      await this.setFocus(focus);
    });
  }

  WordScore(word: string): Promise<number> {
    return this.db.tx(() => this.wordScore(word));
  }

  WordPos(word: string): Promise<number> {
    return this.db.tx(() => this.wordPos(word));
  }

  async WordAt(pos: number): Promise<string> {
    return (await this.s.selWordAt({...this.focusID, $pos: pos}))?.word as string;
  }

  private async maxScore(): Promise<number> {
    return (await this.s.selMaxScore())?.maxScore as number;
  }

  private async maxPos(): Promise<number> {
    return (await this.s.selMaxPos(this.focusID))?.maxPos as number;
  }

  private async moveWord($word: string, dest?: number): Promise<void> {
    const src = await this.wordPos($word)
    if (src === undefined) {
      throw new RangeError(`${$word} not in queue`);
    }
    const maxPos = await this.maxPos();
    dest = Math.min(Math.max(0, dest), maxPos);
    const $first = Math.min(src, dest), last = Math.max(src, dest);
    const $offset = src < dest ? -1 : 1;

    await this.s.rotateWords1({...this.focusID, $offset, $first, $count: last+1-$first});
    await this.s.rotateWords2(this.focusID);
  }

  private async removeWords(words: string[]): Promise<void> {
    const max = await this.maxPos() || 0;
    for (const $word of words) {
      await this.moveWord($word, max)
      await this.s.dequeueWord({...this.focusID, $word})
    }
  }

  private async selectInt(get: AsyncDB.Get, params: AsyncDB.Params, col: string): Promise<number> {
    return ((await get(params)) ?? {})[col] as number;
  }

  private async setFocus($focus: string): Promise<void> {
    await this.s.insertFocus({$focus});
    const $focusID = (await this.s.selFocusID({$focus}))?.focusID as number;
    this.focusID = {$focusID};
}

  private wordScore($word: string): Promise<number> {
    return this.selectInt(this.s.selWordScore, {$word}, 'score');
  }

  private wordPos($word: string): Promise<number> {
    return this.selectInt(this.s.selWordPos, {...this.focusID, $word}, 'pos');
  }
}
