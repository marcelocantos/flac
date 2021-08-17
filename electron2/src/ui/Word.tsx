import React from 'react';

import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Popover from 'react-bootstrap/Popover';

import 分数条 from './Score';

interface Props {
  字: string;
  分数: number;
  定义: JSX.Element;
}

export default function 汉字({字, 分数, 定义}: Props): JSX.Element {
  return (
    <OverlayTrigger
      trigger="click"
      placement="top"
      overlay={props =>
        <Popover className="定义弹出框" {...props} id={`汉字-{字}`}>{定义}</Popover>
      }
    >
      <span className="字">{字}<分数条 分数={分数}/></span>
    </OverlayTrigger>
  );
}
