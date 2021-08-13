const sqlite3 = require('sqlite3');

function promise(obj, method, descr, ...params) {
  const log = true;
  // const log = params.length > 0 && params[0].LOG;
  // if (log) {
  //   delete params[0].LOG;
  // }
  return new Promise((resolve, reject) => {
    obj[method](...params, (error, result) => {
      if (error) {
        !log || console.log({obj, method: obj[method], descr, params, error});
        reject(error);
      } else {
        !log || console.log({obj, method: obj[method], descr, params, result});
        resolve(result);
      }
    })
  });
}

class Statement {
  stmt;
  sql;

  constructor(stmt, sql) {
    this.stmt = stmt;
    this.sql = sql;
  }

  get all() {
    return (params) => {
      return promise(this.stmt, 'all', this.sql, params || {});
    }
  }

  finalize() {
    return new Promise((resolve, reject) => {
      this.stmt.finalize(error => {
        if (error) {
          reject(error);
        } else {
          resolve();
        }
      });
    })
  }

  get get() {
    return async (params) => {
      // Use all() coz get() can leave db locked. https://bit.ly/3lHc0Ic
      return (await this.all(params ?? {}))[0];
    }
  }

  get run() {
    return (params) => {
      return promise(this.stmt, 'run', this.sql, params ?? {});
    }
  }
}

class Database {
  stmts = [];

  constructor(db) { this.db = db; }

  async all(sql, params) {
    const stmt = await this.prepare(sql);
    return await stmt.all(params);
  }

  async close() {
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

  exec(sql) {
    return promise(this.db, 'exec', 'Database', sql);
  }

  async finalize() {
    for (const stmt of this.stmts) {
      await stmt.finalize();
    }
    this.stmts = [];
  }

  async get(sql, params) {
    const stmt = await this.prepare(sql);
    return await stmt.get(params);
  }

  async prepare(sql) {
    return new Promise((resolve, reject) => {
      const stmt = this.db.prepare(sql, (error) => {
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

  async run(sql, params) {
    const stmt = await this.prepare(sql);
    return await stmt.run(params);
  }

  async tx(cb) {
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

function open(filename) {
  return new Promise((resolve, reject) => {
    let db = new sqlite3.Database(filename, (error) => {
      if (error) {
        reject(error)
      } else {
        resolve(new Database(db));
      }
    })
  })
};

exports.open = open;
