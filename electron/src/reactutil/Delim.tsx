import React from 'react';

export default function Delim({list, delim}: {list: React.ReactNode[], delim: string}): JSX.Element {
  return <>{list.map((e, i) => <>{i ? delim : ''}{e}</>)}</>
}
