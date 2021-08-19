import Alts from '../pinyin/Alts';
import Word from '../pinyin/Word';
import { Entries } from '../refdata/Refdata';

export default class Outcome {
  html?: ({分数}: {分数: number}) => JSX.Element;
  Good:     Word[] = [];
	Bad:      Word[] = [];
	TooShort: Word[] = [];
	BadTones: Word[] = [];
	Missing = 0;
	Easy = false;

  constructor(
    public Word:    string,
    public Entries: Entries,
  ){}

  get 及格(): boolean {
    return this.Bad.length + this.TooShort.length + this.BadTones.length + this.Missing === 0;
  }

  get 不及格(): boolean {
    return this.Bad.length + this.BadTones.length > 0;
  }

  get WordAlts(): Alts {
    return new Alts(Object.keys(this.Entries.entries).map(raw => new Word(raw)).sort());
  }
}
