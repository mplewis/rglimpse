# rglimpse

View the status of your [rtorrent](https://github.com/rakshasa/rtorrent) torrents at-a-glance on your mobile device.

# Usage

Start the Docker container with the following environment variables:

| Name                    | Required? | Default          | Example                | Description                                                                                                                                         |
| ----------------------- | --------- | ---------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `RTORRENT_HOST`         | yes       |                  | `rtorrent.example.com` | The hostname of the rtorrent server                                                                                                                 |
| `RTORRENT_PORT`         | no        | `9080`           | `4444`                 | The port of the rtorrent server                                                                                                                     |
| `RTORRENT_USERNAME`     | yes       |                  | `admin`                | The username of the rtorrent RPC user                                                                                                               |
| `RTORRENT_PASSWORD`     | yes       |                  | `swordfish`            | The password of the rtorrent RPC user                                                                                                               |
| `RTORRENT_HTTPS`        | no        | _HTTP_           | `true`                 | If set, uses HTTPS to connect to rtorrent                                                                                                           |
| `HOST`                  | no        | _all interfaces_ | `127.0.0.1`            | The host on which the rglimpse app listens                                                                                                          |
| `PORT`                  | no        | `9081`           | `5555`                 | The port of the rglimpse app                                                                                                                        |
| `MAX_CLIENTS`           | no        | `16`             | `32`                   | The maximum number of workers to use when fetching detailed torrent info. Increasing this number can speed up refreshes but may stress your server. |
| `DEFAULT_INTERVAL_SECS` | no        | `30`             | `120`                  | How often to refresh torrent info, in seconds                                                                                                       |
