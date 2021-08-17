import React, { useEffect, useState } from 'react';

import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';

import 结果清单 from './Results';
import 回答 from './Answer';

import refdata from '../refdata/Refdata';

import { Database } from '../renderer/data/Proxy';

import './Root.css';

const log = true;

const entries = refdata.dict.entries;

const data = new Database();

export default function App(): JSX.Element {
  const [word, setWord] = useState("");
  const [score, setScore] = useState(0);

  const entry = entries[word];

  async function updateWordAndScore(): Promise<void> {
    const {word, score} = await data.HeadWord;
    if (log) console.log({word, score});
    setWord(word);
    setScore(score ?? 0);
  }

  useEffect(() => {
    updateWordAndScore();
  })

  async function submit(answer: string): Promise<string | boolean> {
    if (answer in entry.entries) {
      await data.UpdateScoreAndPos(word, 1 + 2*score, 5);
      await updateWordAndScore();
      return "";
    } else {
      return true;
    }
  }

  return (
    <Container fluid className="Container">
      <Row>
        <p className="welcome">欢迎来到flac，我们一起学中文吧！</p>
      </Row>
      <Col className="main">
        <结果清单 log={[]} streak={[]}/>
      </Col>
      <Row className="input">
        <回答 字={word || "..."} 分数={score} 提交={submit}/>
      </Row>
    </Container>
  );
}
