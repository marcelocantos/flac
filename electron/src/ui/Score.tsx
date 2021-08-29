import React from 'react';

let svgID = 0;

const logBase = Math.log(4000);

type 分数条特性 = {
  分数: number;
  比例?: number;
}

export default function 分数条({分数, 比例}: 分数条特性): JSX.Element {
  const 对数分数 = Math.log(1 + 分数) / logBase;

  const id = `剪子${svgID++}`;

  return (
    <svg viewBox="0 0 8 21" xmlns="http://www.w3.org/2000/svg"
        style={{
          height: "1.1em",
          verticalAlign: "text-bottom",
        }}
      >
      <clipPath id={id} clipPathUnits="objectBoundingBox">
          <rect y={1 - 对数分数} width="1" height={对数分数}/>
      </clipPath>

      <rect width="100" height="100" fill="white"/>

      <path id="progress"
          d="M4,1 h1 a2,2 0 0 1 2,2 v15 a2,2 0 0 1 -2,2 h-1 a2,2 0 0 1 -2,-2 v-15 a2,2 0 0 1 2,-2 z"
          fill="rgb(0, 204, 34)"
          clipPath={`url(#${id})`}
      />
      <path
          d="M4,1 h1 a2,2 0 0 1 2,2 v15 a2,2 0 0 1 -2,2 h-1 a2,2 0 0 1 -2,-2 v-15 a2,2 0 0 1 2,-2 z"
          fill="#0000"
          stroke="gray"
          strokeWidth="1"
      />
    </svg>
  )
}
