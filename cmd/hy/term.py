import contextlib
import sys
from termios import *
import tty

_stack = []

@contextlib.contextmanager
def raw():
    fd = sys.stdin.fileno()
    mode = tcgetattr(fd)
    _stack.append(mode)
    try:
        mode = tcgetattr(fd)
        mode[3] &= ~(ECHO | ICANON | ISIG)
        tcsetattr(fd, TCSAFLUSH, mode)
        yield
    finally:
        tcsetattr(fd, TCSADRAIN, _stack.pop())

@raw()
def getch():
    return sys.stdin.read(1)
