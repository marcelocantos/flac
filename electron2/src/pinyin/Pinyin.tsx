import React from 'react';

import PinyinRE from './PinyinRE';

import './Pinyin.css';

const baseVowels = 'AEIOUÜaeiouü';

const mark: ((v: string) => string)[] = [
  '',
  'ĀĒĪŌŪǕāēīōūǖ',
  'ÁÉÍÓÚǗáéíóúǘ',
  'ǍĚǏǑǓǙǎěǐǒǔǚ', // Breve forms: ăĕĭŏŭ-ĂĔĬŎŬ-
  'ÀÈÌÒÙǛàèìòùǜ',
  baseVowels,
].map(vowels => (v: string) => vowels[baseVowels.indexOf(v)] ?? '')

interface PinyinProps {
  pinyin: Pinyin | string | {syllable: string, tone: number};
  [attr: string]: unknown;
}

export default class Pinyin {
  readonly pinyin: string;
  readonly syllable: string;
  readonly tone: number;
  readonly raw: string;
  readonly consumed?: number;

  constructor(pinyin: string | {syllable: string, tone: number}) {
    if (typeof pinyin === 'string') {
      const groups = PinyinRE.exec(pinyin);
      if (groups == null || (!groups[1] && groups[3].length > 1)) {
        throw new Error(`${pinyin}: invalid pinyin`);
      }
      this.consumed = groups[0].length;
      pinyin = {
        syllable: groups[2].replace('v', 'ü').replace('u:', 'ü'),
        tone: Number.parseInt(groups[3]),
      }
    }

    const {syllable, tone} = pinyin;

    // https://en.wikipedia.org/wiki/Pinyin#Rules_for_placing_the_tone_mark
    this.pinyin   = syllable.replace(/[aeo]|(?<=i)u|(?<=u)i|[iuü]/i, mark[tone]);
    this.syllable = syllable.replace(/ü/i, 'v');
    this.tone     = tone;
    this.raw      = `${this.syllable}${this.tone}`;
  }

  html(props: {[attr: string]: unknown}): JSX.Element {
    return <span {...props} className={`拼音 拼音调${this.tone}`}>{this.pinyin}</span>;
  }

  static HTML({pinyin, ...props}: PinyinProps): JSX.Element {
    if (!(pinyin instanceof Pinyin)) {
      pinyin = new Pinyin(pinyin);
    }
    // Force typecast because type-narrowing failed.
    return (pinyin as Pinyin).html(props);
  }

  static compare(a: Pinyin, b: Pinyin): number {
    const aLower = a.syllable.toLowerCase();
    const bLower = b.syllable.toLowerCase();
    if (aLower < bLower) {
      return -1;
    } else if (aLower > bLower) {
      return 1;
    } else if (a.syllable < b.syllable) {
      return -1;
    } else if (a.syllable > b.syllable) {
      return 1;
    } else {
      return a.tone - b.tone
    }
  }
}
