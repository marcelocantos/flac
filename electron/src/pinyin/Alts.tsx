import React from 'react';
import Word from './Word';
import Pinyin from './Pinyin';

const pinyinsRE = /^([a-zü]+)([1-5]+)$/i;

export default class Alts {
  public words: Word[];

  constructor(words: string | Word[]) {
    if (typeof words == "string") {
      const g = words.match(pinyinsRE);
      if (!g) {
        throw new Error(`${words}: not valid pinyin`);
      }

      let tones = 0;
      for (const t of g[2]) {
        tones |= 1 << Number.parseInt(t);
      }

      const syllable = g[1];
      const arr: Word[] = [];
      for (let tone = 1; tone <= 5; tone++) {
        if (tones & (1 << tone)) {
          const p = new Word([new Pinyin({syllable, tone})]);
          arr.push(p);
        }
      }
      this.words = arr;
    } else {
      this.words = words;
    }
  }

  get length(): number { return this.words.length; }
  get pinyin(): string { return this.words.map(w => w.pinyin).join('/'); }
  get raw   (): string { return this.words.map(w => w.raw   ).join('/'); }

  get html() {
    return <span>{this.words.map((c, i) => <>{i ? '/' : ''}{c.html}</>)}</span>;
  }

  static compare(a: Alts, b: Alts): number {
    const n = Math.max(a.length, b.length);
    for (let i = 0; i < n; ++i) {
      const c = Word.compare(a.words[i], b.words[i]);
      if (c !== 0) {
        return c;
      }
    }
    return a.length - b.length;
  }

  *[Symbol.iterator]() {
    for (const word of this.words) {
      yield word;
    }
  }
}
