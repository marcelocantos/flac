#!/usr/bin/env python3

import pickle
import random
import re
import sqlite3

class default:
    class _naught:
        pass

    def __init__(self, d, ctor):
        self._d = d
        self._ctor = ctor

    def __getitem__(self, key):
        v = self._d.get(key, default._naught())
        if isinstance(v, default._naught):
            self._d[key] = v = self._ctor()
        return v

def load_data():
    forward = {
        (pinyin, 1 + i): tone
        for line in open('flac.data').readlines()
        if line.rstrip() and 'ā' not in line
        for cols in [line.rstrip().split('\t')]
        for (pinyin, tones) in [(cols[0], cols[1:])]
        for (i, tone) in enumerate(tones)
        if tone
    }
    reverse = {}
    for (pinyin, tone), words in forward.items():
        for word in words:
            default(default(reverse, lambda: {})[word], set)[pinyin].add(tone)
    return (forward, reverse)

class Prompter:
    def ask(self, word):
        self.word = word
        self.text = input('\n\033[A%s — \033[K' % (self.word,))
        return self.text

    def check(self, final, ok, fmt=None, *args):
        if fmt is None:
            fmt = ''
        outcome = (
            '\033[%dC❌' % (max(8 - len(self.text), 0),) if not ok else
            '\033[%dD✅\033[K' % (len(self.text),) if final else
            '')
        format = '%s %s' if '\v' in fmt else '%s\v %s'
        self.message(format % (outcome, '' if ok else fmt % args), wait=not ok)
        return ok

    def message(self, m, wait):
        indent = 2*len(self.word) + 3 + len(self.text)
        print('\033[A\033[%dC%s' % (indent, m.replace('\v', '')), end='')
        if wait:
            input()
            print('\033[A\033[%dC\033[K%s' % (indent, m.split('\v')[0]))
        else:
            print()

pinyinColors = {
    1: 31,
    2: 32,
    3: 34,
    4: 35,
    5: 30,
}

def pinyinColor(pinyin, tone):
    return '\033[1;%dm%s\033[0m' % (pinyinColors[tone], pinyin)

vowels = {
    'a': ' āáǎàa',
    'e': ' ēéěèe',
    'i': ' īíǐìi',
    'o': ' ōóǒòo',
    'u': ' ūúǔùu',
    'ü': ' ǖǘǚǜü',
}

def accent(pinyin, tone):
    chars = list(pinyin)
    v = sum(c in vowels for c in chars)
    # https://en.wikipedia.org/wiki/Pinyin#Rules_for_placing_the_tone_mark
    if v == 1:
        for i, c in enumerate(chars):
            if c in vowels:
                chars[i] = vowels[c][tone]
                break
    elif 'a' in pinyin or 'e' in pinyin:
        for i, c in enumerate(chars):
            if c in 'ae':
                chars[i] = vowels[c][tone]
                break
    elif 'ou' in pinyin:
        for i, c in enumerate(chars):
            if c == 'o':
                chars[i] = vowels[c][tone]
                break
    else:
        for i in range(len(chars) - 1, -1, -1):
            c = chars[i]
            if c in vowels:
                chars[i] = vowels[c][tone]
                break
    return pinyinColor(''.join(chars), tone)

def pinyinTones(pinyin, tones):
    return '/'.join(accent(pinyin, t) for t in sorted(tones))

# class scoredb:
#     def __init__(self):
#         self.db = sqlite3.connect('flac.db')
#         self.db.execute('create table if not exists word_has_score (word primary key, score)')

#     def getscore(self, word):
#         with self.db.cursor() as cur:
#             cur.execute('select score from word_has_score where word = ?', (word,))
#             results = cur.fetchall()
#             return results[0][0] if results else 0

#     def setscore(self, word, score):
#         with self.db.cursor() as cur:
#             cur.execute('replace into word_has_score (word, score) values (?, ?)', (word, score))

def clamp(lo, hi):
    return lambda x: min(max(lo, x), hi)

# SRSQueue manages test words in a priority queue, bumping words up and down the
# queue depending the result of each test.
class SRSQueue:
    def __init__(self, words):
        try:
            with open('queue.pickle', 'rb') as f:
                data = pickle.load(f)
                self.queue = [w for w in data['queue'] if w in words]
                self.scores = data['scores']
        except FileNotFoundError:
            self.queue = []
            self.scores = {}
        new = list(words - set(self.queue))
        random.shuffle(new)
        self.queue += new
        last = len(self.queue) - 1
        self.clamp = clamp(-last, last)

    def __enter__(self):
        return self

    def __exit__(self, *ex):
        with open('queue.pickle', 'wb') as f:
            pickle.dump({
                'queue': self.queue,
                'scores': self.scores,
            }, f)

    def head(self):
        return self.queue[0]

    # good > 0, skip < 0
    def score(self):
        return self.scores.get(self.head(), 0)

    # Bump back a little with slow exponential growth.
    def good(self):
        self.bump(max(self.score(), 2)*3//2)

    # Bump back a lot with rapid exponential growth.
    def skip(self):
        self.bump(min(self.score(), -50)*4)

    # Bump forward all the way to head with rapid exponential approach.
    def bad(self):
        self.bump(max(self.score()//8, 1))

    def bump(self, score):
        self.scores[self.head()] = self.clamp(score)
        i = random.randint(abs(score), self.clamp(abs(score)*3//2))
        if i:
            w = self.queue.pop(0)
            self.queue.insert(i, w)

def main():
    wsRE = re.compile(r'[\s/]+')
    pinyinRE = re.compile(r'([a-zü]+)(\d+)')

    (forward, reverse) = load_data()
    # print(forward)
    # print(reverse)
    # db = scoredb()

    def check(c, pinyins):
        rev = {}
        for p in pinyins:
            (word, tones) = pinyinRE.match(p).groups()
            default(rev, set)[word].update(int(t) for t in tones)
        # print(rev)
        # print(reverse[c])
        return rev == reverse[c]
                    # c in forward.get((word, tone), '')

    def correction(phrase):
        return ' '.join(
            '/'.join(pinyinTones(*i) for i in reverse[c].items())
            for c in phrase
        )

    def lookup(c, pinyins):
        lookups = ', '.join(
            '%s = %s' % (accent(pinyin, tone), forward.get((pinyin, tone), '∅'))
            for p in pinyins
            for (pinyin, tones) in pinyinRE.findall(p)
            for tone in sorted(tones)
            for tone in [int(tone)]
            if c not in forward.get((pinyin, tone), '')
        )
        return lookups and '(' + lookups + ')'

    prompter = Prompter()
    with SRSQueue(set(reverse)) as q:
        while True:
            char = q.head()
            color = '0' if q.scores.get(char) else '1;37'
            p = '\033[%sm%s\033[0m' % (color, char)
            while True:
                # Print new characters in bright white.
                text = prompter.ask(p)
                if text:
                    break
                prompter.check(False, False, '\v' + correction(char))
                q.bad()
                continue

            pinyins = [''.join(p).replace('v', 'ü') for p in pinyinRE.findall(text)]
            if not prompter.check(
                False,
                sum(len(w) for w in pinyins) == len(wsRE.sub('', text)),
                '\033[1;31munrecognised elements:\033[0m %s', re.sub('\s+', ' ', pinyinRE.sub(' ', text).strip()),
            ):
                print('\033[A', end='')
            elif prompter.check(
                True,
                check(char, pinyins),
                '%s\v %s', lookup(char, pinyins), correction(char)
            ):
                q.good()
            else:
                q.bad()

if __name__ == '__main__':
    main()
