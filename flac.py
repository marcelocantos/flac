#!/usr/bin/env python3

import collections
import contextlib
import math
import pickle
import random
import re
import shutil
import sys

def termwidth():
    return shutil.get_terminal_size().columns

def dots(score):
    return math.log(max(1, score), 1.5)

def scorerepr(score):
    if score <= 0:
        return ''
    logscore = int(dots(score))
    s = ''
    while logscore > 7:
        s += '⣿'
        # s += '█'
        logscore -= 8
    
    return (s + ['','⡀','⡄','⡆','⡇','⣇','⣧','⣷'][logscore])
    # return s + ' ▏▎▍▌▋▊▉'[logscore]

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
    entries = {
        pinyin: tones
        for line in open('flac.data').readlines()
        if line.rstrip() and 'ā' not in line
        for cols in [line.rstrip().split('\t')]
        for (pinyin, tones) in [(cols[0], cols[1:])]
    }

    syllables = set(entries.keys())

    forward = {
        (pinyin, 1 + i): tone
        for (pinyin, tones) in entries.items()
        for (i, tone) in enumerate(tones)
        if tone
    }

    reverse = {}
    for (pinyin, tone), words in forward.items():
        for word in words:
            default(default(reverse, lambda: {})[word], set)[pinyin].add(tone)

    return (syllables, forward, reverse)

class Prompter:
    def ask(self, word, score, togo):
        new = score == 0
        color = '1;37' if new else '0'
        self.word = word
        # print('' % (), end='')
        self.text = input(
            '\n\033[A\033[K\033[%dG%d\033[G\033[%sm%s\033[1;30m%-2s\033[0m — '
            % (termwidth() - 4, togo, color, word, scorerepr(score)))
        return self.text, color

    def checknowait(self, final, ok, fmt=None, *args):
        return self._check(final, ok, True, fmt, *args)

    def check(self, final, ok, fmt=None, *args):
        return self._check(final, ok, False, fmt, *args)

    def _check(self, final, ok, neverwait, fmt=None, *args):
        if fmt is None:
            fmt = ''
        outcome = (
            '\033[%dC❌' % (max(10 - len(self.text), 0),) if not ok else
            '\033[%dD✅\033[K' % (len(self.text),) if final else
            '')
        format = '%s %s' if '\v' in fmt else '%s\v %s'
        self.message(format % (outcome, '' if ok else fmt % args), wait=not ok and not neverwait)
        return ok

    def message(self, m, wait):
        indent = 2*len(self.word) + 5 + len(self.text)
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

# autopickle loads a pickle (with a default fallback) and saves it out later.
@contextlib.contextmanager
def autopickle(filename, default):
    try:
        with open(filename, 'rb') as f:
            data = pickle.load(f)
    except FileNotFoundError:
        data = default

    yield data

    with open(filename, 'wb') as f:
        pickle.dump(data, f)

# queuedata loads queue data from the pickle
@contextlib.contextmanager
def queuedata(words):
    with autopickle('queue.pickle', {'queue': [], 'scores': {}}) as data:
        queue = [w for w in data['queue'] if w in words]
        new = list(words - set(queue))
        if new:
            random.shuffle(new)
            queue += new
            data['queue'] = queue
        yield data


