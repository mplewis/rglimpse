# rglimpse

View the status of your [rtorrent](https://github.com/rakshasa/rtorrent) torrents at-a-glance on your mobile device.

# Usage

Start the Docker container with the following environment variables:

| Name                  | Required? | Description                                                                                                                                                        | Example                |
| --------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------- |
| RTORRENT_HOST         | yes       | The hostname of the rtorrent server                                                                                                                                | `rtorrent.example.com` |
| RTORRENT_USERNAME     | yes       | The username of the rtorrent RPC user                                                                                                                              | `admin`                |
| RTORRENT_PASSWORD     | yes       | The password of the rtorrent RPC user                                                                                                                              | `swordfish`            |
| RTORRENT_PORT         | no        | The port of the rtorrent server. Defaults to 9080.                                                                                                                 | `4444`                 |
| RTORRENT_HTTPS        | no        | If set, uses HTTPS to connect to rtorrent. If unset, defaults to HTTP.                                                                                             | `true`                 |
| HOST                  | no        | The host on which the rglimpse app listens. Defaults to all interfaces.                                                                                            | `127.0.0.1`            |
| PORT                  | no        | The port of the rglimpse app. Defaults to 9081.                                                                                                                    | `5555`                 |
| MAX_CLIENTS           | no        | The maximum number of workers to use when fetching detailed torrent info. Increasing this number can speed up refreshes but may stress your server. Default is 16. | 32                     |
| DEFAULT_INTERVAL_SECS | no        | How often to refresh torrent info, in seconds. Default is 30.                                                                                                      | 120                    |
