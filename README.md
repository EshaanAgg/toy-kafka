[![progress-banner](https://backend.codecrafters.io/progress/kafka/7c5ee849-f767-4cb4-b45e-658a273a16e1)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Toy Kafka

This my Kafka server built as a part of the ["Build Your Own Kafka" challenge](https://codecrafters.io/challenges/kafka) in GoLang.

In this challenge, I build a toy Kafka clone that's capable of accepting and responding to multiple standard Kafka requests. This also helps to learn about encoding and decoding messages using the Kafka wire protocol. 

**Note**: If you want to try it out for yourself, head over to [codecrafters.io](https://codecrafters.io) and start the challenge.


## Features

- All the requests are validated for their binary format by the broker. If an invalid request is made, the broker does not return a response, but instead logs a detailed error message to the console showing where the parsing, or the handling of the request failed.
- Any [tagged fields](https://cwiki.apache.org/confluence/display/KAFKA/KIP-482%3A+The+Kafka+Protocol+should+Support+Optional+Tagged+Fields) on the request are ignored by the broker, and handles the request as if they were not present. This is done to ensure that the broker can handle requests from different versions of the Kafka protocol, and to allow for future extensions of the protocol without breaking compatibility with older versions.

### AutoDecode

Since the request & response formats in Kafka protocol are deeply nested and binary encoded, it is very easy to make mistakes while encoding and decoding the messages. To make this easier, I have implemented a simple [`AutoDecodeBody`](./app/datatypes/request/auto_decode.go) function that can be used to decode the bodies of the request by just defining the appropriate structs for them. This is achieved by using the `reflect` package in Go, and using the `kafka` tag to identify the fields that need to be decoded.

Supported primitive types for decoding:
- `int8`
- `int16`
- `int32`
- `int64`
- `varint`
- `varuint`
- `string`
- `compact_string`
- `nullable_string`
- `uuid`

In additional to using these primitive types, you can use the following aggregator types to decode complex types:
- `struct`: Represents a struct with a set of fields. The value should be a pointer to a struct type. All the struct's fields are expected to have empty tagged fields at their end.
- `inline_struct`: Similar to `struct`, but it is not expected to have empty tagged fields at the end. 
- `compact_array`: Represents a compact array of values. This type must be specified as `compact_array:type` where `type` is the type of the elements in the array.

It is the duty of the user to ensure that the types used for decoding are compatible with the types used for encoding. The decoder will not check for compatibility, and will simply decode the data as specified by the user. This means that if the user specifies an incompatible type, the decoder may produce unexpected results or even panic.