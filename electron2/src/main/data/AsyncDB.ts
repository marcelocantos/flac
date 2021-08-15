import * as sqlite3 from 'sqlite3';

export type Params = {[_: string]: unknown};
export type Row = {[_: string]: unknown};
export type Result = {[_: string]: unknown};

export type All = (params?: Params) => Promise<Row[]>;
export type Get = (params?: Params) => Promise<Row | undefined>;
export type Run = (params?: Params) => Promise<Result>;

function promise<T>(
  obj: unknown,
  method: string,
  descr: string,
  ...params: (Params | string)[]
): Promise<T> {
  let log = false;
  if (params.length > 0 && typeof params[0] !== "string") {
    log = params[0] && 'LOG' in params[0];
    if (log) {
      delete params[0].LOG;
    }
  }
  return new Promise((resolve, reject) => {
    (obj as any)[method](...params, (error: unknown, result: unknown) => {
      if (error) {
        !log || console.log({obj, method, descr, params, error});
        reject(error);
      } else {
        !log || console.log({obj, method, descr, params, result});
        resolve(result as T);
      }
    })
  });
}

export class Statement {
  constructor(
    public stmt: sqlite3.Statement,
    public sql: string,
  ){}

  get all(): All {
    return (params?: Params): Promise<Row[]> => {
      return promise<Row[]>(this.stmt, 'all', this.sql, params ?? {});
    };
  }

  finalize(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.stmt.finalize((error: unknown) => {
        if (error) {
          reject(error);
        } else {
          resolve();
        }
      });
    })
  }

  get get(): Get {
    return async (params?: Params): Promise<Row | undefined> => {
      return (await this.all(params ?? {}))[0];
    };
  }

  get run(): Run {
    return async (params?: Params): Promise<Result> => {
      return promise(this.stmt, 'run', this.sql, params);
    };
  }
}

export class Database {
  stmts: Statement[] = [];

  constructor(
    public db: sqlite3.Database,
  ){}

  async all(sql: string, params?: Params): Promise<Row[]> {
    const stmt = await this.prepare(sql);
    return await stmt.all(params);
  }

  async close(): Promise<void> {
    await this.finalize();
    return new Promise((resolve, reject) =>
      this.db.close((error) => {
        if (error) {
          reject(error);
        } else {
          resolve();
        }
      })
    );
  }

  exec(sql: string): Promise<unknown> {
    return promise(this.db, 'exec', 'Database', sql);
  }

  async finalize(): Promise<void> {
    for (const stmt of this.stmts) {
      await stmt.finalize();
    }
    this.stmts = [];
  }

  async get(sql: string, params?: Params): Promise<unknown> {
    const stmt = await this.prepare(sql);
    return await stmt.get(params);
  }

  async prepare(sql: string): Promise<Statement> {
    return new Promise((resolve, reject) => {
      const stmt = this.db.prepare(sql, (error: unknown) => {
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

  async run(sql: string, params?: Params): Promise<unknown> {
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
  return new Promise((resolve, reject) => {
    const db = new sqlite3.Database(filename, (error: unknown) => {
      if (error) {
        reject(error)
      } else {
        resolve(new Database(db));
      }
    })
  })
}
