import Alts from './Alts';
import Pinyin from './Pinyin';
import Tones from './Tones';
import Word from './Word';

export default class Token {
	constructor(
		public readonly syllable: string,
		public readonly tones?:   Tones,
	) {}

	get string(): string {
		return `${this.syllable}(${this.tones ? this.tones.string : ''})`
	}

	get alts(): Alts {
		const syllable = this.syllable;
		return new Alts([...this.tones??[]].map(tone =>
			new Word([new Pinyin({syllable, tone})])
		));
	}
}
