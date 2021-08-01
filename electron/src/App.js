import './App.css';

import React, { useState, useRef } from 'react';

import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Row from 'react-bootstrap/Row';

import refdata from './refdata.json';

console.log({refdata});

function App() {
  const [input, setInput] = useState("");
  const [error, setError] = useState(false);

  const answer = useRef(null);

  function checkInput(e) {
    const value = e.target.value;
    const error = !/^([a-z]+\d)*([a-z]+\d?)?$/.test(value);
    setError(error);
    // if (!error) {
      setInput(e.target.value);
    // }
  }

  function submit(e) {
    e.preventDefault();
    console.log(answer.current.value);
    setError(answer.current.value !== "sheng1huo2");
  }

  return (
    <Container fluid className="Container">
      <Row>
        <h1>flac: learn 中文</h1>
      </Row>
      <Col className="main">
        <p className="welcome">欢迎来到flac，一起学中文吧！</p>
      </Col>
      <Row className="input">
        <Form>
          <Form.Label>Enter the pinyin for <strong>生活</strong>.</Form.Label>
          <InputGroup>
            {/* <InputGroup.Text id="input">生活</InputGroup.Text> */}
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
            <Button disabled={error} type="submit" onClick={submit}>
              Submit
            </Button>
          </InputGroup>
        </Form>
      </Row>
    </Container>
  );
}

export default App;
