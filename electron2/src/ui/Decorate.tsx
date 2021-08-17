import React from 'react';

import Table from 'react-bootstrap/Table';

import Pinyin from '../pinyin/Pinyin';
import Delim from '../reactutil/Delim';
import { Entries } from '../refdata/Refdata';

import 汉和拼音字 from './HanAndPinyinWord';

export function 定义({def}: {def: string}): JSX.Element {
	def = def.replace("'", "’");
	def = def.replace(/Taiwan pr. /gu, "🇹🇼  ");
	def = def.replace(/(?:\p{Script=Han}+\|)(\p{Script=Han}+)/gu, "$1");
	def = def.replace(/\bCL:(\p{Script=Han}+)/gu, "🆑:$1");
	def = def.replace(/\bclassifier for\b/gu, "🆑➤");
	const segments: JSX.Element[] = [];
	for (let i = 0; ; i++) {
		const m = def.match(/^(.*?)(\p{Script=Han}+)?\[((?:\w+\d\s+)*\w+\d)\](.*)/iu);
		if (!m) {
			segments.push(<React.Fragment key={i}>{def}</React.Fragment>);
			break;
		}
		const [, 前, 汉, 拼音, 后] = m;
		segments.push(
			<React.Fragment key={i}>
				{前}<汉和拼音字 汉={汉} 拼音={拼音}/>
			</React.Fragment>
		);
		def = 后;
	}
	return <>{segments}</>;
}

export function 索引的({i, children}: {i: number, children?: React.ReactNode}): JSX.Element {
	return <>{' '}<sup className="定义序数词">{i}</sup>{children}</>;
}

export function 定义清单({清单}: {清单: string[]}): JSX.Element {
	class group {
		readonly regex: RegExp;

		constructor(
			public readonly prefix:  string = "",
			public readonly suffix:  string = "",
			public readonly replace?: string,
			public first:   number = -1,
		){
			this.regex = new RegExp(`^${prefix}(.*)${suffix}$`)
		}
	}

	const groups: group[] = [
		new group("to ", "", "to… "),
		new group("abbr. for ", "", "abbr… "),
		new group("classifier for "),
		new group("(grammatical equivalent of ", ")", "(gramm ≣… "),
		new group("(indicates ", ")", "(indic… "),
	];

	const grouped: {group?: group, defs: string[]}[] = [];
	for (const def of 清单) {
		let matched = false;
		for (const group of groups) {
			if (def.match(group.regex)) {
				matched = true;
				if (group.first < 0) {
					group.first = grouped.length;
					grouped.push({group, defs: []});
				}
				grouped[group.first].defs.push(def.replace(group.regex, '$1'));
				break;
			}
		}
		if (!matched) {
			grouped.push({defs: [def]});
		}
	}

	let n = 1;

	return (
		<table className="定义清单"><tbody>
			{grouped.map((g, i) =>
				<tr key={i}>{
					g.group
					? <>
							<td/>
							<td>
								{g.group.replace ?? g.group.prefix}
								<Delim delim=", "
									list={g.defs.map((d, i) =>
										<索引的 key={i} i={n++}><定义 def={d}/></索引的>
									)}
								/>
								{g.group.suffix}
							</td>
						</>
					: <>
							<td><索引的 i={n++}/></td>
							<td><定义 def={g.defs[0]}/></td>
						</>
				}</tr>
			)}
		</tbody></table>
	);
}

interface 条目清单特性 {
	清单: Entries,
}

export function 条目清单({清单}: 条目清单特性): JSX.Element {
	return (
		<Table>
			<tbody>
				{Object.keys(清单.entries).sort().map(条目名 =>
					<tr key={条目名}>
						<th>{<Pinyin.HTML pinyin={条目名}/>}</th>
						<td><定义清单 清单={清单.entries[条目名].definitions}/></td>
					</tr>
				)}
			</tbody>
		</Table>
	);
}
