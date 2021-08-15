import React from 'react';

interface Props {
  score: number;
}

export default function Score({score}: Props): JSX.Element {
  const logScore = Math.log(1 + score) / Math.log(1000);

  return (
    <svg viewBox="0 0 8 21" xmlns="http://www.w3.org/2000/svg"
        style={{
          height: "21px",
          verticalAlign: "text-bottom",
        }}
      >
      <clipPath id="clipper" clipPathUnits="objectBoundingBox">
          <rect y={1 - logScore} width="1" height={logScore}/>
      </clipPath>

      <rect width="100" height="100" fill="white"/>

      <path id="progress"
          d="M4,1 h1 a2,2 0 0 1 2,2 v15 a2,2 0 0 1 -2,2 h-1 a2,2 0 0 1 -2,-2 v-15 a2,2 0 0 1 2,-2 z"
          fill="rgb(0, 204, 34)"
          clipPath="url(#clipper)"
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
