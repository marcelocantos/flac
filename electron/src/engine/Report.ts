import Proxy from '../renderer/data/Proxy';
import { Refdata } from '../refdata/Refdata';
import Outcome from '../outcome/Outcome';

const 记录 = false;

export default class 汇报类 {
  refreshCount: number;

  历史: Outcome[][] = [];
  好组: Outcome[] = [];

  onScoreChangedFunc: (word: string, score: number) => void;

  constructor(
    private db: Proxy,
    private rd: Refdata,
  ){}

  async 好(字: string, 产物: Outcome, 容易: boolean): Promise<void> {
    产物.Score = await this.bump(字, score => {
      // 产物.html.分数(score);
      return {score: Math.max(2, 2 * (score ?? 0)), move: true};
    });

    this.好组.push(产物);
  }

  async 不好(产物: Outcome, 容易: boolean, 尝试包装器: {尝试: number}): Promise<void> {
    if (产物.不及格) {
      try {
        const penalty = Math.sqrt(1 + 尝试包装器.尝试);
        尝试包装器.尝试++;

        // Multiply score by 1/2√(1 + attempt).
        产物.Score = await this.bump(产物.Word, score => ({
          score: Math.max(1, (score ?? 0) / (2 * penalty)),
          move: false,
        }));
      } finally {
        if (this.好组.length > 0) {
          this.历史.push(this.好组);
          this.好组 = [];
        }
        this.历史.push([产物]);
      }
    }
  }

  async 放弃(结果: Outcome): Promise<void> {
    await this.bump(结果.Word, score => {
      return {score: Math.max(1, score / 8), move: false};
    })
  }

  private async bump(word: string, bump: (score: number) => {score: number, move: boolean}): Promise<number> {
    const {score, move} = bump(await this.db.WordScore(word));
    let pos = -1;
    if (move) {
      pos = score + Math.floor(Math.random() * (1+score*3/2-score));
    }

    if (记录) console.log('setScoreAndPos', {word, score, pos});
    await this.db.UpdateScoreAndPos(word, score, pos);
    return score;
  }
}