# SRSQueue manages test words in a priority queue, bumping words up and down the
# queue depending the result of each test.
class SRSQueue:
    def __init__(self, data):
        self.data = data
        self.queue = data['queue']
        self.scores = data['scores']
        last = len(self.queue) - 1
        self.clamp = clamp(-last, last)
        self.goods = [0, 0]
        self.bads = [0, 0]
        self.skips = 0
        self.chars = set()

    def __enter__(self):
        return self

    def __exit__(self, *ex):
        self.data['queue'] = [q for q in self.queue if q]
        self.data['scores'] = self.scores

    def report(self):
        return ('''
            ┌─────┤scores├─────┐
              right :  %3d %3d
              wrong :  %3d \033[%sm%3d\033[0m
              huh?  :  %3d
             ──────────────────
              chars : %4d
              Σchars: %4d
            └──────────────────┘
            %s
        '''.strip().replace('\n            ', '\n') % (
            self.goods[False], self.goods[True],
            self.bads[False], '1;31' if self.bads[True] else '', self.bads[True],
            self.skips,
            len(self.chars),
            len(self.scores),
            self.histogram(),
        ))

    def histogram(self):
        hist = collections.defaultdict(int)
        maxs = 0
        for s in self.scores.values():
            s = int(dots(s))
            maxs = max(maxs, s)
            hist[s] += 1
        bars = [
            '█' * (n//8) + ('','▁','▂','▃','▄','▅','▆','▇')[n%8]
            for s in range(maxs + 1)
            for n in [int(4*math.log(hist.get(s) + 1, 2) if s in hist else 0)]
        ]
        maxbars = max(len(b) for b in bars)
        bars = [b.ljust(maxbars) for b in bars]
        return '\n'.join(reversed([''.join(c) for c in zip(*bars)]))

    def head(self):
        return self.queue[0]

    # good > 0, skip < 0
    def score(self, char):
        return self.scores.get(char, 0)

    # Bump back a little with slow exponential growth.
    def good(self, char, easy):
        self.bump(char, max(self.score(char), 2)*3//2)
        self.goods[easy] += 1

    # Bump forward all the way to head with slow exponential approach.
    def bad(self, char, easy, attempt):
        # Multiply score by 2/(3*sqrt(attempt)).
        penalty = int(10 * attempt**0.5)
        self.bump(char, max(self.score(char)*20//3//penalty, 1))
        self.bads[easy] += 1

    # Bump forward all the way to head with rapid exponential approach.
    def skip(self, char):
        self.bump(char, max(self.score(char)//8, 1))
        self.skips += 1

    def bump(self, char, score, maxpos=math.inf):
        score = self.clamp(score)
        self.chars.add(char)
        self.scores[char] = score
        p = min(abs(score), maxpos)
        i = random.randint(p, self.clamp(p*3//2))
        if i:
            self.queue.remove(char)
            self.queue.insert(i, char)

def focusqueue(queue, chars):
    queuechars = set(queue)
    focuschars = chars & queuechars
    nonchars = chars - queuechars
    print('non-chars:', ''.join(sorted(nonchars)))
    blurchars = [c for c in queue if c not in focuschars]
    focuschars = [c for c in queue if c in focuschars]
    print('focus (%d): %s' % (len(focuschars), ''.join(focuschars)))
    return focuschars + [None] + blurchars, focuschars

def focusreport(fscores):
    fscores = sorted(fscores)
    n = len(fscores)
    percentiles = [0]
    if n >= 5:
        percentiles.append(n//4)
    if n >= 3:
        percentiles.append(n//2)
    if n >= 5:
        percentiles.append(n*3//4)
    if n >= 2:
        percentiles.append(n - 1)
    return '  '.join('%4.1f' % (dots(fscores[p]),) for p in percentiles)

def main():
    wsRE = re.compile(r'[\s/]+')
    pinyinRE = re.compile(r'([a-zü]+)(\d+)')

    (syllables, forward, reverse) = load_data()

    fscores = None
    with queuedata(set(reverse)) as data:
        scores = data['scores']

        focus = sys.argv[1:2] == ['--focus']
        if focus:
            data['queue'], focuschars = focusqueue(data['queue'], set(sys.argv[2]))
        else:
            focuschars = data['queue']
        fscores = [scores.get(c, 0) for c in focuschars]

        def check(c, pinyins):
            rev = {}
            for p in pinyins:
                (word, tones) = pinyinRE.match(p).groups()
                default(rev, set)[word].update(int(t) for t in tones)
            return rev == reverse[c]

        def correction(phrase):
            return ' '.join(
                '/'.join(pinyinTones(*i) for i in reverse[c].items())
                for c in phrase
            )

        def lookup(c, pinyintones):
            lookups = ', '.join(
                '%s = %s' % (
                    accent(pinyin, tone),
                    forward.get((pinyin, tone), '\033[1;30m∅\033[0m'),
                )
                for (pinyin, tones) in pinyintones
                for tone in sorted(tones)
                for tone in [int(tone)]
                if c not in forward.get((pinyin, tone), '')
            )
            return lookups and '(' + lookups + ')'

        def lookuptext(c, pinyins):
            return lookup(c, (pt for p in pinyins for pt in pinyinRE.findall(p)))

        prompter = Prompter()
        with SRSQueue(data) as q:
            goods = ''
            tests = 0
            prevChar = None
            rounds = 100
            while True:
                togo = q.queue.index(None) if focus else rounds - tests
                if togo <= 0:
                    break
                char = q.head()
                done = False
                attempt = 0
                failed = False
                while True:
                    attempt += 1
                    score = q.scores.get(char, 0)
                    # Print new characters in bold.
                    try:
                        text, color = prompter.ask(char, score, togo)
                    except EOFError:
                        done = True
                        break

                    skipped = not text
                    if skipped:
                        prompter.check(False, False, '\v%s %s' % (
                            correction(char),
                            lookup('!', reverse[char].items()),
                        ))
                        goods = ''
                        q.skip(char)
                        prevChar = char
                        continue

                    easy = prevChar == char
                    if not easy:
                        prevChar = char
                        tests += 1

                    pinyins = [''.join(p).replace('v', 'ü') for p in pinyinRE.findall(text)]
                    sylls = [
                        s
                        for p in pinyins
                        for m in [pinyinRE.match(p)]
                        if m
                        for s in [m.group(1)]
                    ]
                    if not prompter.check(
                        False,
                        sum(len(w) for w in pinyins) == len(wsRE.sub('', text)),
                        '\033[1;31munrecognised elements:\033[0m %s', re.sub('\s+', ' ', pinyinRE.sub(' ', text).strip()),
                    ) or not prompter.check(
                        False,
                        all(s in syllables for s in sylls),
                        '\033[1;31minvalid syllable%s: %s',
                        's' if sum(s not in syllables for s in sylls) > 1 else '',
                        ' '.join(s for s in sylls if s not in syllables),
                    ):
                        print('\033[A', end='')
                    elif prompter.checknowait(True, check(char, pinyins), ''):
                        q.good(char, easy)
                        if goods:
                            print('\033[A', end='')
                        goods += char
                        score = q.scores.get(char, 0)
                        sr = scorerepr(score)
                        goodslen = 2*len(goods) + len(sr)
                        lines = 1 + (goodslen - 1)//termwidth()
                        print('\033[%dA\033[1;32m%s\033[0;32m%s\033[0m\033[J' % (lines, goods, sr))
                        break
                    else:
                        failed = True
                        q.bad(char, easy, attempt)
                        goods = ''
                        score = q.scores.get(char, 0)
                        print('\033[A\033[5C\033[%s31;9m%s\033[0m %s \033[1;30m%s\033[0m\033[K'
                            % (color, text, lookuptext(char, pinyins), scorerepr(score)))
                if done:
                    print()
                    break

            print('Completed %d rounds. 🎉' % (tests,))
            print(q.report())

            print('I:', focusreport(fscores))
            fscores = [scores.get(c, 0) for c in focuschars]
            print('O:', focusreport(fscores))

if __name__ == '__main__':
    main()
