package errors

var UNAUTHORIZED_BODY = []byte(`{"error":"invalid or expired token"}`)
var NOT_FOUND_BODY = []byte(`{"error":"router+method not"}`)
var BAD_REQUEST_BODY = []byte(`{"error":"request contains invalid attributes"}`)
