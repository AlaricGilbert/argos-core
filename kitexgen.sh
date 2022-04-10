mkdir -p sniffer
mkdir -p master

cd sniffer 
kitex -module github.com/AlaricGilbert/argos-core -service argos.sniffer ../thrift/sniffer.thrift
cd ../master
kitex -module github.com/AlaricGilbert/argos-core -service argos.master ../thrift/master.thrift
cd ..