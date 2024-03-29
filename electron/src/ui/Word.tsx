import React from 'react';

import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Popover from 'react-bootstrap/Popover';

import 分数条 from './Score';

type Props = {
  className?: string;
  字: string;
  分数: number;
  定义: JSX.Element;
}

export default function 汉字({className, 字, 分数, 定义}: Props): JSX.Element {
  className = (className ? className + " 字" : "字");
  const 字分数 = <span className={className}>{字}<分数条 分数={分数}/></span>;
  return (
    定义
    ? (
      <OverlayTrigger
        placement="auto"
        overlay={props =>
          <Popover className="定义弹出框" {...props} id={`汉字-{字}`}>
            <Popover.Body>{定义}</Popover.Body>
          </Popover>
        }
      >
        {字分数}
      </OverlayTrigger>
    )
    : 字分数
  );
}
