namespace go protocols
namespace py protocols

struct SimpleResponse {
  1: required bool OK;
  2: optional i64  ID;
}

service Snowflake {
  SimpleResponse GetNextID();

  void Ping();
}
