import React from 'react';
import { fireEvent, render, act } from '@testing-library/react';
import Root from './Root';

it('renders learn react link', async () => {
  const app = render(<Root />);
  const linkElement = app.getByText(/欢迎来到flac，我们一起学中文吧！/i);
  expect(linkElement).toBeInTheDocument();

  const answerLabel = app.getByText(/^enter the /i);
  const answerInput = app.getByLabelText(/^enter the /i);
  const submitButton = app.getByText(/^submit$/i);

  async function key(value: string): Promise<void> {
    await act(async () => {
      fireEvent.change(answerInput, {
        target: {value: value},
      });
    });
  }

  async function submit(): Promise<void> {
    await act(async () => {
      fireEvent.click(submitButton);
    });
  }

  // Submit the correct pinyin for 第.
  expect(answerInput).toHaveValue('');
  await key('di4');
  expect(answerInput).toHaveValue('di4');
  await submit();
  expect(answerInput).toHaveValue('');

  // Label should now be prompting for 的.
  expect(answerLabel).toHaveTextContent(/的/);

  // Again, submit the correct pinyin for 第, which is also valid for 的.
  await key('di4');
  await submit();
  expect(answerInput).toHaveValue(''); // Should remain unchanged after submit.

  expect(answerLabel).toHaveTextContent(/了/);

  // Submit the correct pinyin for 第 again, but this time it should have been
  // the pinyin for 了.
  await key('di4');
  expect(answerInput).toHaveValue('di4');
  await submit();
  expect(answerInput).toHaveValue('di4'); // Should remain unchanged after submit.

  // Now get it right.
  await key('le5');
  await submit();
  expect(answerInput).toHaveValue(''); // Should remain unchanged after submit.
});
