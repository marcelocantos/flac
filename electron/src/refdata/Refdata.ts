import refdata from './refdata.json';

export interface Definitions {
  definitions: string[];
}

export interface Entries {
  entries: {[key: string]: Definitions};
  traditional: string;
}

export interface Refdata {
  dict: {
    ambiguousWords: {[key: string]: boolean},
    entries: {[key: string]: Entries},
    pinyinToSimplified: {[key: string]: {
      words: string[],
    }},
    traditionalToSimplified: {[key: string]: string},
    validSyllables: {[key: string]: boolean},
  };
  wordList: {
    words: string[],
    frequencies: {[key: string]: number},
  };
}

export default refdata as Refdata;
