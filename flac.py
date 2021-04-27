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
    reverse = collections.defaultdict(lambda: collections.defaultdict(set))
    for (pinyin, tone), words in forward.items():
        for word in words:
            reverse[word][pinyin].add(tone)
    return (forward, reverse)

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
    (forward, reverse) = load_data()
    print(forward)
    print(reverse)
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
                c in forward.get((word, tone), '')
                for (c, pinyin) in zip(phrase, pinyins)
                for (word, tone) in [(pinyin[:-1], int(pinyin[-1]))]
            ),
            'mismatched pinyin, correct form is: ' + ' '.join([''.join('%s%s' % (pinyin, ''.join(str(t) for t in tones)) for (pinyin, tones) in reverse[c].items()) for c in phrase])
        ):
            break

if __name__ == '__main__':
    main()
