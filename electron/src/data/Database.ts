import * as AsyncDB from './AsyncDB';

export default class Database {
  focusID: {$focusID: number};

  constructor(
    public db: any,

    public selMaxScore:    AsyncDB.Getter,
    public selMaxPos:      AsyncDB.Getter,
    public selWordScore:   AsyncDB.Getter,
    public selWordPos:     AsyncDB.Getter,
    public selWordAt:      AsyncDB.Getter,
    public selFocusID:     AsyncDB.Getter,

    public selQueuedWords: AsyncDB.Aller,

    public insertFocus:    AsyncDB.Runner,
    public enqueueWord:    AsyncDB.Runner,
    public dequeueWord:    AsyncDB.Runner,
    public updateScore:    AsyncDB.Runner,
    public rotateWords1:   AsyncDB.Runner,
    public rotateWords2:   AsyncDB.Runner,
  ){
    this.focusID = {$focusID: 0};
  }

  static async build(db: AsyncDB.Database, focus: string, words: string[]): Promise<Database> {
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
      db,

      (await db.prepare(
        `/*selMaxScore*/ SELECT MAX(score) AS maxScore FROM word_score`)).get,
      (await db.prepare(
        `/*selMaxPos*/ SELECT MAX(pos) AS maxPos FROM focus_queue WHERE focusID = $focusID`)).get,
      (await db.prepare(
        `/*selWordScore*/ SELECT score FROM word_score WHERE word = $word`)).get,
      (await db.prepare(
        `/*selWordPos*/ SELECT pos FROM focus_queue WHERE word = $word AND focusID = $focusID`)).get,
      (await db.prepare(
        `/*selWordAt*/ SELECT word FROM focus_queue WHERE pos = $pos AND focusID = $focusID`)).get,
      (await db.prepare(
        `/*selFocusID*/ SELECT focusID FROM focus WHERE focus = $focus`)).get,

      (await db.prepare(
        `/*selQueuedWords*/ SELECT word FROM focus_queue WHERE focusID = $focusID`)).all,

      (await db.prepare(
        `/*insertFocus*/ INSERT OR IGNORE INTO focus (focus) VALUES ($focus)`)).run,
      (await db.prepare(
        `/*enqueueWord*/ INSERT INTO focus_queue (focusID, pos, word) VALUES ($focusID, $pos, $word)`)).run,
      (await db.prepare(
        `/*dequeueWord*/ DELETE FROM focus_queue WHERE word = $word AND focusID = $focusID`)).run,
      (await db.prepare(
        `/*updateScore*/ INSERT OR REPLACE INTO word_score (word, score) VALUES ($word, $score)`)).run,
      (await db.prepare(`
        /*rotateWords1*/ UPDATE focus_queue
        SET pos = -1-((pos - $first + $count + $offset) % $count + $first)
        WHERE pos BETWEEN $first AND $first + $count - 1
              AND focusID = $focusID
      `)).run,
      (await db.prepare(`
        /*rotateWords2*/ UPDATE focus_queue
        SET pos = -1-pos
        WHERE pos < 0 AND focusID = $focusID
      `)).run,
    );

    await db.tx(async () => {
      await d.SetFocus(focus);

      let maxPos = await d.maxPos() || -1;
      const positions: {[id: string]: any} = {};

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

      let remove: string[] = [];
      for (const {word} of await d.selQueuedWords({...d.focusID})) {
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

  HeadWord(): Promise<string | undefined> {
    return this.WordAt(0);
  }

  MaxScore(): Promise<number | undefined> {
    return this.maxScore();
  }

  MaxPos(): Promise<number | undefined> {
    return this.maxPos();
  }

  MoveWord(word: string, dest: number): Promise<void> {
    return this.db.tx(() => this.moveWord(word, dest));
  }

  async UpdateScoreAndPos($word: string, $score: number, $dest?: number): Promise<void> {
    await this.updateScore({$word, $score});
    if (typeof $dest !== "undefined") {
      await this.moveWord($word, $dest);
    }
  }

  async SetFocus($focus: string): Promise<void> {
    await this.insertFocus({$focus});
    this.focusID = {$focusID: (await this.selFocusID({$focus}))?.focusID};
  }

  WordScore($word: string): Promise<number | undefined> {
    return this.db.tx(() => this.wordScore($word));
  }

  WordPos($word: string): Promise<number | undefined> {
    return this.db.tx(() => this.wordPos($word));
  }

  async WordAt($pos: number): Promise<string | undefined> {
    return (await this.selWordAt({...this.focusID, $pos}))?.word;
  }

  private async maxScore(): Promise<number | undefined> {
    return (await this.selMaxScore())?.maxScore;
  }

  private async maxPos(): Promise<number | undefined> {
    return (await this.selMaxPos(this.focusID))?.maxPos;
  }

  private async moveWord($word: string, dest: number): Promise<void> {
    const src = await this.wordPos($word)
    if (src === undefined) {
      throw new RangeError(`${$word} not in queue`);
    }
    const maxPos = await this.maxPos() as number;
    dest = Math.min(Math.max(0, dest), maxPos);
    let $first = Math.min(src, dest), last = Math.max(src, dest);
    let $offset = src < dest ? -1 : 1;

    await this.rotateWords1({...this.focusID, $offset, $first, $count: last+1-$first});
    await this.rotateWords2(this.focusID);
  }

  private async removeWords(words: string[]): Promise<void> {
    const max = await this.maxPos() || 0;
    for (const $word of words) {
      await this.moveWord($word, max)
      await this.dequeueWord({...this.focusID, $word})
    }
  }

  private async selectInt(get: AsyncDB.Getter, params: {[key: string]: any}, col: string): Promise<number> {
    return ((await get(params)) ?? {})[col];
  }

  private wordScore($word: string): Promise<number | undefined> {
    return this.selectInt(this.selWordScore, {$word}, 'score');
  }

  private wordPos($word: string): Promise<number | undefined> {
    return this.selectInt(this.selWordPos, {...this.focusID, $word}, 'pos');
  }
}
