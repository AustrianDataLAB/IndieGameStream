# How to mount games into a CloudRetro instance
We are using the image `ghcr.io/giongto35/cloud-game/cloud-game:v3.0.5` (see https://github.com/giongto35/cloud-game/tree/v3.0.5).
You can use the image as the old image, but you have to additionally mount the game file into `/usr/local/share/cloud-game/assets/games`.
It may be useful to set `CLOUD_GAME_WORKER_DEBUG` and `CLOUD_GAME_COORDINATOR_DEBUG` to true. Furthermore, `CLOUD_GAME_WEBRTC_LOGLEVEL=-1` gives trace level logs for WebRTC.