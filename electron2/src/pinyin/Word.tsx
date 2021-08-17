import React from 'react';
import Pinyin from './Pinyin';

interface WordProps {
  word: Word | string | Pinyin[];
  [attrs: string]: unknown;
}

export default class Word {
  readonly chars: Pinyin[];

  constructor(word: string | Pinyin[]) {
    if (typeof word === "string") {
      const arr: Pinyin[] = [];
      while (word !== "") {
        const p = new Pinyin(word);
        word = word.slice(p.consumed);
        arr.push(p);
      }
      this.chars = arr;
    } else {
      this.chars = word;
    }
  }

  get length(): number { return this.chars.length; }
  get pinyin(): string { return this.chars.map(c => c.pinyin).join(' '); }
  get raw   (): string { return this.chars.map(c => c.raw   ).join(' '); }

  html(props: {[attrs: string]: unknown}): JSX.Element {
    return (
      <React.Fragment {...props}>
        {this.chars.map((c, i) => c.html({key: i}))}
      </React.Fragment>
    );
  }

  static HTML({word, ...props}: WordProps): JSX.Element {
    if (!(word instanceof Word)) {
      word = new Word(word);
    }
    return word.html(props);
  }

  static compare(a: Word, b: Word): number {
    const n = Math.max(a.length, b.length);
    for (let i = 0; i < n; ++i) {
      const c = Pinyin.compare(a.chars[i], b.chars[i]);
      if (c !== 0) {
        return c;
      }
    }
    return a.length - b.length;
  }

  slice(start?: number | undefined, end?: number | undefined): Word {
    return new Word(this.chars.slice(start, end));
  }

  *[Symbol.iterator](): Iterator<Pinyin> {
    for (const char of this.chars) {
      yield char;
    }
  }
}
