import { Entries } from '../refdata/Refdata';

function 随机选择<T>(数组: T[]): T {
  return 数组[Math.floor(Math.random() * 数组.length)];
}

export default function 随机定义(
  汉字: string,
  条目组: Entries,
): {定义: string, 条目组: Entries} {
  const 拼音清单 = Object.keys(条目组.entries);
  if (拼音清单.length === 1) {
    return {定义: "", 条目组: 条目组};
  }

  const 新条目: Entries = {...条目组, entries: {}};

  const 拼音 = 随机选择(拼音清单);
  const 定义清单 = 条目组.entries[拼音];
  新条目.entries = {[拼音]: 定义清单};

  const 候选定义: string[] = [];
  let 见 = -1;
  const 正则表达式 = new RegExp(
    `^(?:also written |also pr. |CL:)|^(?:(?:unofficial )?variant of|see) .*\\b${拼音}\\b`,
    'iu');
  const pinyinRE = new RegExp(`\\b${拼音}\\b`, 'giu');
  for (let 定义 of 定义清单.definitions) {
    console.log({定义, 拼音});
    定义 = 定义.replaceAll(pinyinRE, "🙈");
    if (定义.match(正则表达式)) {
      候选定义.push(定义);
      见 = 候选定义.length;
    } else if (定义.startsWith("surname ")) {
      候选定义.push("surname");
    } else {
      候选定义.push(定义);
    }
  }

  // "(see|variant of) ..." aren't great choices of definitions to test. Avoid
  // unless no other options reman.
  if (见 != -1 && 候选定义.length > 1) {
    候选定义.splice(见, 1);
  }

  if (候选定义.length === 0) {
    throw new Error(`no useful definitions for ${汉字}: ${定义清单}`);
  }
  return {定义: 随机选择(候选定义), 条目组: 新条目};
}
