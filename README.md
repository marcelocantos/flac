# flac

Command-line flash-card application for learning 中文字.

Why the name? It's just a truncation of "flash-card". Kind of silly, really.

## Install

- (Recommended) Download the program for your platform from
  <https://github.com/marcelocantos/flac/releases>.

- (For geeks) Install from source:

    1. Install Go. Options:

        - <https://golang.org/doc/install>
        - Mac with [Homebrew](https://brew.sh/):

           ```sh
           brew install go
           ```

    1. Clone this repo and run the following command:

        ```sh
        go install ./cmd/flac`.
        ```

## Run

```bash
./flac
```

Press Ctrl-C when you get bored. Progress will be saved as you go. The next time
you run flac, it will continue where it left off.

<!-- Not yet implemented in Go version.
Test yourself with a fixed corpus of characters:

```bash
./flac.py --focus 秋天的后半夜，月亮下去了，太阳还没有出，只剩下一片乌蓝的天；除了夜游的东西，什么都睡着。华老栓忽然坐起身，擦着火柴，点上遍身油腻的灯盏，茶馆的两间屋子里，便弥满了青白的光。
```
-->

## HOWTO

For each character shown, enter the numerical pinyin form. E.g., if shown 扎,
for which the pinyin is _zā_, _zhā_ and _zhá_, enter `za1zha12` (`za1 zha12` and
`zha1 za1 zha2` are also acceptable, as long as all forms are present).

For multi-character words, alternate pinyin forms must be separated by `/`. For
instance, the answer for 不是 can be entered as `bu2shi5/bu4shi4` (spaces
optional). The two forms can be entered in either order.

When you don't know the answer, just press <kbd>esc</kbd> and the answer will
be shown in accented pinyin.

## ETC

### Scoring

Braille-dots indicate a character's score, e.g.: ⣿⡄ indicates a score of 10.
Each correct answer bumps the score up by one dot. Each incorrect answer bumps
it down by slightly more than one dot initially, and more dots for each
consecutive incorrect answer. Pressing <kbd>esc</kbd> to reveal the answer
bumps the score down by several dots.

A partial answer (incomple word or missing alternatives) won't change the score,
and will offer a hint that more input is required.

An answer that is correct aside from tones, will be scored as an incorrect
answer, but will also provide a hint that only the tones need to be corrected.

A character's score also indicates how far down the queue it will jump after a
correct answer, in a crude approximation of
[SRS][6].

### Random vs contextual learning

Learning groups of words in the context of a narrative or theme is
[considered][7] the most
effective way to build a foreign language vocabulary. It is considerably better
than learning related words, such as colors, fruits, or modes of transport.

Yet another approach, learning words randomly, with no connection to each other,
is slightly worse than learning them contextually, but still much better than
learning related words. It is also vastly easier to implement for large
collection of words, since it doesn't require human effort to construct
meaningful stories. This is the approach taken by flac.

### Sources

- The list of words used for testing is sourced from an attachment in a [Pleco
  Forums](https://www.plecoforums.com/) post titled [_Word frequency list based
  on a 15 billion character corpus: BCC (BLCU Chinese Corpus)_][1]. The post
  itself derives the attachements from the Beijing Language and Culture
  University's 15-billion character corpus ([zip archive][2]).

  Flac uses the top 10,000 words from the global word-frequency list.

- Definitions are sourced from [CC-CEDICT][3]. A number of missing terms
  were added by hand (see [addenda.txt][4] from the flac source), mostly by
  looking them up in [Google Translate](https://translate.google.com/).

#### Updating source data

Currently, flac binds all source data directly into the program. Updating the
lists or choosing different sources requires rebuilding from source. If you'd
like instructions for how to do this, [raise an issue][5].

### Saving progress

Flac saves progress as you go. The program creates a database file, `flac.db`,
in the current directory, which updates the word queue and their scores as you
go, so that flac can remember your progress one session to the next. Flac
doesn't currently support multiple profiles, but if you rename this file, flac
will create a fresh file and start from scratch. You can thus manage profiles
manually by fiddling with copies of `flac.db`.

## TODO

In rough order of priority:

1. Test knowledge of definitions (currently only shows definitions after a
   correct response).
   - Could be text entry or multiple choice.
   - Could be a second step after getting the pinyin form(s) correct.
1. Instead of requiring entry of all pinyin alternative forms, test for a
   specific form with a unique definition as the hint.
1. Test knowledge of the 汉字 for a given pinyin form (with or without hints).
   - Since we can't support the user drawing characters on a text UI, this will
     have to be multiple choice.
   - Actually, we can do mouse support, so perhaps we could support drawing with
     box-drawing characters!
1. Save session history view across sessions.
1. Reintroduce `--focus` option from the old Python version. But the new version
   should create a separate queue instead of polluting the normal queue.
1. Auto-detect flagging performance and end the session with a suggestion to
   take a break.
1. Test phrases.
1. Integrate [TTS][8] and [PortAudio][9]. Uses:
   - Hints
   - Listening exercises
1. Use speech recognition to test speaking.
1. Research science-based SRS models and implement one.

  [1]: https://www.plecoforums.com/threads/word-frequency-list-based-on-a-15-billion-character-corpus-bcc-blcu-chinese-corpus.5859/
  [2]: http://bcc.blcu.edu.cn/downloads/resources/BCC_LEX_Zh.zip
  [3]: https://www.mdbg.net/chinese/dictionary?page=cc-cedict
  [4]: https://github.com/marcelocantos/flac/blob/master/refdata/addenda.txt
  [5]: https://github.com/marcelocantos/flac/issues/new
  [6]: https://en.wikipedia.org/wiki/Spaced_repetition
  [7]: https://blog.fluent-forever.com/base-vocabulary-list/
  [8]: http://www.voicerss.org/api/demo.aspx
  [9]: https://pkg.go.dev/github.com/gordonklaus/portaudio?utm_source=godoc
