import React from 'react';

import Table from 'react-bootstrap/Table';

import Word from '../pinyin/Word';
import Delim from '../reactutil/Delim';
import { Entries } from '../refdata/Refdata';

import 汉和拼音字 from './HanAndPinyinWord';

const 记录 = false;

type 装饰定义特性 = {
  定义: string;
  不见恶?: string;
  量?: number;
}

export function 装饰定义({定义, 不见恶, 量}: 装饰定义特性): JSX.Element {
  定义 = 定义.replace("'", "’");
  定义 = 定义.replace(/Taiwan pr. /gu, "🇹🇼  ");
  定义 = 定义.replace(/(?:\p{Script=Han}+\|)(\p{Script=Han}+)/gu, "$1");
  定义 = 定义.replace(/\bCL:(\p{Script=Han}+)/gu, "🆑:$1");
  定义 = 定义.replace(/\bclassifier for\b/gu, "🆑➤");
  const segments: JSX.Element[] = [];
  for (let i = 0; ; i++) {
    const m = 定义.match(/^(.*?)(\p{Script=Han}+)?\[((?:(?:🙈|\w+\d)\s+)*(?:🙈|\w+\d))\](.*)/iu);
    if (!m) {
      segments.push(<React.Fragment key={i}>{定义}</React.Fragment>);
      break;
    }
    const [, 前, 汉, 拼音, 后] = m;
    segments.push(
      <React.Fragment key={i}>
        {前}<汉和拼音字 汉={汉} 拼音={拼音}/>
      </React.Fragment>
    );
    定义 = 后;
  }
  return <>{segments}{(量 ?? 1) > 1 && <> (×{量})</>}</>;
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
    new group("classifier for ", "", "🆑➤ "),
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
        grouped[group.first].defs.push(def);
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
          g.group && g.defs.length > 1
          ? <>
              <td/>
              <td>
                {g.group.replace ?? g.group.prefix}
                <Delim delim=", "
                  list={g.defs.map((d, i) =>
                    <索引的 key={i} i={n++}>
                      <装饰定义 定义={d.replace(g.group.regex, '$1')}/>
                    </索引的>
                  )}
                />
                {g.group.suffix}
              </td>
            </>
          : <>
              <td><索引的 i={n++}/></td>
              <td><装饰定义 定义={g.defs[0]}/></td>
            </>
        }</tr>
      )}
    </tbody></table>
  );
}

interface 条目清单特性 {
  清单: Entries,
}

function WordCompare(a: string, b: string): number {
  return Word.compare(new Word(a), new Word(b));
}

export function 条目清单({清单}: 条目清单特性): JSX.Element {
  if (记录) console.log(Object.keys(清单.entries))
  return (
    <Table>
      <tbody>
        {Object.keys(清单.entries).sort(WordCompare).map(条目名 =>
          <tr key={条目名}>
            <th>{<Word.HTML word={条目名}/>}</th>
            <td><定义清单 清单={清单.entries[条目名].definitions}/></td>
          </tr>
        )}
      </tbody>
    </Table>
  );
}
