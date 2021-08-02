import refdata from './refdata.json';

export type Refdata = {
  dict: {
    ambiguousWords: {[key: string]: boolean},
    entries: {[key: string]: {
      entries: {[key: string]: {
        definitions: string[],
      }},
      traditional: string,
    }},
    pinyinToSimplified: {[key: string]: {
      words: string[],
    }},
    traditionalToSimplified: {[key: string]: string},
    validSyllables: {[key: string]: boolean},
  }
  wordList: {
    words: string[],
    frequencies: {[key: string]: number},
  },
};

export default refdata as Refdata;
