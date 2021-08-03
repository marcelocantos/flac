import React from 'react';

import PinyinRE from './PinyinRE';

const toneColors = [
  '',
  'red',
  'green',
  'blue',
  'purple',
  'black',
];

const baseVowels = 'AEIOUÜaeiouü';

const mark = [
  '',
  'ĀĒĪŌŪǕāēīōūǖ',
  'ÁÉÍÓÚǗáéíóúǘ',
  'ǍĚǏǑǓǙǎěǐǒǔǚ', // Breve forms: ăĕĭŏŭ-ĂĔĬŎŬ-
  'ÀÈÌÒÙǛàèìòùǜ',
  baseVowels,
].map(vowels => (v: string) => vowels[baseVowels.indexOf(v)] ?? '')

export default class Pinyin {
  readonly pinyin: string;
  readonly syllable: string;
  readonly tone: number;
  readonly raw: string;

  constructor(arg: string | {syllable: string, tone: number}) {
    if (typeof arg === 'string') {
      let groups = PinyinRE.exec(arg);
      if (groups == null || (!groups[1] && groups[3].length > 1)) {
        throw new Error(`${arg}: invalid pinyin`);
      }
      arg = {
        syllable: groups[2].replace('v', 'ü').replace('u:', 'ü'),
        tone: Number.parseInt(groups[3]),
      }
    }

    let {syllable, tone} = arg;

    // https://en.wikipedia.org/wiki/Pinyin#Rules_for_placing_the_tone_mark
    this.pinyin   = syllable.replace(/[aeo]|(?<=i)u|(?<=u)i|[iuü]/i, mark[tone]);
    this.syllable = syllable.replace(/ü/i, 'v');
    this.tone   = tone;
    this.raw      = `${this.pinyin}${this.tone}`;
  }

  get color (): string { return toneColors[this.tone]; }

  get html() {
    return <span style={{color: this.color}}>{this.pinyin}</span>;
  }

  static compare(a: Pinyin, b: Pinyin): number {
    const aLower = a.syllable.toLowerCase();
    const bLower = b.syllable.toLowerCase();
    return (
      aLower < bLower ? -1 :
      aLower > bLower ? 1 :
      a.syllable < b.syllable ? -1 :
      a.syllable > b.syllable ? 1 :
      a.tone - b.tone
    );
  }
};
