cp -n master/config/config.example master/config/config.go

sh kitexgen.sh

mkdir -p output


mkdir -p sniffer/output
cd sniffer
go build -o ../output/argos.sniffer

cd ../master
go build -o ../output/argos.master

cd ..