go build -o cnvtr *.go
for i in ${*}
do
	./cnvtr ${i} > ${i}.re
done
