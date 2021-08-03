import Lex from './Lex';

test('Test Lex', () => {
	assertLex("", "");
	assertLex("shi4", "shi-4");
	assertLex("shi4de5", "shi-4 de-5");
	assertLex("shi4de5", "shi-4 de-5");
	assertLex("dou1/Du1/du1", "dou-1", "Du-1", "du-1");
	assertLex("dou1/Du1/du1", "dou-1", "Du-1", "du-1");
	assertLex("yi1 kong3 zhi1 jian4", "yi-1 kong-3 zhi-1 jian-4");
	assertLex("xu1yao4/xiang3", "xu-1 yao-4", "xiang-3");
	assertLex("xu1 yao4/xiang3", "xu-1 yao-4", "xiang-3");
	assertLex(" xu1 yao4 / xiang3 ", "xu-1 yao-4", "xiang-3");

	assertLex("jiang14qiang1", "jiang-14 qiang-1");
	assertLex("jiang14/qiang1", "jiang-14", "qiang-1");
});

function assertLex(src: string, ...expected: string[]): boolean {
	const tokenses = Lex(src);
	let actual: string[] = [];
	for (const tokens of tokenses) {
		let chars: string[] = [];
		for (const token of tokens.tokens) {
			chars.push(token.string);
		}
		actual.push(chars.join(' '));
	}
	return expected === actual;
}
