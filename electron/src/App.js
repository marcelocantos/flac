import logo from './logo.svg';
import './App.css';

import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import InputGroup from 'react-bootstrap/InputGroup';
import FormControl from 'react-bootstrap/FormControl';
import Row from 'react-bootstrap/Row';

function App() {
  return (
    <Container fluid className="Container">
      <Row>
        <h1>flac: learn 中文</h1>
      </Row>
      <Col className="main">
        <p class="welcome">欢迎来到flac，一起学中文吧！</p>
      </Col>
      <Row className="input">
        <InputGroup>
          <InputGroup.Text id="input">生活</InputGroup.Text>
          <FormControl
            placeholder="拼音"
            aria-label="pinyin"
            aria-describedby="input"
          />
          <Button>Submit</Button>
        </InputGroup>
      </Row>
    </Container>
  );
}

export default App;
