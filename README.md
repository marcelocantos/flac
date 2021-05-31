# flac

Command-line flash-card application for learning 中文字.

# Install

Clone this repo.

# Run

Install Python 3, then…

```bash
./flac.py
```

# HOWTO

For each character shown, enter the numerical pinyin form. E.g., if shown 扎,
for which the pinyin is "zā", "zhā" and "zhá", enter `za1zha12` (`za1 zha12` and
`zha1za1 zha2` are also acceptable).

When you don't know the answer, just press <kbd>return</kbd> and the answer will
be shown in accented pinyin.

Braille-dots indicate a character's score, e.g.: ⣿⡄ indicates a score of 10.
Each correct answer bumps the score up by one dot. Each incorrect answer bumps
it down by one dot initially, but more dots for each consecutive incorrect
answer. Pressing <kbd>return</kbd> to reveal the answer bumps the score down by
several dots. A character's score also indicates how far down the queue it will
jump after a correct answer (a poor approximation of
[SRS](https://en.wikipedia.org/wiki/Spaced_repetition)). 

# TODO

- Compound words
- Test knowledge of definitions (currently only shows definitions after a
  correct response).
