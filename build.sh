cp -n master/config/config.example master/config/config.go

sh kitexgen.sh

mkdir -p sniffer/output
cd sniffer
go build -o output/bin/argos.sniffer

cd ../master
sh build.sh

cd ..