import * as AsyncDB from './AsyncDB';
import Database from './Database';

const words = [
	"第", "的", "了", "在", "是", "我", "和", "有", "你", "个", "也", "这", "不",
	"他", "上", "人", "中", "就", "年", "为", "对", "说", "都", "要", "到", "着",
	"住", "与", "将", "日", "我们", "好", "月", "会", "大", "来", "还", "等", "而",
	"地", "自己", "后", "两", "一", "被", "没有", "去", "但", "从", "很", "给", "时",
	"以", "中国",
];

it('build and close', async () => {
	const d = await prepareDatabase();
	expect(await d.MaxPos).toEqual(words.length - 1);
	for (let i = 0; i < words.length; ++i) {
		expect(await d.WordPos(words[i])).toEqual(i);
	}
	await d.close();
});

it('build partial', async () => {
	const d = await prepareDatabase(undefined, words.slice(0, 10));
	expect(await d.MaxPos).toEqual(9);
	for (let i = 0; i < 10; ++i) {
		expect(await d.WordPos(words[i])).toEqual(i);
	}
	await d.close();
});

it('build incrementally', async () => {
	const db = await AsyncDB.open(":memory:");
	await prepareDatabase(db, words.slice(0, 10));
	db.finalize();
	const d = await prepareDatabase(db, words);
	expect(await d.MaxPos).toEqual(words.length - 1);
	for (let i = 0; i < words.length; ++i) {
		expect(await d.WordPos(words[i])).toEqual(i);
	}
	await d.close();
});

it('build different words', async () => {
	const db = await AsyncDB.open(":memory:");
	await prepareDatabase(db, words.slice(20));
	db.finalize();
	const d = await prepareDatabase(db, words.slice(0, 30));
	expect(await d.MaxPos).toEqual(29);
	await d.close();
});

it('WordAt', async () => {
	const d = await prepareDatabase();

	expect(await d.WordAt(0)).toEqual("第");
	expect(await d.HeadWord).toEqual({word: "第", score: undefined});
	expect(await d.WordAt(1)).toEqual("的");
});

it('WordPos', async () => {
	const d = await prepareDatabase();

	expect(await d.WordPos("第")).toEqual(0);
	expect(await d.WordPos("元")).toEqual(undefined);
});

it('MoveWord', async () => {
	const d = await prepareDatabase();
	await expectWordsAt(d, 0, ["第", "的", "了"]);
	await d.MoveWord("第", 1); await expectWordsAt(d, 0, ["的", "第", "了"]);
	await d.MoveWord("的", 3); await expectWordsAt(d, 0, ["第", "了", "在", "的"]);
	await d.MoveWord("了", 3); await expectWordsAt(d, 0, ["第", "在", "的", "了"]);
	await d.MoveWord("了", 1); await expectWordsAt(d, 0, ["第", "了", "在", "的"]);
	await d.MoveWord("的", 0); await expectWordsAt(d, 0, ["的", "第", "了", "在"]);
	await d.MoveWord("第", 0); await expectWordsAt(d, 0, ["第", "的", "了", "在"]);
});

it('MoveWord from end', async () => {
	const d = await prepareDatabase();
	await d.MoveWord("中国", 10);
	await expectWordsAt(d, 0, words.slice(0, 10).concat(["中国"], words.slice(10, -1)));
});

it('MoveWord to end', async () => {
	const d = await prepareDatabase();
	await d.MoveWord("也", words.length - 1);
	await expectWordsAt(d, 0, words.slice(0, 10).concat(words.slice(11), ["也"]));
});

it('MoveWord past end', async () => {
	const d = await prepareDatabase();
	await d.MoveWord("也", words.length+99);
	await expectWordsAt(d, 0, words.slice(0, 10).concat(words.slice(11), ["也"]));
});

it('MoveWord from start to end', async () => {
	const d = await prepareDatabase();
	await d.MoveWord("第", words.length - 1);
	await expectWordsAt(d, 0, words.slice(1).concat(["第"]));
});

it('MoveWord from end to start', async () => {
	const d = await prepareDatabase();
	await d.MoveWord("中国", 0);
	await expectWordsAt(d, 0, ["中国"].concat(words.slice(0, -1)));
});

it('MoveWord empty queue', async () => {
	const d = await prepareDatabase();
	await expect(d.MoveWord("元", 10)).rejects.toThrow();
});

it('MoveWord missing word', async () => {
	const d = await prepareDatabase();
	await expect(d.MoveWord("元", 10)).rejects.toThrow();
});

it('MaxScore no scores', async () => {
	const d = await prepareDatabase();
	expect(await d.MaxScore).toEqual(null);
});

it('MaxScore', async () => {
	const d = await prepareDatabase();
	await d.UpdateScoreAndPos("第", 42, 10);
	expect(await d.MaxScore).toEqual(42);
	expect(await d.MaxScore).toEqual(42);
	expect(await d.WordScore("第")).toEqual(42);
	expect(await d.MaxScore).toEqual(42);
	// expect(await d.WordPos("第")).toEqual(10);
	expect(await d.MaxScore).toEqual(42);

	await d.UpdateScore("的", 56);
	expect(await d.MaxScore).toEqual(56);

	await d.UpdateScore("了", 20);
	expect(await d.MaxScore).toEqual(56);

	await d.UpdateScore("的", 15);
	expect(await d.MaxScore).toEqual(42);

	await d.UpdateScore("第", 10);
	expect(await d.MaxScore).toEqual(20);
});

async function prepareDatabase(db?: AsyncDB.Database, w?: string[]): Promise<Database> {
	return await Database.build(db || await AsyncDB.open(":memory:"), "", w || words);
}

async function expectWordsAt(d: Database, begin: number, words: string[]) {
	for (let i = 0; i < words.length; ++i) {
		const word = words[i];
		const pos = await d.WordPos(word);
		expect(pos).toEqual(begin + i);
	}
}
