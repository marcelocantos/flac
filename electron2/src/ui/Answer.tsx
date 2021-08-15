import React, { useState, useRef } from 'react';

import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

import { inputRE, inputCharRE } from './InputRE';

import refdata from '../refdata/Refdata';

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

interface AnswerProps {
  word: string;
  submit: (answer: string) => Promise<string | boolean>;
}

export default function Answer({word, submit}: AnswerProps): JSX.Element {
  const [input, setInput] = useState("");
  const [error, setError] = useState(false);

  const answer = useRef(null);

  function accept(text: string): boolean {
    const m = text.match(inputRE);
    if (!m || !validPrefixes.has(m[2])) {
      return false;
    }
    for (const [, m] of text.matchAll(inputCharRE)) {
      if (!validPrefixes.has(m)) {
        return false;
      }
    }
    return true;
  }

  function checkInput(e: React.ChangeEvent<any>) {
    const value = e.target.value as string;
    const accepted = !!value && accept(value);
    const error = !!value && !accepted;

    setError(error);
    setInput(e.target.value);
  }

  async function onClick(e: React.MouseEvent<HTMLElement>) {
    e.preventDefault();
    const current = answer.current;
    const result = await submit(current.value);
    if (typeof result === "string") {
      setInput(result);
    } else {
      setError(result);
    }
    current.focus();
  }

  return (
    <Form>
      <Form.Label htmlFor="answer">
        Enter the pinyin for <strong>{word}</strong>.
      </Form.Label>
      <InputGroup>
        <InputGroup.Text style={{color: "#666"}}>pinyin&nbsp;→</InputGroup.Text>
        <Form.Control
          id="answer"
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
            disabled={!/\d$/.test(input) || error}
            type="submit"
            onClick={onClick}
            tabIndex={-1}
          >
          Submit
        </Button>
      </InputGroup>
    </Form>
  );
}
