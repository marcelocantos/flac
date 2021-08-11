import Alts from '../pinyin/Alts';
import Word from '../pinyin/Word';
import { Entries } from '../refdata/Refdata';

export default class Outcome {
  Good:     Word[] = [];
	Bad:      Word[] = [];
	TooShort: Word[] = [];
	BadTones: Word[] = [];
	Missing:  number = 0;
	Easy:     boolean = false;

  constructor(
    public Word:    string,
    public Entries: Entries,
  ){}

  get Pass(): boolean {
    return this.Bad.length + this.TooShort.length + this.BadTones.length + this.Missing === 0;
  }

  get Fail(): boolean {
    return this.Bad.length + this.BadTones.length > 0;
  }

  get Correction() {
    return <>{this.Word} = {this.WordAlts.html}</>;
  }

  get WordAlts(): Alts {
    return new Alts(Object.keys(this.Entries.entries).map(raw => new Word(raw)).sort());
  }
}
