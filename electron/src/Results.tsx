function Results({log, streak}: {
    log: any[],
    streak: string[]
}) {
    return (
        <div>
            {log.map((e: any) => <div>{e}</div>)}
            <div style={{display: "flex", flexWrap: "wrap"}}>
                {streak.map((s: string) =>
                    <div>{s}</div>
                )}
            </div>
        </div>
    )
}

export default Results;
