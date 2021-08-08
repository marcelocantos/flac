import sqlite3 from 'sqlite3';

export type Params = {[param: string]: any};
export type Row = {[param: string]: any};
export type Aller = (params?: Params) => Promise<Row[]>;
export type Getter = (params?: Params) => Promise<Row | undefined>;
export type Runner = (params?: Params) => Promise<void>;

function promise<T>(obj: any, method: string, descr: string, ...params: any): Promise<T> {
  // const log = params.length > 0 && params[0].LOG;
  // if (log) {
  //   delete params[0].LOG;
  // }
  return new Promise<T>((resolve, reject) => {
    obj[method](...params, (error?: any, result?: any) => {
      if (error) {
        // !log || console.log({obj, method: obj[method], descr, params, error});
        reject(error);
      } else {
        // !log || console.log({obj, method: obj[method], descr, params, result});
        resolve(result);
      }
    })
  });
}

export class Statement {
  constructor(
    private readonly stmt: any,
    public readonly sql: string,
  ) {}

  get all(): Aller {
    return (params?: Params): Promise<Row[]> => {
      return promise<Row[]>(this.stmt, 'all', this.sql, params ?? {});
    }
  }

  finalize(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      this.stmt.finalize((error?: any) => {
        if (error) {
          reject(error);
        } else {
          resolve();
        }
      });
    })
  }

  get get(): Getter {
    return async (params?: Params): Promise<Row | undefined> => {
      // Use all() coz get() can leave db locked. https://bit.ly/3lHc0Ic
      return (await this.all(params ?? {}))[0];
    }
  }

  get run(): Runner {
    return (params?: Params): Promise<void> => {
      return promise<void>(this.stmt, 'run', this.sql, params ?? {});
    }
  }
}

export class Database {
  private stmts: Statement[] = [];

  constructor(
    private readonly db: any,
  ) {}

  async all(sql: string, params?: Params): Promise<Row[]> {
    const stmt = await this.prepare(sql);
    return await stmt.all(params);
  }

  async close(): Promise<void> {
    await this.finalize();
    return new Promise<void>((resolve, reject) =>
      this.db.close((error?: any) => {
        if (error) {
          reject(error);
        } else {
          resolve();
        }
      })
    );
  }

  exec(sql: string): Promise<void> {
    return promise<void>(this.db, 'exec', 'Database', sql);
  }

  async finalize(): Promise<void> {
    for (const stmt of this.stmts) {
      await stmt.finalize();
    }
    this.stmts = [];
  }

  async get(sql: string, params?: Params): Promise<Row | undefined> {
    const stmt = await this.prepare(sql);
    return await stmt.get(params);
  }

  async prepare(sql: string): Promise<Statement> {
    return new Promise<Statement>((resolve, reject) => {
      const stmt = this.db.prepare(sql, (error?: any) => {
        if (error) {
          reject(error);
        } else {
          const result = new Statement(stmt, sql);
          this.stmts.push(result);
          resolve(result);
        }
      });
    })
  }

  async run(sql: string, params?: Params): Promise<void> {
    const stmt = await this.prepare(sql);
    return await stmt.run(params);
  }

  async tx<T>(cb: () => Promise<T>): Promise<T> {
    await this.run("BEGIN");
    try {
      const ret = await cb();
      await this.run("COMMIT");
      return ret;
    } catch (error) {
      await this.run("ROLLBACK");
      throw error;
    }
  }
}

export function open(filename: string): Promise<Database> {
  return new Promise<Database>((resolve, reject) => {
    let db = new sqlite3.Database(filename, (error?: any) => {
      if (error) {
        reject(error)
      } else {
        resolve(new Database(db));
      }
    })
  })
}
