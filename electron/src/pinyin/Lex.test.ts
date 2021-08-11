import Lex from './Lex';

it('Lex', () => {
	expect(lex("")).toEqual("");
	expect(lex("shi4")).toEqual("shi(4)");
	expect(lex("shi4de5")).toEqual("shi(4) de(5)");
	expect(lex("shi4de5")).toEqual("shi(4) de(5)");
	expect(lex("dou1/Du1/du1")).toEqual("dou(1) / Du(1) / du(1)");
	expect(lex("dou1/Du1/du1")).toEqual("dou(1) / Du(1) / du(1)");
	expect(lex("yi1 kong3 zhi1 jian4")).toEqual("yi(1) kong(3) zhi(1) jian(4)");
	expect(lex("xu1yao4/xiang3")).toEqual("xu(1) yao(4) / xiang(3)");
	expect(lex("xu1 yao4/xiang3")).toEqual("xu(1) yao(4) / xiang(3)");
	expect(lex(" xu1 yao4 / xiang3 ")).toEqual("xu(1) yao(4) / xiang(3)");

	expect(lex("jiang14qiang1")).toEqual("jiang(14) qiang(1)");
	expect(lex("jiang14/qiang1")).toEqual("jiang(14) / qiang(1)");

	expect(lex("yi1 shi4 yi1 , er4 shi4 er4"))
		.toEqual("yi(1) shi(4) yi(1) ,() er(4) shi(4) er(4)");

	expect(lex("Ya4 dang1 · Si1 mi4")).toEqual("Ya(4) dang(1) ·() Si(1) mi(4)")
});

it('bad inputs', () => {
	expect(() => lex("1")).toThrow();
	expect(() => lex("shi1-de5")).toThrow();
});

function lex(src: string): string {
	const tokenses = Lex(src);
	let parts: string[] = [];
	for (const tokens of tokenses) {
		let chars: string[] = [];
		for (const token of tokens.tokens) {
			chars.push(token.string);
		}
		parts.push(chars.join(' '));
	}
	return parts.join(" / ");
}
