package ginredisratelimiter

import "errors"

var ErrRedisError = errors.New("REDIS OPERATION FAIL")
var ErrIpNotRecognize = errors.New("IP NOT RECOGNIZE")
var TooManyRequest = errors.New("TOO MANY REQUEST")