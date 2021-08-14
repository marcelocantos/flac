// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';

type MyWindow = typeof window & {
  api: {
    call: (channel: string, ...args: unknown[]) => unknown,
  },
};

// Expose protected methods that allow the renderer process to use
// the ipcRenderer without exposing the entire object
(window as MyWindow).api = {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  call: async (_channel: string, ..._args: unknown[]) => {
    return 42;
  },
};
