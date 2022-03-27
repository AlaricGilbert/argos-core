namespace go api

struct TimeRequest {
	1: i64 timestamp
}

struct TimeResponse {
	1: i64 sendTimestamp
    2: i64 recvTimestamp
    3: i64 respTimestamp
}

struct TransactionRequest {
    1: i64 timestamp
    2: list<i8> transaction
}

struct TransactionResponse {
    1: optional i64 timestamp
}

struct Address {
    1: list<i8> ip
    2: i32 port
}

struct AddressRequest {

}

struct AddressResponse {
    1: list<Address> address
}

service ArgosSniffer {
    TimeResponse time(1: TimeRequest req)
    TransactionResponse transaction(1: TimeRequest req) 
    AddressResponse address(1: AddressRequest req)
}