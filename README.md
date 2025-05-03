[![progress-banner](https://backend.codecrafters.io/progress/kafka/7c5ee849-f767-4cb4-b45e-658a273a16e1)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Toy Kafka

This my Go Kafka server built as a part of the ["Build Your Own Kafka" Challenge](https://codecrafters.io/challenges/kafka).

In this challenge, I build a toy Kafka clone that's capable of accepting and responding to multiple standard Kafka requests. This also helps to learn about encoding and decoding messages using the Kafka wire protocol. 

**Note**: If yyou want to try it out for yourself, head over to [codecrafters.io](https://codecrafters.io) and start the challenge.


## Features

- All the requests are validated for their binary format by the broker. If an invalid request is made, the broker does not return a response, but instead logs a detailed error message to the console showing where the parsing, or the handling of the request failed.
- Any [tagged fields](https://cwiki.apache.org/confluence/display/KAFKA/KIP-482%3A+The+Kafka+Protocol+should+Support+Optional+Tagged+Fields) on the request are ignored by the broker, and handles the request as if they were not present. This is done to ensure that the broker can handle requests from different versions of the Kafka protocol, and to allow for future extensions of the protocol without breaking compatibility with older versions.