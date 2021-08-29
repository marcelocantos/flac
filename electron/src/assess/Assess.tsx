import Alts from '../pinyin/Alts';
import Word from '../pinyin/Word';
import Pinyin from '../pinyin/Pinyin';
import Lex from '../pinyin/Lex';
import Outcome from '../outcome/Outcome';
import { Entries } from '../refdata/Refdata';

const 记录 = false;

export default function Assess(
	word: string,
	entries: Entries,
	answer: string,
): Outcome {
	if (记录) console.log({word, entries, answer});
	const o = new Outcome(word, entries);
	const answerAlts = AnswerAlts(word, answer);
	if (answerAlts.length > 0) {
		assess(entries, answerAlts, o);
	}
	return o;
}

function AnswerAlts(word: string, answer: string): Alts {
	const tokenses = Lex(answer);
	const words: Word[] = [];
	if (word.length === 1) {
		for (const tokens of tokenses) {
			for (const word of tokens) {
				words.push(...word);
			}
		}
	} else {
		for (const tokens of tokenses) {
			const altses: Word[][] = [];
			for (const token of tokens) {
				altses.push([...token]);
			}
			answerProduct(words, altses);
		}
	}
	return new Alts(words);
}

function answerProduct(words: Word[], altses: Word[][], pinyins: Pinyin[] = []) {
	if (altses.length === 0) {
		words.push(new Word(pinyins));
	} else for (const alt of altses[0]) {
		answerProduct(words, altses.slice(1), pinyins.concat([alt.chars[0]]));
	}
}

function assess(
	entries: Entries,
	answerAlts: Alts,
	o: Outcome,
) {
	const answerMap: {[_: string]: Word} = {};
	for (const alt of answerAlts) {
		answerMap[alt.raw] = alt
	}

	const defMap = new Set<string>(Object.keys(entries.entries));

	const partialDefs = new Set<string>();

	for (const answer in answerMap) {
		const alt = answerMap[answer];
		if (defMap.has(answer)) {
			o.Good.push(alt);
		} else {
			let tooShort = false;
			let badTones = false;
			for (const def in defMap) {
				const word = new Word(def);
				if (alt.length < word.length && alt.raw === word.slice(0, alt.length).raw) {
					partialDefs.add(def);
					tooShort = true;
				} else if (alt.length === word.length) {
					let syllableErrors = 0;
					let tonalErrors = 0;
					word.chars.forEach((p, i) => {
						if (alt.chars[i].syllable !== p.syllable) {
							syllableErrors++;
						}
						if (alt.chars[i].tone !== p.tone) {
							tonalErrors++;
						}
					})
					if (syllableErrors === 0 && tonalErrors > 0) {
						badTones = true;
					}
				}
			}
			if (tooShort) {
				o.TooShort.push(alt);
			} else if (badTones) {
				o.BadTones.push(alt);
				o.Bad.push(alt);
			} else {
				o.Bad.push(alt);
			}
		}
	}

	for (const def in defMap) {
		if (def in answerMap && !partialDefs.has(def)) {
			o.Missing++;
		}
	}
}
