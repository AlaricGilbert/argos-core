namespace go base

// StatusOK is the status code for a successful request.
const i32 StatusOK = 0
// Status Code 10000 to 19999 are reserved for system errors.
const i32 StatusInternalError = 10000
const i32 StatusInvalidArgument = 10001
// Status Code 20000 to 20999 are reserved for master errors.
const i32 StatusIdentifierConflict = 20001
const i32 StatusTransactionConflict = 20002

const string MessageInternalError = "internal error"
const string MessageInvalidArgument = "invalid argument"
const string MessageIdentifierConflict = "another client with same identifier already connected"
const string MessageTransactionConflict = "transaction with same txid already exists"


struct TcpAddress {
    1: binary ip
    2: i32 port
}

struct ResponseStatus {
    1: i32 code
    2: string message
}

struct Transaction {
    1: i64 timestamp
    2: binary txid
    3: TcpAddress from
}