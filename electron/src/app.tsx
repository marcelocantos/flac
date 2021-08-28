import React from 'react';
import * as ReactDOM from 'react-dom';

import Root from './ui/Root';

function render() {
  ReactDOM.render(
    <React.StrictMode>
      <Root />
    </React.StrictMode>,
    document.getElementById('root'),
  );
}

render();