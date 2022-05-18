include "base.thrift"
namespace go master

struct TimeSync {
	1: i64 sendTimestamp
    2: i64 recvTimestamp
    3: i64 respTimestamp
}

struct PingRequest {
	1: string identifier
    2: i64 timestamp
    3: optional i64 deltaTime
}

struct PingResponse {
    1: base.ResponseStatus status
    2: string   protocol
    3: TimeSync timeSync
}

struct ReportRequest {
    1: string identifier
    2: string method
    3: string protocol
    4: base.Transaction transaction
}

struct ReportResponse {
    1: base.ResponseStatus status
}

service ArgosMaster {
    PingResponse ping(1: PingRequest req)
    ReportResponse report(1: ReportRequest req)
}