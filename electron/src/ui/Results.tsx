import React from 'react';

import 汇报类 from '../engine/Report';
import Word from '../pinyin/Word';
import Alts from '../pinyin/Alts';

import { 条目清单 } from './Decorate';
import 汉字 from './Word';

type Props = {
  汇报: 汇报类;
}

export default function 结果清单({汇报}: Props): JSX.Element {
  if (!汇报) {
    return null;
  }
  let 结果组组 = 汇报.历史;
  if (汇报.好组.length > 0) {
    结果组组 = 结果组组.concat([汇报.好组]);
  }
  return (
    <div>
      {结果组组.map((结果组, i) =>
        <div key={i} style={{display: "flex", flexWrap: "wrap"}}>
          {结果组[0].不及格
            ? 结果组.map((结果, i) => {
                let entry: JSX.Element;
                try {
                  entry = Word.HTML({word: 结果.Answer});
                } catch {
                  entry = Alts.HTML({alts: 结果.Answer});
                }
                return (
                  <React.Fragment key={i}>
                    ❌&nbsp;&nbsp;
                    {entry}
                    &nbsp; ≠ &nbsp;
                    <汉字
                      className="bad"
                      字={结果.Word}
                      分数={结果.Score}
                      定义={<条目清单 清单={结果.Entries}/>}
                    />
                  </React.Fragment>
                );
              })
            : 结果组.map((结果, i) =>
                <React.Fragment key={i}>
                  <汉字
                    className="good"
                    字={结果.Word}
                    分数={结果.Score}
                    定义={<条目清单 清单={结果.Entries}/>}
                  />
                  &nbsp;&nbsp;{' '}
                </React.Fragment>
              )
          }
        </div>
      )}
    </div>
  )
}
