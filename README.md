# livestream-rtmp

Live streaming RTMP server written in Go.

This server accepts inbound RTMP connections (from apps like OBS), and then serves the uploaded media to viewers as HLS in realtime.

## Development
To develop this project locally, you'll need to have Go 1.16 or newer installed on your local machine.

In the root directory of the project, create a file named `.env` with the following contents:

```env
API_PASSCODE=3454f56ygfdsertyuio076rseryui76
```

These are just example values. The passcode can be any string, as long as it matches the passcode required by the `livestream-api` server you're running.

Once you've got a `.env` file, just run this to start the server:

```sh
go run .
```

## Streaming via RTMP

To stream into the RTMP server, you first need to have a stream ready-to-go on your `livestream-api` database. Then, take the **stream key** associated with your stream, and point your OBS stream to this URL:

```
rtmp://127.0.0.1/<STREAM-KEY-HERE>
```

## Playback via HLS (HTTP Live Streaming)

When you created your stream on the `livestream-api`, your stream was given an **identifier** value. The identifier is public (unlike the stream key), and you'll need this to view your stream.

By default, in development mode, you can view your stream at this URL:

```
https://127.0.0.1:8081/play/<STREAM-IDENTIFIER-HERE>
```

## Notes

You need a running `livestream-api` instance in order for this RTMP server to run correctly. Make sure you've set that up and it's working first.
