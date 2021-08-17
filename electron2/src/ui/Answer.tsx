import React, { useState, useRef } from 'react';

import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

import refdata from '../refdata/Refdata';

import { 输入RE, 输入字RE } from './InputRE';
import 汉字 from './Word';
import { 条目清单 } from './Decorate';

const validSyllables = refdata.dict.validSyllables;

const validPrefixes = (() => {
  const ret = new Set<string>(['']);
  for (const s in validSyllables) {
    for (let i = 1; i <= s.length; i++) {
      ret.add(s.slice(0, i));
    }
  }
  return ret
})();

interface 回答特性 {
  字: string;
  分数: number;
  提交: (回答: string) => Promise<string | boolean>;
}

export default function 回答({字, 分数, 提交}: 回答特性): JSX.Element {
  const [输入, 设置输入] = useState("");
  const [错误, 设置错误] = useState(false);

  const 回答 = useRef(null);

  function 接受(文字: string): boolean {
    const m = 文字.match(输入RE);
    if (!m || !validPrefixes.has(m[2])) {
      return false;
    }
    for (const [, m] of 文字.matchAll(输入字RE)) {
      if (!validPrefixes.has(m)) {
        return false;
      }
    }
    return true;
  }

  function 检查输入(e: React.ChangeEvent<any>) {
    const 值 = e.target.value as string;
    const 接受的 = !!值 && 接受(值);
    const 错误 = !!值 && !接受的;

    设置错误(错误);
    设置输入(e.target.value);
  }

  async function onClick(e: React.MouseEvent<HTMLElement>) {
    e.preventDefault();
    const 当前 = 回答.current;
    const 结果 = await 提交(当前.value);
    if (typeof 结果 === "string") {
      设置输入(结果);
    } else {
      设置错误(结果);
    }
    当前.focus();
  }

  return (
    <Form>
      <Form.Label htmlFor="回答">
        Enter the pinyin for{' '}
        <汉字 字={字} 分数={分数} 定义={<条目清单 清单={refdata.dict.entries[字]}/>}/>.
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
