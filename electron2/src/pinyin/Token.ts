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

	*[Symbol.iterator](): Iterator<Word> {
		const syllable = this.syllable;
		for (const tone of this.tones??[]) {
			yield new Word([new Pinyin({syllable, tone})]);
		}
	}
}
