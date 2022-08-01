#!/usr/bin/env python3

# Detect characters subsumed by bisyllable (i.e., not in HSK list as singles).

# Built-in
import dicts

# Local
import sheets

conn = sheets.Connection()

data = conn.get([
    f'{s.tab}!A1:D5'
    for s in sheets.statuses
])
words = [
    word
    for t in data
    for row in t[0]
    for cell in row
    if cell
    for word in cell.split()
]

chars = {c for w in words for c in w}

def subsumed(dct):
    dictchars = {w for w in dct.keys() if len(w) == 1}
    return len(chars - dictchars)

all = len(chars)

print('characters subsumed in multisyllabic words:')
print(f'CEDICT: {100*subsumed(dicts.cedict)//all}%')
print(f'HSK: {100*subsumed(dicts.hskdict)//all}%')
