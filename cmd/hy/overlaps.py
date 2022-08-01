#!/usr/bin/env python3

import collections
import dicts

revmap = collections.defaultdict(list)

for (word, entries) in dicts.cedict.items():
    for (tone, entry) in entries:
        for defn in definitions:
            pass
