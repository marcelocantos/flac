import { useState } from 'react';

import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';

import Results from './ui/Results';
import Answer from './ui/Answer';

import refdata from './refdata/Refdata';

import './App.css';

const entries = refdata.dict.entries;
const words = refdata.wordList.words;

export default function App() {
  const [wordIndex, setWordIndex] = useState(0);

  const word = words[wordIndex];
  const entry = entries[word];

  function submit(answer: string): string | boolean {
    if (answer in entry.entries) {
      setWordIndex(wordIndex + 1);
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
        <Answer word={word} submit={submit}/>
      </Row>
    </Container>
  );
}
