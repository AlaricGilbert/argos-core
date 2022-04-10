include "base.thrift"
namespace go sniffer

struct TimeRequest {
	1: i64 timestamp
}

struct TimeResponse {
	1: i64 sendTimestamp
    2: i64 recvTimestamp
    3: i64 respTimestamp
}

struct TransactionNotify {
    1: string identifier
    2: base.Transaction transaction
}

struct TransactionResponse {
    1: optional i64 timestamp
}

struct AddressRequest {

}

struct AddressResponse {
    1: list<base.TcpAddress> addresses
}

service ArgosSniffer {
    TimeResponse time(1: TimeRequest req)
    void transactionNotify(1: TransactionNotify req) 
    AddressResponse address(1: AddressRequest req)
}