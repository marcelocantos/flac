import React from 'react';

interface Props {
  分数: number;
}

export default function 分数条({分数}: Props): JSX.Element {
  const 对数分数 = Math.log(1 + 分数) / Math.log(1000);

  return (
    <svg viewBox="0 0 8 21" xmlns="http://www.w3.org/2000/svg"
        style={{
          height: "21px",
          verticalAlign: "text-bottom",
        }}
      >
      <clipPath id="剪子" clipPathUnits="objectBoundingBox">
          <rect y={1 - 对数分数} width="1" height={对数分数}/>
      </clipPath>

      <rect width="100" height="100" fill="white"/>

      <path id="progress"
          d="M4,1 h1 a2,2 0 0 1 2,2 v15 a2,2 0 0 1 -2,2 h-1 a2,2 0 0 1 -2,-2 v-15 a2,2 0 0 1 2,-2 z"
          fill="rgb(0, 204, 34)"
          clipPath="url(#剪子)"
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
