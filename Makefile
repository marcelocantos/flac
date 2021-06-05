.PHONY : gen
gen : internal/proto/cedict/cedict.pb.go


%.pb.go : %.proto
	protoc --go_out=. $<
