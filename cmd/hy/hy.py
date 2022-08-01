#!/usr/bin/env python3

# Built-in
import collections
import functools
import json
import logging
import os
import random
import re
import sys
import time

# Local
import ansi
from data import splicer
import dicts
import sheets
import term
import text
import tts

STATE_FILE_PATH = '.hy.state'

def gammavariate(µ, sd):
    k = µ**2 / sd
    θ = sd / k
    # Parameters ⍺ and β actually correspond to k and θ as described in
    # https://en.wikipedia.org/wiki/Gamma_distribution
    return random.gammavariate(k, θ)

# PDFs peak at 9, 6, 4, 4, 3, 3, 3, 2, 2, 2, ...
# for score = 0, -2, -4, -6, ...
def scorebump(score):
    µ = 50/(4-score)**1.1
    return 2 + int(gammavariate(µ, µ/2))

def put(s):
    print(s, end='', flush=True)

def report(before, after):
    for (j, (t, b)) in enumerate(zip(sheets.statuses, after)):
        words = [
            f'{s.color}{c}'
            for (i, (s, a)) in sorted(
                enumerate(zip(sheets.statuses, before)),
                key=lambda x: x[0] != j,
            )
            for c in set(a) & set(b)
        ]
        text.wrapWords(
            list(words),
            prefix=f'{t.color}{t.key}: {ansi.lo}',
        )

def outcomeMask(outcome):
    return 1 << 'plkjn'.index(outcome)

outcomeLabels = {
    'n': f'{ansi.rgb(240,  70,  70).hi}✘',
    'j': f'{ansi.rgb(180, 100, 100).hi}✗',
    'k': f'{ansi.green.hi}✔',
    'l': f'{ansi.yellow.hi}★',
    'p': f'{ansi.black.hi}‧',
}

