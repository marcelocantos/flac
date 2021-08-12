import * as InputRE from './InputRE';

it('empty', () => {
  expectMatch("()");
});

it('bad start', () => {
  expectNoMatch("1");
  expectNoMatch("12");
  expectNoMatch("/");
  expectNoMatch("!");
});

it('good start', () => {
  expectMatch("(a)");
  expectMatch("(a1)");
  expectMatch("(a2) ");
  expectMatch("(a3)      ");
  expectMatch("(a4)");
  expectMatch("(a5)");
  expectMatch("(a12)");
  expectMatch("(yi)");
  expectMatch("(yi1)");
  expectMatch("(yi12)");
  expectMatch("(zhuang12345)");
});

it('bad tones', () => {
  expectNoMatch("a0");
  expectNoMatch("a6");
  expectNoMatch("a9");
  expectNoMatch("yi0");
  expectNoMatch("yi6");
  expectNoMatch("yi9");
  expectNoMatch("zhuang0123456789");
});

it('bad first char', () => {
  expectNoMatch("yi/");
  expectNoMatch("12/");
  expectNoMatch("yi9/");
});

it('good first char', () => {
  expectMatch("yi1/()");
  expectMatch("yi12 / ()");
  expectMatch("yi1 / (s)   ");
  expectMatch("yi1/(shi)");
  expectMatch("yi1/(shi4)");
});

it('several chars', () => {
  expectMatch("jiang1/()");
  expectMatch("jiang14/()");
  expectMatch("jiang14/(q)");
  expectMatch("jiang14/(qiang)");
  expectMatch("jiang14/(qiang1)");
  expectMatch("yi1ding1(bu4)");
  expectMatch("yi1ding1bu4(s)");
  expectMatch("yi1ding1bu4 (s)");
  expectMatch("yi1 ding1 bu4 (s)");
  expectMatch("yi1 ding1 bu4 (shi)");
  expectMatch("yi1 ding1 bu4 (shi2)");

  expectMatch("nuo2(na)");
  expectMatch("nuo2(na1)");
  expectMatch("nuo2(na13)");
  expectMatch("nuo2(na134)");

  expectMatch("Ya4 dang1 ·()");
  expectMatch("Ya4 dang1 · ()");
  expectMatch("Ya4 dang1 · (S)");
  expectMatch("Ya4 dang1 · (Si1)");
  expectMatch("Ya4 dang1 · (Si1) ");
  expectMatch("Ya4 dang1 · Si1 (mi4)");
  expectMatch("Ya4 dang1 · Si1 (mi4)");

  expectMatch("yi1 bu4 zuo4,()");
  expectMatch("yi1 bu4 zuo4, ()");
  expectMatch("yi1 bu4 zuo4 , ()");
  expectMatch("yi1 bu4 zuo4 , (e)");
  expectMatch("yi1 bu4 zuo4 , (er4)");
  expectMatch("yi1 bu4 zuo4 , (er4) ");
  expectMatch("yi1 bu4 zuo4 , er4 bu4 (xiu1)");
});

it('bad several chars', () => {
  expectNoMatch("1 ding1 bu4 shi2");
  expectNoMatch("yi1 ding bu4 shi2");
  expectNoMatch("yi1 ding1 bu4 shi0");

  expectNoMatch("yi1 bu4 zuo,");
});

const expectedInputRE = /^(.*)\((.*)\)(.*)$/;

// expectMatch takes patterns of the form "prefix(lastchar)", expecting that
// inputRE matches prefix+lastchar and captures lastchar as group 1.
function expectMatch(pattern: string) {
  const e = pattern.match(expectedInputRE);
  expect(e).not.toBeUndefined();
  if (e) {
    const all = e[1] + e[2] + e[3];
    const expected = e[2];

    const m = all.match(InputRE.inputRE);
    expect(m).not.toBeUndefined();
    if (m) {
      expect(m[1]).toEqual(expected);
    }
  }
}

function expectNoMatch(input: string) {
  const e = input.match(InputRE.inputRE);
  expect(e).not.toBeUndefined();
}
