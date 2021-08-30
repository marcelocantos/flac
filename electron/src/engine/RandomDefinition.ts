import { Entries } from '../refdata/Refdata';

const 记录 = false;

function 随机选择<T>(数组: T[]): T {
  return 数组[Math.floor(Math.random() * 数组.length)];
}

export default function 随机定义(
  汉字: string,
  条目组: Entries,
): {定义: string, 条目组: Entries} {
  if (记录) console.log({随机定义: {汉字, 条目组}});
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
  let 量词 = -1;
  const 看RE = new RegExp(
    `^(?:also written |also pr. |CL:)|^(?:(?:unofficial )?variant of|see) .*\\b${拼音}\\b`,
    'iu');
  const pinyinRE = new RegExp(`\\b${拼音}\\b`, 'giu');
  for (let 定义 of 定义清单.definitions) {
    定义 = 定义.replaceAll(pinyinRE, "🙈");
    if (记录) console.log({定义, 拼音});
    if (定义.match(看RE)) {
      if (记录) console.log('看RE');
      见 = 候选定义.length;
      候选定义.push(定义);
    } else if (定义.startsWith("CL:")) {
      if (记录) console.log('CL:');
      量词 = 候选定义.length;
      候选定义.push(定义);
    } else if (定义.startsWith("surname ")) {
      if (记录) console.log('surname');
      候选定义.push("surname");
    } else {
      候选定义.push(定义);
    }
  }

  // "CL:..." aren't great choices of definitions to test. Avoid
  // unless no other options reman.
  if (量词 != -1 && 候选定义.length > 1) {
    候选定义.splice(量词, 1);
  }

  // "(see|variant of) ..." aren't great choices of definitions to test. Avoid
  // unless no other options reman.
  if (见 != -1 && 候选定义.length > 1) {
    候选定义.splice(见, 1);
  }

  if (候选定义.length === 0) {
    throw new Error(`no useful definitions for ${汉字}: ${定义清单}`);
  }

  const 定义 = 随机选择(候选定义);
  for (const 拼音2 of 拼音清单) {
    if (拼音2 !== 拼音) {
      for (const 定义2 of 条目组.entries[拼音2].definitions) {
        if (定义 === "surname" ? 定义2.startsWith("surname ") : 定义2 === 定义) {
          新条目.entries[拼音2] = 条目组.entries[拼音2];
          break;
        }
      }
    }
  }
  return {定义, 条目组: 新条目};
}
