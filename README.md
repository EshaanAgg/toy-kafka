[![progress-banner](https://backend.codecrafters.io/progress/kafka/7c5ee849-f767-4cb4-b45e-658a273a16e1)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Toy Kafka

This my Kafka server built as a part of the ["Build Your Own Kafka" challenge](https://codecrafters.io/challenges/kafka) in GoLang.

In this challenge, I build a toy Kafka clone that's capable of accepting and responding to multiple standard Kafka requests. This also helps to learn about encoding and decoding messages using the Kafka wire protocol. 

**Note**: If you want to try it out for yourself, head over to [codecrafters.io](https://codecrafters.io) and start the challenge.


## Features

- All the requests are validated for their binary format by the broker. If an invalid request is made, the broker does not return a response, but instead logs a detailed error message to the console showing where the parsing, or the handling of the request failed.
- Any [tagged fields](https://cwiki.apache.org/confluence/display/KAFKA/KIP-482%3A+The+Kafka+Protocol+should+Support+Optional+Tagged+Fields) on the request are ignored by the broker, and handles the request as if they were not present. This is done to ensure that the broker can handle requests from different versions of the Kafka protocol, and to allow for future extensions of the protocol without breaking compatibility with older versions.

### Auto Encode/Decode

Since the request & response formats in Kafka protocol are deeply nested and binary encoded, it is very easy to make mistakes while encoding and decoding the messages. To make this easier, the request object used in the implementation has a [`AutoDecode`](./app/datatypes/request/auto_encode.go) method and the response object has an [`AutoEncode`](./app/datatypes/response/auto_encode.go) method that can be used to handle the parsing for the same automatically. This is done by using the `kafka` tag on the fields of the request and response structs to specify the type of the field, and the order in which they should be encoded or decoded.

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
- `struct`: Represents a struct with a set of fields. All the struct's fields are expected to have empty tagged fields at their end. All the structs used in a request should be pointers to the struct type, where as in response they must be the struct type itself.
- `inline_struct`: Similar to `struct`, but it is not expected to have empty tagged fields at the end. The note aboout pointer vs value types still applies. 
- `compact_array`: Represents a compact array of values. This type must be specified as `compact_array:type` where `type` is the type of the elements in the array.

It is the duty of the user to ensure that the types used for structs are compatible with the struct definition. The decoder will not check for compatibility, and will simply decode the data as specified by the user. This means that if the user specifies an incompatible type, the decoder may produce unexpected results or even panic.