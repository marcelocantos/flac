import Word from './Word';
import Pinyin from './Pinyin';

it('new Word', () => {
  expect(new Word("wo3 men5")).toEqual(
    new Word([new Pinyin("wo3 "), new Pinyin("men5")])
  );
});
