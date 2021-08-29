import React from 'react';

import 汉字 from './Word';
import 汇报类 from '../engine/Report';

type Props = {
    汇报: 汇报类;
}

export default function 结果清单({汇报}: Props): JSX.Element {
    if (!汇报) {
        return null;
    }
    const 结果组组 = 汇报.历史.concat([汇报.好组]);
    return (
        <div>
            {结果组组.map((结果组, i) =>
                <div key={i} style={{display: "flex", flexWrap: "wrap"}}>
                    {结果组.map((结果, i) =>
                        <React.Fragment key={i}>
                            {结果.不及格 && '❌ '}
                            <汉字 字={结果.Word} 分数={结果.Score} 定义={<></>}/>
                            &nbsp;&nbsp;{' '}
                        </React.Fragment>
                    )}
                </div>
            )}
        </div>
    )
}
