import React, { useState, useRef } from 'react';

import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

import refdata from '../refdata/Refdata';

import { 输入RE, 输入字RE } from './InputRE';
import 汉字 from './Word';
import * as Decorate from './Decorate';

const 有效音节 = refdata.dict.validSyllables;

const 有效前缀 = (() => {
  const ret = new Set<string>(['']);
  for (const s in 有效音节) {
    for (let i = 1; i <= s.length; i++) {
      ret.add(s.slice(0, i));
    }
  }
  return ret
})();

type 回答特性 = {
  词: string;
  分数: number;
  定义?: string;
  量?: number;
  提交: (回答: string) => Promise<boolean>;
}

export default function 回答({词, 分数, 定义, 量, 提交}: 回答特性): JSX.Element {
  const [输入, 设置输入] = useState("");
  const [错误, 设置错误] = useState(false);

  const 回答 = useRef(null);

  function 接受(文字: string): boolean {
    const 匹配 = 文字.match(输入RE);
    return 匹配 && 有效前缀.has(匹配[2]) && ![...文字.matchAll(输入字RE)].some(([, 匹配]) => !有效前缀.has(匹配));
  }

  function 检查输入(事件: React.ChangeEvent<any>) {
    const 值 = 事件.target.value as string;
    const 接受的 = !!值 && 接受(值);
    const 错误 = !!值 && !接受的;

    设置输入(事件.target.value);
    设置错误(错误);
  }

  async function onClick(事件: React.MouseEvent<HTMLElement>) {
    事件.preventDefault();
    const 当前 = 回答.current;
    if (await 提交(当前.value)) {
      设置输入("");
    } else {
      设置错误(true);
    }
    当前.focus();
  }

  return (
    <Form>
      <Form.Label htmlFor="回答">
        Enter the pinyin for{' '}
        <汉字 字={词} 分数={分数}
          定义={<Decorate.条目清单 清单={refdata.dict.entries[词]}/>}
        />
        {定义 && <>&nbsp;&nbsp;👉&nbsp;&nbsp;<Decorate.装饰定义 定义={定义} 不见恶={词} 量={量}/></>}
        .
      </Form.Label>
      <InputGroup>
        <InputGroup.Text style={{color: "#666"}}>拼音 →</InputGroup.Text>
        <Form.Control
          id="回答"
          ref={回答}
          isInvalid={错误}
          autoFocus={true}
          value={输入}
          onChange={检查输入}
          size="lg"
          spellCheck={false}
          aria-label="回答"
          aria-describedby="input"
        />
        <Button
            disabled={!/\d$/.test(输入) || 错误}
            type="submit"
            onClick={onClick}
            tabIndex={-1}
          >
          提交
        </Button>
      </InputGroup>
    </Form>
  );
}
