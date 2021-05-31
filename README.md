# flac

Command-line flash-card application for learning 中文字.

# Install

Clone this repo.

# Run

Install Python 3, then…

Test yourself till you get 100 correct responses.

```bash
./flac.py
```

Test yourself, focusing on a fixed corpus of characters:

```bash
./flac.py --focus 秋天的后半夜，月亮下去了，太阳还没有出，只剩下一片乌蓝的天；除了夜游的东西，什么都睡着。华老栓忽然坐起身，擦着火柴，点上遍身油腻的灯盏，茶馆的两间屋子里，便弥满了青白的光。
```

The first time you run it, flac orders characters randomly. Thereafter, it
preserves the queue from one run to the next. It will even remember the focused
character order when you rerun flac without `--focus`.

# HOWTO

For each character shown, enter the numerical pinyin form. E.g., if shown 扎,
for which the pinyin is "zā", "zhā" and "zhá", enter `za1zha12` (`za1 zha12` and
`zha1za1 zha2` are also acceptable).

When you don't know the answer, just press <kbd>return</kbd> and the answer will
be shown in accented pinyin.

To exit before finishing a run, press <kbd>Ctrl-D</kbd>. <kbd>Ctrl-C</kbd> also
exits, but won't save your progress.

# ETC

Braille-dots indicate a character's score, e.g.: ⣿⡄ indicates a score of 10.
Each correct answer bumps the score up by one dot. Each incorrect answer bumps
it down by one dot initially, but more dots for each consecutive incorrect
answer. Pressing <kbd>return</kbd> to reveal the answer bumps the score down by
several dots. A character's score also indicates how far down the queue it will
jump after a correct answer (a poor approximation of
[SRS](https://en.wikipedia.org/wiki/Spaced_repetition)). 

# TODO

- Offer to save progress after pressing Ctrl-C.
- Test compound words.
- Test knowledge of definitions (currently only shows definitions after a
  correct response).
- Switch queue and score persistence from pickling to sqlite.
