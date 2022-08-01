import re

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
    tone = int(tone)
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

pinyinRE = re.compile(r'(?i)([a-zü]+)(\d+)')
pinyinsRE = re.compile(r'(?i)(?:(?:[\u3000-\u9FFF]+\||)?([\u3000-\u9FFF]+))?\[((?:[a-zü]+\d+\s+)*[a-zü]+\d+)\]')
tradcharRE = re.compile(r'(?:[\u3000-\u9FFF]+\|)(?=[\u3000-\u9FFF]+)')

def accents(pinyins):
    return pinyinRE.sub(lambda m: accent(*m.groups()), pinyins)

def accentsInPhrase(phrase):
    phrase = tradcharRE.sub('', phrase)
    return pinyinsRE.sub(
        lambda m: f"\033[1m{m[1] or ''}[\033[0m{accents(m[2])}\033[1m]\033[0m",
        phrase)
