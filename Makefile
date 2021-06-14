.PHONY : all
all : flac

empty :=
space := $(empty) $(empty)

REFDATA_PB = internal/pkg/proto/refdata/refdata.pb.go

# addenda.txt must come second, so its removals get applied.
REFDATA_SRCS = \
	refdata/cedict_1_0_ts_utf-8_mdbg.txt \
	refdata/addenda.txt

REFDATA_CACHE = internal/pkg/refdata/refdata.cache

WORDS = refdata/words.txt

GEN_FILES = $(REFDATA_CACHE) $(WORDS)

GO_SRCS = $(REFDATA_PB) $(shell find . -name '*.go')

FLAC_SRCS = $(GO_SRCS) $(GEN_FILES)

PRECACHE_DIRS = \
	internal/cmd/precache \
	internal/pkg/pinyin \
	internal/pkg/proto/refdata

PRECACHE_SRCS = $(foreach dir,$(PRECACHE_DIRS),$(wildcard $(dir)/*.go))

# generate

.PHONY : gen
gen : $(REFDATA_PB) $(GEN_FILES)

$(REFDATA_CACHE) : precache $(WORDS) $(REFDATA_SRCS)
	./precache -o $@ $(WORDS) $(subst $(space),:,$(REFDATA_SRCS))

GLOBAL_WORDFREQ_NAME = global_wordfreq-release_utf-8-txt.2593
GLOBAL_WORDFREQ_URL = \
	https://www.plecoforums.com/download/$(GLOBAL_WORDFREQ_NAME)/
GLOBAL_WORDFREQ = refdata/$(GLOBAL_WORDFREQ_NAME)

$(GLOBAL_WORDFREQ) :
	curl --silent --output $@ $(GLOBAL_WORDFREQ_URL)
	[ -s $@ ] || (rm -f $@ && echo "Empty file!" && false)

$(WORDS) : $(GLOBAL_WORDFREQ)
	head -n 10000 $< | awk '//{print $1}' > $@ || rm -f $@

%.pb.go : %.proto
	protoc --go_out=. $<

# binaries

flac : $(FLAC_SRCS)
	go build ./cmd/flac

precache : $(PRECACHE_SRCS)
	go build ./internal/cmd/precache

# test

.PHONY : test
test :
	go test ./...

# clean

.PHONY : clean
clean : clean-gen clean-bin

.PHONY : deep-clean
deep-clean : clean clean-pb

.PHONY : clean-gen
clean-gen :
	rm -f $(GEN_FILES)

.PHONY : clean-bin
clean-bin :
	rm -f flac precache

.PHONY : clean-pb
clean-pb :
	rm -f $(REFDATA_PB)
