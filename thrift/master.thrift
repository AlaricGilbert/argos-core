include "base.thrift"
namespace go master

struct PingRequest {
	1: string identifier
    2: list<string> protocols
}

struct NodeTask {
    1: string protocol
    2: list<base.TcpAddress> nodes
    3: optional base.TcpAddress main
}

struct PingResponse {
    1: base.ResponseStatus status
    2: NodeTask task
}

struct ReportRequest {
    1: string identifier
    2: base.Transaction transaction
}

struct ReportResponse {
    1: base.ResponseStatus status
}

service ArgosMaster {
    PingResponse ping(1: PingRequest req)
    ReportResponse report(1: ReportRequest req)
}