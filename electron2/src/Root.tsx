import React, { useState, useEffect } from 'react';

type MyWindow = typeof window & {
  api: (channel: string, ...args: unknown[]) => unknown,
};

export default function Main(): JSX.Element {
  const [num, setNum] = useState(0);

  useEffect(() => {
    (async () => {
      const n = await (window as MyWindow).api.call("data");
      setNum(n);
    })();
  })

  return <h2>Hello from React! {num}</h2>;
}
