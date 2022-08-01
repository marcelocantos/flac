import collections
import os
import re

import ansi
from pinyin import accents, accentsInPhrase, pinyinsRE
import text

nonDigitsRE = re.compile(r'\D+')

DictEntry = collections.namedtuple('CedictEntry', 'word tones pinyin pinyinN definitions')

def loadcedict():
    partsRE = re.compile(r'^[\u4e00-\u9fff]+ ([\u4e00-\u9fff]+) \[([^]]+)\] /(.*)/$')
    cedict = collections.defaultdict(lambda: collections.defaultdict(list))
    for f in ['cedict_1_0_ts_utf-8_mdbg.txt', 'addenda.txt']:
        for line in open(f'../../refdata/{f}').readlines():
            line = line.strip().replace('u:', 'ü')
            for parts in [partsRE.match(line)]:
                if parts:
                    (word, pinyin, definitions) = parts.groups()
                    tones = nonDigitsRE.sub('', pinyin)
                    cedict[word][tones].append(DictEntry(
                        word=word,
                        tones=tones,
                        pinyin=accents(pinyin),
                        pinyinN=pinyin,
                        definitions=definitions.split('/'),
                    ))
    return cedict

cedict = loadcedict()

def loadhskdict():
    hskdict = collections.defaultdict(lambda: collections.defaultdict(list))
    for level in range(1, 7):
        filename = f'HSK Official With Definitions 2012 L{level}.txt'
        for line in open(os.path.expandvars(f'$HOME/play/hsk/{filename}')).readlines():
            if line[:1] == '\ufeff':
                line = line[1:]
            (word, _, pinyin, accentedPinyin, definition) = line.split('\t')[:5]
            tones=(nonDigitsRE.sub('', pinyin) + '5')[:2]
            hskdict[word][tones].append(DictEntry(
                word=word,
                tones=tones,
                pinyin=accentedPinyin,
                pinyinN=pinyin,
                definitions=[definition],
            ))
    return hskdict

hskdict = loadhskdict()

surnameRE = re.compile(r'^(?i)(surname) \w+$')
variantRE = re.compile(r'^(?:(?:old )?variant of|see) (?:[\w🙈]+)(?:\|([\w🙈]+)\[[^\]]*\])?$')

def see(dct, word):
    dictword = dct[word]
    if len(dictword) == 1:
        entries = next(iter(dictword.values()))
        if len(entries) == 1:
            entry = entries[0]
            defns = entry.definitions
            if len(defns) == 1:
                defn = defns[0]
                if defn.startswith('see '):
                    m = pinyinsRE.match(defn[len('see '):])
                    if m:
                        word = m[1]
                        dictword = dct[word]
    return word, dictword

def formatdefinition(
    dct,
    word,
    tones,
    tonesep=f" {ansi.hi}◆{ansi.lo} ",
    defsep=f" {ansi.hi}/{ansi.lo} ",
    light=False,
    pinyin=None,
):
    seeword, dictword = see(dct, word)
    prefix = f'see {ansi.hi}{seeword}{ansi.lo} → ' if seeword != word else ''
    word = seeword

    def denoise(s):
        m = variantRE.match(s)
        if m and m[1] == word:
            return ''
        return s

    if light:
        def hide(s):
            s = surnameRE.sub(r'\1 🙈', s)
            s = pinyinsRE.sub(r'🙈', s)
            s = variantRE.sub('', s)
            return denoise(s)
        def trim(e):
            return e.pinyinN in pinyin
    else:
        hide = denoise
        trim = None

    if not pinyin:
        pinyin = {e.pinyinN for ee in dictword.values() for e in ee}

    entries = [
        f'{"" if light else f"[{e.pinyin}] "}{defsep.join(defs)}'
        for (_, ee) in sorted(dictword.items(), key=lambda te: te[0] != tones)
        for e in sorted(filter(trim, ee), key=lambda e: [e.pinyinN not in pinyin, e.pinyinN[0].isupper()])
        for defs in [list(filter(None, [accentsInPhrase(hide(d)) for d in e.definitions]))]
        if defs
    ]
    return f'{prefix}{tonesep.join(filter(None, entries))}'

def syllableDefs(word, tones, light):
    word, _ = see(cedict, word)
    pinyins = zip(*[e.pinyinN.split() for e in cedict[word][tones]])
    return f''.join(
        f'\n{ansi.white.hi}{c}{ansi.lo} 👉 {formatdefinition(cedict, c, t, light=light, pinyin=p)}'
        for (c, t, p) in zip(text.dedup(word), tones, pinyins)
    )
