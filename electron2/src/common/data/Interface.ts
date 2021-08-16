export interface Database {
  close(): Promise<void>
  get HeadWord(): Promise<{word: string, score?: number} | undefined>
  get MaxScore(): Promise<number | undefined>
  get MaxPos(): Promise<number | undefined>
  MoveWord(word: string, dest: number): Promise<void>
  UpdateScore(word: string, score: number): Promise<void>
  UpdateScoreAndPos(word: string, score: number, dest: number): Promise<void>
  SetFocus(focus: string): Promise<void>
  WordScore(word: string): Promise<number | undefined>
  WordPos(word: string): Promise<number | undefined>
  WordAt(pos: number): Promise<string | undefined>
}
