import time

from reloading import reloading


def test1():
    for i in range(100):
        print(f"|i={i}|")
        time.sleep(1.0)


if __name__ == "__main__":
    test1()
