import React, { useEffect, useState } from 'react';

import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';

import Results from './ui/Results';
import Answer from './ui/Answer';

import refdata from './refdata/Refdata';

import { Database } from './renderer/data/Proxy';

import './Root.css';

const entries = refdata.dict.entries;

const data = new Database();

export default function App(): JSX.Element {
  const [word, setWord] = useState("");

  const entry = entries[word];

  useEffect(() => {
    (async () => {
      setWord(await data.HeadWord);
    })();
  })

  async function submit(answer: string): Promise<string | boolean> {
    if (answer in entry.entries) {
      await data.MoveWord(word, 5);
      setWord(await data.HeadWord);
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
        <Results log={[]} streak={[]}/>
      </Col>
      <Row className="input">
        <Answer word={word || "..."} submit={submit}/>
      </Row>
    </Container>
  );
}
