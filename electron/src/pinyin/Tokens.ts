import Token from './Token';

export default class Tokens {
  constructor(
    public readonly tokens: Token[]
  ) {}

  get string(): string {
    return this.tokens.map(t => t.string).join(' ');
  }

  *[Symbol.iterator](): Iterator<Token> {
    for (const token of this.tokens) {
      yield token;
    }
  }
}
