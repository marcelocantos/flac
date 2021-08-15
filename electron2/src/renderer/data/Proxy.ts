import * as Interface from '../../common/data/Interface';

type Api = {
  call: (channel: string, ...args: unknown[]) => Promise<unknown>,
};

type MyWindow = typeof window & {
  api: Api,
};

class Bridge {
  api: Api;

  constructor() {
    this.api = (window as MyWindow).api;
  }

  async call<T>(method: string, ...params: unknown[]): Promise<T> {
    params = ([method as unknown]).concat(params);
    return await this.api.call("call", ...params) as T;
  }

  async get<T>(method: string): Promise<T> {
    return await this.api.call("get", method) as T;
  }
}

export class Database implements Interface.Database {
  bridge: Bridge;

  constructor() {
    this.bridge = new Bridge();
  }

  close(): Promise<void> {
    return this.bridge.call<void>("close");
  }

  get HeadWord(): Promise<{word: string, score: number}> {
    return this.bridge.get<{word: string, score: number}>("HeadWord");
  }

  get MaxScore(): Promise<number> {
    return this.bridge.get<number>("MaxScore");
  }

  get MaxPos(): Promise<number> {
    return this.bridge.get<number>("MaxPos");
  }

  MoveWord(word: string, dest: number): Promise<void> {
    return this.bridge.call<void>("MoveWord", word, dest);
  }

  UpdateScore(word: string, score: number): Promise<void> {
    return this.bridge.call<void>("UpdateScore", word, score);
  }

  UpdateScoreAndPos(word: string, score: number, dest: number): Promise<void> {
    return this.bridge.call<void>("UpdateScoreAndPos", word, score, dest);
  }

  SetFocus(focus: string): Promise<void> {
    return this.bridge.call<void>("SetFocus", focus);
  }

  WordScore(word: string): Promise<number> {
    return this.bridge.call<number>("WordScore", word);
  }

  WordPos(word: string): Promise<number> {
    return this.bridge.call<number>("WordPos", word);
  }

  WordAt(pos: number): Promise<string> {
    return this.bridge.call<string>("WordAt", pos);
  }
}
