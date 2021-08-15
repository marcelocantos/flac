export interface Database {
  close(): Promise<void>
  get HeadWord(): Promise<{word: string, score: number}>
  get MaxScore(): Promise<number>
  get MaxPos(): Promise<number>
  MoveWord(word: string, dest: number): Promise<void>
  UpdateScore(word: string, score: number): Promise<void>
  UpdateScoreAndPos(word: string, score: number, dest: number): Promise<void>
  SetFocus(focus: string): Promise<void>
  WordScore(word: string): Promise<number>
  WordPos(word: string): Promise<number>
  WordAt(pos: number): Promise<string>
}
