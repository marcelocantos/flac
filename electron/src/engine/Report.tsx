import React from 'react';

import Proxy from '../renderer/data/Proxy';
import { Refdata } from '../refdata/Refdata';
import Outcome from '../outcome/Outcome';

function logScore(score: number): number {
	return Math.log(score) / Math.log(4000)
}

export interface 汇报项目 {
  get html(): JSX.Element;
}

class 错误汇报 implements 汇报项目 {
  constructor(
    private 字: string,
    private 回答: string,
  ){}

  get html(): JSX.Element {
    return <>{this.字} ≠ {this.回答}</>;
  }
}

class 好汇报 implements 汇报项目 {
  constructor(
    private 好组: {字: string, 分数: number}[],
  ){}

  get html(): JSX.Element {
    return <></>;
  }
}

export default class 汇报类 {
	refreshCount: number;

	历史:   string[];
	好清单: string[];
	// msgs:    string[];

	onScoreChangedFunc: (word: string, score: number) => void;

  constructor(
    private db: Proxy,
    private rd: Refdata,
  ){}

  async 好(字: string, 产物: Outcome, 容易: boolean): Promise<void> {
    this.bump(字, score => {
      // 产物.html.分数(score);
      return {score: Math.max(2, 2 * score), move: true};
    });

    const score = await this.score(字);

    // this.appendGoods(字, score);
    this.ClearMessages();
  }

  async 不好(o: Outcome, easy: boolean, 包装的尝试: {尝试: number}): Promise<void> {
    // defer this.refresh()()

    // if o.Fail() {
    //   if err := this.bad(o, easy, attempt); err != null {
    //     return err
    //   }
    // }

    // this.ClearMessages()

    // if len(o.Bad) > 0 {
    //   prefix := strings.Repeat(" ", 3+2*len([]rune(o.Word))+2)
    //   top := prefix
    //   var corrections []string[]
    //   for _, word := range o.Bad {
    //     wordLen := len([]rune(word.String()))
    //     middle := (wordLen - 1) / 2
    //     tail := wordLen - middle - 1
    //     var correction string
    //     if dancis, has := this.rd.Dict.PinyinToSimplified[word.RawString()]; has {
    //       correction = strings.Join(dancis.Words, " ")
    //     } else {
    //       correction = "∅"
    //     }
    //     top = fmt.Sprintf("%s %s┬%s", top, strings.Repeat("─", middle), strings.Repeat("─", tail))
    //     corrections = append(corrections, string[]{
    //       fmt.Sprintf("%s %*s╘👉 %s", prefix, middle, "", correction),
    //     })
    //     prefix = fmt.Sprintf("%s %s│%s", prefix, strings.Repeat(" ", middle), strings.Repeat(" ", tail))
    //   }
    //   for i := len(corrections) - 1; i >= 0; i-- {
    //     for _, line := range corrections[i] {
    //       this.appendMessage("%s", line)
    //     }
    //   }

    //   this.appendHistory(fmt.Sprintf(
    //     "❌ %s ≠ %s\034❌ [#999999::]%[1]s ≠ [#999999::d]%[3]s[-::-]",
    //     o.Word, o.Bad.ColorString("u"), o.Bad.String()))
    // }
    // if len(o.TooShort) > 0 {
    //   this.appendMessage("⚠️  Missing characters: %s...", o.TooShort.ColorString(""))
    // }
    // if len(o.Bad) == 0 && o.Missing > len(o.TooShort)+len(o.BadTones) {
    //   this.appendMessage("⚠️  Missing alternative%s[-::]", pluralS(o.Missing))
    // }
    // if len(o.BadTones) > 0 {
    //   this.appendMessage("[:silver:]🎵[:-:] Only tone(s) need correcting!")
    // }
  }

  GiveUp(outcome: Outcome) {
    this.setMessages("TODO: outcome.更正");
    return this.bump(outcome.Word, score => {
      return {score: Math.max(1, score / 8), move: false};
    })
  }

  appendGoods(...好清单: string[]) {
    if (好清单.length > 0) {
      this.好清单.push(...好清单);
    }
  }

  clearGoods() {
    if (this.好清单.length > 0) {
      this.appendHistory(...this.goodsReport());
      this.好清单 = [];
    }
  }

  goodsReport(): string[] {
    if (this.好清单.length === 0) {
      return [];
    }
    return [this.好清单.join(" ")];
  }

  appendHistory(...lines: string[]) {
    if (lines.length > 0) {
      this.历史.push(...lines);
    }
  }

  async bump(word: string, bump: (score: number) => {score: number, move: boolean}): Promise<void> {
    const {score, move} = bump(await this.score(word));
    let pos = -1;
    if (move) {
      pos = score + Math.floor(Math.random() * (1+score*3/2-score));
    }

    await this.setScoreAndPos(word, score, pos);
  }

  async score(word: string): Promise<number> {
    return await this.db.WordScore(word);
  }

  async setScoreAndPos(word: string, score: number, pos: number): Promise<void> {
    await this.db.UpdateScoreAndPos(word, score, pos);
  }

  async bad(outcome: Outcome, easy: boolean, attempt: {attempt: number}): Promise<void> {
    try {
      const penalty = Math.sqrt(1 + attempt.attempt);
      attempt.attempt++;

      // Multiply score by 1/2√(1 + attempt).
      await this.bump(outcome.Word, score => ({score: Math.max(1, score / (2 * penalty)), move: false}));

      this.ClearMessages()
    } finally {
      this.clearGoods();
    }
  }

  ClearMessages() {
    // this.msgs = [];
  }

  setMessages(...messages: string[]) {
    // this.msgs = messages;
  }
}
