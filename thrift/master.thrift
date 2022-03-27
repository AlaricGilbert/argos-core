namespace go api

struct ResponseStatus {
    1: i32 code
    2: string message
}

struct RegisterRequest {
	1: string identifier
}

struct RegisterResponse {
    1: ResponseStatus status
}

struct PingRequest {
	1: string identifier
}

struct Address {
    1: list<i8> ip
    2: i32 port
}

struct NodeTask {
    1: string protocol
    2: list<Address> nodes
    3: optional Address main
}

struct PingResponse {
    1: ResponseStatus status
    2: NodeTask task
}

struct ReportRequest {
    1: string identifier
    2: list<i8> transaction
    3: Address source
}

struct ReportResponse {
    1: ResponseStatus status
}

service ArgosMaster {
    RegisterResponse register(1: RegisterRequest req)
    PingResponse ping(1: PingRequest req)
    ReportResponse report(1: ReportRequest req)
}