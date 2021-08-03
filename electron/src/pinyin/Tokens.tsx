import Token from './Token';

export default class Tokens {
	constructor(
		public readonly tokens: Token[]
	) {}

	get string(): string {
		return this.tokens.map(t => t.string).join(' ');
	}
}
