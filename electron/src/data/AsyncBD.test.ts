import * as AsyncDB from './AsyncDB';

it('close', async () => {
  const db = await AsyncDB.open(':memory:');
  await db.close();
});

it('close twice throws', async () => {
  const db = await AsyncDB.open(':memory:');
  await db.close();
  await expect(db.close()).rejects.toThrow();
});

it('throws with bad filename', async () => {
  await expect(AsyncDB.open('./nowhere/nothing')).rejects.toThrow();
});

it('throws with bad query', async () => {
  const db = await AsyncDB.open(':memory:');
  await expect(db.run('INSERT INTO foo (x) VALUES (42)')).rejects.toThrow();
});

it('throws calling stmt.finalize twice', async () => {
  const db = await AsyncDB.open(':memory:');
  await db.run('CREATE TABLE foo (x int);');
  const stmt = await db.prepare('INSERT INTO foo (x) VALUES (42)');
  await stmt.finalize();
  await expect(stmt.finalize()).rejects.toThrow();
});

it('get throws with bad params', async () => {
  const db = await AsyncDB.open(':memory:');
  await db.run('CREATE TABLE foo (x int);');
  await expect(db.get('SELECT * FROM foo WHERE x = $x', {$y: 42})).rejects.toThrow();
});

it('all throws with bad params', async () => {
  const db = await AsyncDB.open(':memory:');
  await db.run('CREATE TABLE foo (x int);');
  await expect(db.all('SELECT * FROM foo WHERE x = $x', {$y: 42})).rejects.toThrow();
});

it('all with no params', async () => {
  const db = await AsyncDB.open(':memory:');
  await db.run('CREATE TABLE foo (x int);');
  await db.run('INSERT INTO foo (x) VALUES (42);');
  expect(await db.all('SELECT * FROM foo')).toEqual([{x: 42}]);
});
