import { fireEvent, render } from '@testing-library/react';
import App from './App';

it('renders learn react link', () => {
  const app = render(<App />);
  const linkElement = app.getByText(/欢迎来到flac，我们一起学中文吧！/i);
  expect(linkElement).toBeInTheDocument();

  const answerLabel = app.getByText(/^enter the /i);
  const answerInput = app.getByLabelText(/^enter the /i);
  const submitButton = app.getByText(/^submit$/i);

  function key(value: string) {
    fireEvent.change(answerInput, {
      target: {value: value},
    });
  }

  function submit() {
    fireEvent.click(submitButton);
  }

  // Submit the correct pinyin for 第.
  expect(answerInput).toHaveValue('');
  key('di4');
  expect(answerInput).toHaveValue('di4');
  submit();
  expect(answerInput).toHaveValue('');

  // Label should now be prompting for 的.
  expect(answerLabel).toHaveTextContent(/的/);

  // Again, submit the correct pinyin for 第, which is also valid for 的.
  key('di4');
  submit();
  expect(answerInput).toHaveValue(''); // Should remain unchanged after submit.

  expect(answerLabel).toHaveTextContent(/了/);

  // Submit the correct pinyin for 第 again, but this time it should have been
  // the pinyin for 了.
  key('di4');
  expect(answerInput).toHaveValue('di4');
  submit();
  expect(answerInput).toHaveValue('di4'); // Should remain unchanged after submit.

  // Now get it right.
  key('le5');
  submit();
  expect(answerInput).toHaveValue(''); // Should remain unchanged after submit.
});
