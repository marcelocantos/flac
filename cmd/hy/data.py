class splicer:
    def __init__(self, a):
        self._a = a

    def __getitem__(self, slices):
        if isinstance(slices, slice):
            slices = [slices]

        return [
            x
            for s in slices
            for x in (self._a[s] if isinstance(s, slice) else [self._a[s]])
        ]
