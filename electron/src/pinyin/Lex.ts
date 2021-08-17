import Token from './Token';
import Tokens from './Tokens';
import Tones from './Tones';
import PinyinRE from './PinyinRE';

export default function Lex(raw: string): Tokens[] {
	const ret: Tokens[] = [];
	let tokens: Token[] = [];
	while (raw) {
		const groups = raw.match(PinyinRE);
		if (!groups) {
			throw new Error(`${raw}: invalid pinyin`)
		}
		switch (groups[1]) {
		case undefined:
			tokens.push(new Token(groups[2], new Tones(groups[3])));
			break;
		case "/":
			ret.push(new Tokens(tokens));
			tokens = [];
			break;
		default:
			tokens.push(new Token(groups[1]))
		}
		raw = raw.slice(groups[0].length);
	}
	ret.push(new Tokens(tokens));
	return ret
}
