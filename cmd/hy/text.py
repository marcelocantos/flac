import os
import re
import wcwidth

termwidth = os.get_terminal_size().columns

def wrapWords(words, width=None, color='', prefix=''):
    if width is None:
        width = termwidth
    pwidth = sum(printedWidths(prefix))
    sep = f'\n{"":{pwidth}}'
    wpl = (width - pwidth) // 5
    lines = [" ".join(words[i:i+wpl]) for i in range(0, len(words), wpl)]
    print(f'{prefix}{color}{sep.join(lines)}\033[0m')


nonprint = re.compile(r'[\b]|\033(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')

def printedWidths(text):
    return [
        wcwidth.wcswidth(nonprint.sub('', line))
        for line in text.split('\n')
    ]

def dedup(lst):
    seen = set()
    deduped = lst[:0]
    for (i, w) in enumerate(lst):
        if w not in seen:
            seen.add(w)
            deduped += lst[i:i+1]
    return deduped
