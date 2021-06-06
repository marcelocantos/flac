.PHONY : all
all : flac

empty :=
space := $(empty) $(empty)

CEDICT_PB = internal/proto/refdata/refdata.pb.go

CEDICT_SRCS = \
	refdata/cedict_1_0_ts_utf-8_mdbg.txt \
	refdata/addenda.txt

REFDATA_CACHE = internal/refdata/refdata.cache

WORDS = refdata/words.txt

GEN_FILES = $(CEDICT_PB) $(REFDATA_CACHE) $(WORDS)

GO_SRCS = $(CEDICT_PB) $(shell find . -name '*.go')

FLAC_SRCS = $(GO_SRCS) $(GEN_FILES)

# generate

.PHONY : gen
gen : $(GEN_FILES)

$(REFDATA_CACHE) : precache $(WORDS) $(CEDICT_SRCS)
	./precache -o $@ $(WORDS) $(subst $(space),:,$(CEDICT_SRCS))

GLOBAL_WORDFREQ_NAME = global_wordfreq-release_utf-8-txt.2593
GLOBAL_WORDFREQ_URL = \
	https://www.plecoforums.com/download/$(GLOBAL_WORDFREQ_NAME)/
GLOBAL_WORDFREQ = refdata/$(GLOBAL_WORDFREQ_NAME)

$(GLOBAL_WORDFREQ) :
	curl --silent --output $@ $(GLOBAL_WORDFREQ_URL)
	[ -s $@ ] || (rm -f $@ && echo "Empty file!" && false)

$(WORDS) : $(GLOBAL_WORDFREQ)
	head -n 10000 $< | awk '//{print $1}' > $@ || rm -f $@

# binaries

flac : $(FLAC_SRCS)
	go build ./cmd/flac

precache : $(GO_SRCS)
	go build ./internal/cmd/precache

# clean

.PHONY : clean-gen
clean-gen :
	rm -f $(GEN_FILES)

%.pb.go : %.proto
	protoc --go_out=. $<

.PHONY : clean
clean : clean-cache clean-gen

.PHONY : clean-flac
clean-flac :
	rm -f flac

.PHONY : clean-precache
clean-precache :
	rm -f precache

.PHONY : clean-cache
clean-cache :
	rm -f refdata/*.cache
