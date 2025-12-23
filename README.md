# gortsplib-timeout-issue
Related to https://github.com/bluenviron/gortsplib/issues/968.

## Usage
This repository contains:
- A simple RTSP server that sends an RTP packet containing a mock payload every 5 seconds
- A simple RTSP client that connects to the RTSP server and reads the stream with a configurable timeout

To build:
```shell
make
```

Run the RTSP server:
```shell
./rtsp-server
```

Run the RTSP client with default timeout:
```shell
./rtsp-client 
2025/12/23 13:31:05 Dialing RTSP URL rtsp://localhost:8554/stream with timeout 10s
2025/12/23 13:31:07 Received RTP packet with timestamp 0, payload foobar
2025/12/23 13:31:12 Received RTP packet with timestamp 0, payload foobar
2025/12/23 13:31:17 Received RTP packet with timestamp 0, payload foobar
```

Run it with a timeout shorter than the RTP packet interval:
```shell
./rtsp-client --rtsp-timeout=2s
2025/12/23 13:31:53 Dialing RTSP URL rtsp://localhost:8554/stream with timeout 2s
2025/12/23 13:31:55 Failed to run RTSP client: RTSP wait: TCP timeout
```
