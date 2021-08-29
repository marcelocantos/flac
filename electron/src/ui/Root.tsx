import React, { useEffect, useState } from 'react';

import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';

import Assess from '../assess/Assess';
import 随即定义 from '../engine/RandomDefinition';
import 汇报类 from '../engine/Report';
import refdata, { Entries } from '../refdata/Refdata';
import Proxy from '../renderer/data/Proxy';

import 回答 from './Answer';
import { 条目清单 } from './Decorate';
import 结果清单 from './Results';
import 汉字 from './Word';

import './Root.css';

const 记录 = false;

const 条目数据 = refdata.dict.entries;

const 数据 = new Proxy();
const 汇报 = new 汇报类(数据, refdata);

export default function App(): JSX.Element {
  const [字, 设置字] = useState("");
  const [分数, 设置分数] = useState(0);
  const [定义, 设置定义] = useState<string>();
  const [条目组, 设置条目组] = useState<Entries>();
  const [尝试, 设置尝试] = useState(1);

  async function 更新字和分数(新定义: boolean): Promise<void> {
    const {word, score} = await 数据.HeadWord;
    设置字(word);
    设置分数(score ?? 0);
    if (新定义 || typeof 定义 === "undefined") {
      const {定义, 条目组} = 随即定义(word, 条目数据[word]);
      设置定义(定义);
      设置条目组(条目组);
    }
    if (记录) console.log({word, score, 定义, 条目组});
  }

  useEffect(() => {
    更新字和分数(false);
  })

  async function 提交(回答: string): Promise<boolean> {
    const 产物 = Assess(字, 条目组, 回答)
    if (产物.及格) {
      产物.html = ({分数}) => <汉字 字={字} 分数={分数} 定义={<条目清单 清单={条目组}/>}/>;
      await 汇报.好(字, 产物, false);
      await 更新字和分数(true);
      return true;
    } else {
      产物.html = ({分数}) => <汉字 字={字} 分数={分数} 定义={<条目清单 清单={条目组}/>}/>;
      const 尝试包装器 = {尝试};
      await 汇报.不好(产物, false, 尝试包装器);
      设置尝试(尝试包装器.尝试);
      return false;
    }
  }

  const 定义的数目 = 条目组 ? Object.keys(条目组.entries).length : 1;

  return (
    <Container fluid className="Container">
      <Row>
        <p className="welcome">欢迎来到flac，我们一起学中文吧！</p>
      </Row>
      <Col className="main">
        <结果清单 log={[]} streak={[]}/>
      </Col>
      <Row className="input">
        <回答 词={字 || "..."} 分数={分数} 定义={定义} 量={定义的数目} 提交={提交}/>
      </Row>
    </Container>
  );
}
