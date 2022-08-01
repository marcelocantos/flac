ESC = '\033'
CSI = f'{ESC}['

def code(*args):
    (args, cmd) = args[:-1], args[-1]
    return f'{CSI}{";".join(map(str, args))}{cmd}'

def optargs(cmd):
    class C:
        def __str__(self):
            return f'{CSI}{cmd}'

        def __call__(self, *args):
            return code(*args, cmd)

    return C()

up = optargs('A')
down = optargs('B')
forward = optargs('C')
back = optargs('D')

save = f'{CSI}s'
restore = f'{CSI}u'

def color(c):
    class C:
        def __str__(self):
            return f'\033[{c}m'

        @property
        def hi(self):
            return f'\033[1;{c}m'

        @property
        def lo(self):
            return f'\033[0;{c}m'

        @property
        def reset(self):
            return f'\033[0;{c}m'

    return C()

def rgb(r, g, b):
    return color(f'38;2;{r};{g};{b}')

def color256(c):
    return color(f'38;5;{c}')

lo = color(0)
hi = color(1)
italic = color(3)
underline, nounderline = color(4), color(24)
blink, noblink = color(5), color(25)
invert = color(7)

(black, red, green, yellow, blue, magenta, cyan, white) = [color(c) for c in range(30, 38)]

darkgray = black.hi

class eraseMeta(type):
    def __str__(cls):
        return f'{CSI}J'

class erase(metaclass=eraseMeta):
    def __init__(self, n):
        self.n = n

    def __str__(self):
        return f'{CSI}{self.n}J'
