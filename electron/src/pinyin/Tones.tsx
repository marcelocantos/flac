export default class Tones {
	readonly tones: number;
	readonly length: number;

	constructor(s: string) {
		this.tones = 0
		for (const c of s) {
			this.tones |= 1 << Number.parseInt(c);
		}
		let i = 0;
		for (let t = this.tones; t > 0; t &= t - 1) {
			i++
		}
		this.length = i;
	}

	get string(): string {
		return [...this].join('');
	}

	*[Symbol.iterator]() {
		for (let i = 1, t = this.tones >> 1; i <= 5; i++, t >>= 1) {
			if (t & 1) {
				yield i;
			}
		}
	}
}
