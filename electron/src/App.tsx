import './App.css';

import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';

import Results from './ui/Results';
import Answer from './ui/Answer';

export default function App() {
  return (
    <Container fluid className="Container">
      <Row>
        <p className="welcome">欢迎来到flac，一起学中文吧！</p>
      </Row>
      <Col className="main">
        <Results log={[]} streak={[]}/>
      </Col>
      <Row className="input">
        <Answer/>
      </Row>
    </Container>
  );
}
