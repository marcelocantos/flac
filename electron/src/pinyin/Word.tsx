import Pinyin from './Pinyin';

export default class Word {
  readonly chars: Pinyin[];

  constructor(chars: string | Pinyin[]) {
    if (typeof chars === "string") {
      const arr: Pinyin[] = [];
      while (chars != "") {
        const p = new Pinyin(chars);
        chars = chars.slice(p.raw.length);
        arr.push(p);
      }
      this.chars = arr;
    } else {
      this.chars = chars;
    }
  }

  get length(): number { return this.chars.length; }
  get pinyin(): string { return this.chars.map(c => c.pinyin).join(' '); }
  get raw   (): string { return this.chars.map(c => c.raw   ).join(' '); }

  get html() {
    return <span>{this.chars.map(c => c.html)}</span>;
  }

  static compare(a: Word, b: Word): number {
    const n = Math.max(a.length, b.length);
    for (let i = 0; i < n; ++i) {
      const c = Pinyin.compare(a.chars[i], b.chars[i]);
      if (c != 0) {
        return c;
      }
    }
    return a.length - b.length;
  }
}