def main():
    logging.basicConfig(
        filename='hy.log',
        encoding='utf-8',
        level=logging.DEBUG,
        format='%(asctime)s.%(msecs)03d:%(levelname)s:%(name)s:%(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
    )

    conn = sheets.Connection()
    ttsconn = tts.Connection()
    audio = True

    def say(word, always=False):
        if audio or always:
            ttsconn.getWordWav(word).play(nonblocking=True)
            logging.debug(f"said {word}")

    def onoff(state):
        return "on" if state else "off"

    oneshot = False
    logfile = None
    tones = None
    cell = None
    x, y = None, None

    argv = list(reversed(sys.argv[1:]))
    while argv:
        arg = argv.pop()

        if arg == '--check':
            conn.checkAllCells()
            return

        if arg == '-1':
            oneshot = True
        elif arg == '--log':
            logfile = open(argv.pop(), 'a+')
        elif re.match(r'\d\d', arg):
            tones = arg
        else:
            print(f"bad arg: {arg}")
            return

    state = None
    resumeMessage = ''
    if os.path.isfile(STATE_FILE_PATH):
        with open(STATE_FILE_PATH) as f:
            state = json.load(f)
        os.unlink(STATE_FILE_PATH)
        tones = state['tones']
        resumeMessage = f' {ansi.black.hi}(resumed from previous run...){ansi.lo}'

    if tones:
        x = int(tones[0]) - 1
        y = int(tones[1]) - 1
    else:
        (progressScores, _) = conn.get("两字符词!scores")[0]

        (_, (x, y)) = min(
            (v, (i, j))
            for (j, row) in enumerate(progressScores)
            for (i, v) in enumerate(row)
        )

        tones = '1234'[x] + '12345'[y]

    cell = 'ABCD'[x] + '12345'[y]

    result = conn.fetchGroups(cell, tones)
    if result is None:
        return

    (groups, ids) = result

    shuffled = []
    for g in groups:
        g = g.copy()
        random.shuffle(g)
        shuffled.append(g)

    words = [w for g in shuffled[:3] for w in g]
    nwords = len(words)

    eta = nwords * 3 // 2
    etastr = f'{eta % 60}m'
    if eta >= 60:
        etastr = f'{eta // 60}h{etastr}'

    print(f'tones {"-".join(tones)}, {nwords} words, ETA {etastr}, {nwords - len(groups[3])} to test')
    for (s, g) in zip(sheets.statuses, groups):
        text.wrapWords(g, color=s.color, prefix=f'{s.key}: ')
    print()

    origWords = frozenset(words)

    history = []
    done = scores = seen = outcomes = worst = None

    def pushCheckpoint():
        history.append({
            'words': words,
            'done': {s: w.copy() for (s, w) in done.items()},
            'scores': scores.copy(),
            'seen': seen.copy(),
            'outcomes': outcomes.copy(),
            'worst': worst.copy(),
        })

    def setCheckpoint(cp):
        nonlocal words, done, scores, seen, outcomes, worst
        words = cp['words']
        done = cp['done']
        scores = cp['scores']
        seen = cp['seen']
        outcomes = cp['outcomes']
        worst = cp['worst']

    def popCheckpoint():
        setCheckpoint(history.pop())

    # State variables
    if state:
        for h in state['history']:
            history.append({
                'words': tuple(h['words']),
                'done': {s: set(w) for (s, w) in h['done'].items()},
                'scores': collections.defaultdict(int, h['scores']),
                'seen': set(h['seen']),
                'outcomes': collections.defaultdict(int, h['outcomes']),
                'worst': collections.defaultdict(int, h.get('worst', {})),
            })
        setCheckpoint(history[-1])
    else:
        words = tuple(words)
        done = {s.key: set() for s in sheets.statuses}
        scores = collections.defaultdict(int, **{w: 0 for w in words})
        seen = set()
        outcomes = collections.defaultdict(int)
        worst = collections.defaultdict(int)

    doubleCheck = False
    showingInfo = False

    lens = [len(g) for g in groups]
    logging.info(f"groups: {' + '.join(map(str, lens))} = {sum(lens)}")
    logging.info(f"开始了 {tones}: {' ◆ '.join(' '.join(g) for g in groups)}")

    def remaining():
        return sum(s <= 0 for s in scores.values())

    def output(s):
        space = " " * (not s.startswith('\n'))
        lines = sum(
            # Add 5% for safety margin.
            pwidth*105//100 // text.termwidth + 1
            for pwidth in text.printedWidths(s)
        )
        vtabs = '\v' * lines
        modes = f'{ansi.black.hi}{"1" * oneshot}{"⾳" * audio}{ansi.lo}'
        put(f'{vtabs}{ansi.up(lines)}{ansi.save}  {modes}{space}{s}{ansi.erase}{ansi.restore}')

    repeated = False

    hi = ansi.yellow.hi
    lo = ansi.lo
    firstTime = True
    lastFanfare = remaining()
    won = False
    startTime = before = time.time()
    totalTime = 0
    initRemaining = remaining()
    defaultRate = 1/90
    temporalGaps = []

    def fmttime(t):
        return time.strftime("%H:%M", time.localtime(t))

    def eta():
        left = remaining()
        completed = initRemaining - left
        defaultWeight = 0.9 ** completed
        rate = defaultWeight * defaultRate
        if completed:
            rate += (1 - defaultWeight) * completed / totalTime
        return fmttime(time.time() + left / rate)

    # words = ('家庭',) + words
    with term.raw():
        while remaining() and words:
            try:
                if logfile is not None:
                    print(words, file=logfile)

                if not repeated:
                    pushCheckpoint()

                word = words[0]

                status = None
                inStatuses = []
                for s in sheets.statuses:
                    d = done[s.key]
                    if word in d:
                        d.discard(word)
                        status = s
                        inStatuses.append(s.key)
                if len(inStatuses) > 1:
                    logging.error(f'Too many statuses: {inStatuses}')

                if not status:
                    for (group, status) in zip(groups, sheets.statuses):
                        if word in group:
                            break

                level = -scores[word]//2
                scorestr = (
                    f'{ansi.magenta.hi}*'
                        if level == 0 and word in seen
                    else ansi.black.hi + ''.join(
                        chr(ord(d) - ord('0') + ord('₀'))
                        for d in str(level)
                    )
                        if scores[word] <= 0
                    else f'{ansi.hi}₊'
                )
                scorestrwidth = sum(text.printedWidths(scorestr))
                shownew = '' if word in seen else f'{ansi.underline}'
                hidenew = ansi.nounderline
                if not repeated:
                    put(f'{status.color}{shownew}{word}{hidenew}{scorestr}'
                        f'{ansi.lo} {ansi.erase}\v{ansi.up}'
                        f' {ansi.back}' # Clean up end-of-line behaviour.
                    )
                    if firstTime:
                        firstTime = False
                        output(
                            f'{ansi.black.hi} Press {ansi.rgb(199, 196, 0)}?{ansi.black.hi} for help{ansi.lo}'
                            f'{resumeMessage}')
                    elif not won:
                        output('')
                    else:
                        left = remaining()
                        drought = lastFanfare - left
                        # randint(5, 22) yields an average drought length of about 10.
                        if random.randint(5, 22) < drought:
                            # Using a weighted mean of a default rate and the
                            # measured rate is probably overkill, but occasionally
                            # (and often during testing) an early fanfare will
                            # benefit from the noise reduction it offers.
                            # TODO: exponentially weight recent progress higher.
                            another = '' if lastFanfare == initRemaining else 'another '
                            fanfare = (
                                f'🥳 {ansi.blink}🎈 🎈{ansi.noblink} {ansi.white.hi}{another}{drought}{ansi.lo} down, '
                                f'{ansi.white.hi}{left}{ansi.lo} to go! '
                                f'(ETA {eta()}) {ansi.blink}🎈 🎈{ansi.noblink} 🥳'
                            )
                            offset = (text.termwidth - sum(text.printedWidths(fanfare)))//2
                            output(f'\n\n{"":{offset}}{fanfare}')
                            lastFanfare = left
                else:
                    repeated = False

                if not repeated:
                    won = False

                c = term.getch()
                logging.debug('key pressed')

                delta = time.time() - before
                before = time.time()
                # After 30s, your attention is likely focused elsewhere. Ignore it.
                if delta < 30:
                    totalTime += delta
                else:
                    temporalGaps.append((before, delta))

                if c == 'q':
                    if doubleCheck != 'q':
                        doubleCheck = 'q'
                        output("Press 'q' again to quit.")
                        repeated = True
                        continue
                    else:
                        break

                if c == '\x03':
                    print('^C')
                    return

                if c == 'z':
                    with open(STATE_FILE_PATH, 'w') as f:
                        json.dump({
                            'tones': tones,
                            'history': [
                                {
                                    'words': words,
                                    'done': {s: list(w) for (s, w) in done.items()},
                                    'scores': list(scores.items()),
                                    'seen': list(seen),
                                    'outcomes': list(outcomes.items()),
                                }
                                for h in history[-1:]
                            ],
                        }, f)
                    print(f' Progress saved (will auto-resume on next run).{ansi.erase}')
                    return

                if c == '\n':
                    say(word, always=True)
                    repeated = True

                elif c == 'a':
                    audio = not audio
                    space = '  ' * (not audio)
                    output(f'{space}audio {onoff(audio)}')
                    repeated = True

                elif c == '1':
                    oneshot = not oneshot
                    space = ' ' * (not oneshot)
                    output(f'{space}1-shot mode {onoff(oneshot)}{" (you get one shot at each word)" * oneshot}')
                    repeated = True

                elif c in done or c == 'p':
                    if c != 'p' and not doubleCheck:
                        say(word)
                    oldScore = scores[word]

                    logging.debug('outcome key pressed')
                    if c == 'n':
                        newScore = min(scores[word] - 3, -4)
                    elif c == 'j':
                        newScore = min(scores[word] - 2, -2)
                    elif c == 'k':
                        newScore = scores[word] + 2
                    elif c == 'l':
                        newScore = 1
                    else:
                        newScore = scores[word]

                    won = oldScore <= 0 and newScore > 0

                    if status.key == 'l' and c == 'k':
                        put('\a')
                        c = 'l'

                    outcomes[c] += 1
                    outcomeLabel = outcomeLabels[c]

                    worst[word] |= outcomeMask(c)

                    if newScore <= 0:
                        seen.add(word)
                        scores[word] = newScore
                        doubleCheck = False
                        if c == 'k' and not oneshot:
                            bump = scorebump(newScore)
                            words = splicer(words)[1:bump, 0, bump:]
                        else:
                            words = words[1:]
                        if c == 'k':
                            c = 'j'
                    elif oldScore <= 0 and doubleCheck != '*' and showingInfo != ' ':
                        output(
                            f'{ansi.black.hi}double-check:{ansi.lo} '
                            f'{dicts.formatdefinition(dicts.cedict, word, tones)}'
                        )
                        doubleCheck = '*'
                        repeated = True
                        continue
                    else:
                        seen.add(word)
                        scores[word] = newScore
                        words = words[1:]
                        if len(words) < 4:
                            recycledWord = random.choice(list(origWords - set(words)))
                            words += (recycledWord,)

                    if c == 'p':
                        newStatus = status
                    else:
                        done[c].add(word)
                        newStatus = sheets.statusForKey[c]

                    put(f'{ansi.back(5+scorestrwidth)}{newStatus.color}'
                        f'{shownew}{word}{hidenew}{outcomeLabel}{ansi.lo} ')

                elif c in ' ?~dhs': # Toggle definition
                    put(ansi.erase)
                    showingInfo = showingInfo != c and c
                    if not showingInfo:
                        output('')
                        doubleCheck = False
                    elif showingInfo == ' ':
                        output(
                            dicts.formatdefinition(dicts.cedict, word, tones) +
                            '\n' +
                            dicts.syllableDefs(word, tones, False)
                        )
                    elif showingInfo == '?':
                        output('\n' + '; '.join(
                            f'''{"/".join(f"{hi}{key}{lo}" for key in keys if key not in description)}{
                                ' ' * any(key not in description for key in keys)
                            }{
                                functools.reduce(
                                    lambda description, key: description.replace(key, f'{hi}{key}{lo}', 1),
                                    keys,
                                    description,
                                )
                            }'''
                            for (keys, description) in [
                                (['n', 'j', 'k', 'l'], 'outcome (nope, just pronounce, known, pool room)'),
                                (['⌫'], 'undo'),
                                (['p'], 'pass'),
                                (['h'], 'hint'),
                                (['z'], 'snooze'),
                                (['esc'], 'clear message'),
                                (['a'], f'toggle audio (currently {onoff(audio)})'),
                                (['1'], f'toggle 1-shot mode (currently {onoff(oneshot)})'),
                                (['?'], 'help'),
                                (['s'], 'stats'),
                            ] # + ([('~', 'debug info')] if debug else [])
                        ))
                    elif showingInfo == 'h':
                        output('\n' + dicts.syllableDefs(word, tones, True))
                    elif showingInfo == 's':
                        scoreBuckets = collections.defaultdict(set)
                        for (w, s) in scores.items():
                            if s:
                                scoreBuckets[s].add(w)

                        scoresReport = [
                            f"{s}:{''.join(f'{w[:1]}̳{w[1:]}' for (i, w) in enumerate(ww))}"
                            for (s, ww) in sorted(scoreBuckets.items())
                        ]

                        output(
                            f'\n{ansi.hi}progress:{ansi.lo} {remaining()} to go; ETA {eta()}'
                            f'\n{ansi.hi}outcomes:{ansi.lo} {" ".join(f"{o}={outcomes[o]}" for o in "njklp")}'
                            f'\n{ansi.hi}scores:{ansi.lo} {", ".join(scoresReport)}'
                            f'\n{ansi.hi}time gaps:{ansi.lo} {", ".join(f"{int(d)}s @ {fmttime(t)}" for (t, d) in temporalGaps)}'
                            f'{ansi.lo}'
                        )
                    elif showingInfo == '~':
                        output(f'DEBUG: {scorestrwidth = }')

                    repeated = True

                elif c == '\033': # Escape (clear output)
                    output('')
                    showingInfo = False
                    doubleCheck = False
                    repeated = True

                elif c == '\x7f': # Undo
                    popCheckpoint()
                    if history:
                        popCheckpoint()
                        put(ansi.back(11+scorestrwidth))
                    else:
                        put('\a')
                        pushCheckpoint()
                        repeated = True

                else:
                    repeated = True

            finally:
                if not repeated:
                    showingInfo = False
                    doubleCheck = False

    print(f'{ansi.erase}\n')
    unseen = set([w for g in groups for w in g])
    for s in sheets.statuses:
        unseen -= done[s.key]

    worst = {
        k: 5 - v.bit_length()
        for k, v in worst.items()
    }

    wordses = [
        text.dedup([
            # Preserve original order. Not required, but easier to test.
            w
            for j in list(range(i, 4)) + list(range(i))
            for w in groups[j]
            if i == j and (w not in worst or worst.get(w) == 4) or worst.get(w) == i
        ])
        for (i, s) in enumerate(sheets.statuses)
    ]
    report(groups, wordses)

    # Now shuffle them
    for words in wordses:
        random.shuffle(words)

    endTime = time.time()
    ndone = nwords - len(words)
    put(
        f'\nTotal time: '
        f'{int(endTime - startTime)//60} min on the wall, '
        f'{int(totalTime)//60} min active'
        f'; {ndone} words ({int(totalTime)/ndone:.3f} s/word)'
        f'\n'
        f'Save? (Y/n) '
    )
    while True:
        c = term.getch()
        if c in 'Yy\n':
            put('在保存...')
            # Reconnect to avoid token timeout issues.
            try:
                conn = sheets.Connection()
                conn.save(ids, x, y, wordses)
                logging.info(f"保存了 {tones}: {wordses}")
                print('完成了!')
            except:
                print(f'{ansi.save}  Something went wrong!{ansi.restore}')
                continue
        elif c in 'Nnq\x03':
            print('no')
        else:
            continue
        return

if __name__ == '__main__':
    try:
        # put(f'{ansi.CSI}H{ansi.erase}')
        # put(f'{ansi.ESC}]50;SetProfile=Big\a')
        main()
    finally:
        # put(f'{ansi.ESC}]50;SetProfile=Default\a')
        pass
