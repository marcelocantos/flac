#!/usr/bin/env python3

import collections
import re
import sqlite3

def load_data():
    forward = {
        (pinyin, 1 + i): tone
        for line in open('flac.data').readlines()
        if line.rstrip() and 'ā' not in line
        for cols in [line.rstrip().split('\t')]
        for (pinyin, tones) in [(cols[0], cols[1:])]
        for (i, tone) in enumerate(tones)
    }
    # reverse = collections.defaultdict(set)
    # for pinyin, tones in forward.items():
    #     for 
    return forward

class Prompter:
    def ask(self, word):
        self.word = word
        self.text = input(self.word + ' — ')
        return self.text

    def check(self, ok, fmt=None, *args):
        if fmt is None:
            fmt = ''
        indent = 2*len(self.word) + 3 + len(self.text)
        outcome = '✅' if ok else '❌'
        print('\033[A\033[%dC  %s \033[1;31m%s\033[0m' % (indent, outcome, '' if ok else fmt % args))
        return ok

def main():
    data = load_data()
    print(data)
    map = collections.defaultdict()

    prompter = Prompter()
    wsRE = re.compile(r'\s+')
    pinyinRE = re.compile(r'[a-zéü]+\d')
    while True:
        phrase = '你好'
        text = prompter.ask(phrase)
        pinyins = pinyinRE.findall(text)
        if prompter.check(
            sum(len(w) for w in pinyins) == len(wsRE.sub('', text)),
            'unrecognised elements: %s', re.sub('\s+', ' ', pinyinRE.sub(' ', text).strip()),
        ) and prompter.check(
            len(phrase) == len(pinyins),
            'character count mismatch: %s has %d characters; %s has %d pinyins words' % (phrase, len(phrase), text, len(pinyins)),
        ) and prompter.check(
            all(
                c in data.get((word, tone), '')
                for (c, pinyin) in zip(phrase, pinyins)
                for (word, tone) in [(pinyin[:-1], int(pinyin[-1]))]
            ),
            'mismatched pinyin, correct form is: '
        ):
            break

if __name__ == '__main__':
    main()
