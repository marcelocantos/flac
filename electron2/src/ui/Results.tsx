import React from 'react';

interface Props {
    log: JSX.Element[];
    streak: string[];
}

export default function 结果清单({log, streak}: Props): JSX.Element {
    return (
        <div>
            {log.map(e => <div>{e}</div>)}
            <div style={{display: "flex", flexWrap: "wrap"}}>
                {streak.map(s => <div>{s}</div>)}
            </div>
        </div>
    )
}
