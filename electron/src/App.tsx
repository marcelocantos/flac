import './App.css';

import React, { useState, useRef } from 'react';

import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Row from 'react-bootstrap/Row';

import Results from './Results';

import refdata from './Refdata';

console.log({refdata});

const words = refdata.wordList.words;
const entries = refdata.dict.entries;

function App() {
  const [wordIndex, setWordIndex] = useState(0);
  const [input, setInput] = useState("");
  const [error, setError] = useState(false);

  const answer = useRef(null);

  const word = words[wordIndex];
  const entry = entries[word];
  console.log({word, entry});

  function checkInput(e: React.ChangeEvent<any>) {
    const value = e.target.value;
    const error = !/^([a-z]+\d)*([a-z]+\d?)?$/.test(value);
    setError(error);
    setInput(e.target.value);
  }

  function submit(e: React.MouseEvent<HTMLElement>) {
    e.preventDefault();
    const current = answer.current as any;
    if (current) {
      if (current.value in entry.entries) {
        setWordIndex(wordIndex + 1);
        setInput("");
      } else {
        setError(true);
      }
      current.focus();
    }
  }

  return (
    <Container fluid className="Container">
      <Row>
        <p className="welcome">欢迎来到flac，一起学中文吧！</p>
      </Row>
      <Col className="main">

        <Results log={[]} streak={[]}/>
      </Col>
      <Row className="input">
        <Form>
          <Form.Label id="prompt">
            Enter the pinyin for <strong>{words[wordIndex]}</strong>.
          </Form.Label>
          <InputGroup>
            <InputGroup.Text style={{color: "#666"}}>pinyin&nbsp;→</InputGroup.Text>
            <Form.Control
              ref={answer}
              isInvalid={error}
              autoFocus={true}
              value={input}
              onChange={checkInput}
              size="lg"
              spellCheck={false}
              aria-label="answer"
              aria-describedby="input"
            />
            <Button
                disabled={error}
                type="submit"
                onClick={submit}
                tabIndex={-1}
              >
              Submit
            </Button>
          </InputGroup>
        </Form>
      </Row>
    </Container>
  );
}

export default App;
