import React from 'react';

import Word from '../pinyin/Word';
import Pinyin from '../pinyin/Pinyin';

import './ui.css';

interface 汉和拼音字特性 {
  汉: string;
  拼音: Word | string | Pinyin[];
  [attrs: string]: unknown;
}

export default function 汉和拼音字({汉, 拼音}: 汉和拼音字特性): JSX.Element {
  return (
    汉
    ? <div className="汉和拼音字">
        {<div className="汉">{汉}</div>}
        <div className="拼音"><Word.HTML word={拼音}/></div>
      </div>
    : <Word.HTML word={拼音}/>
  );
}
