curl --include \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Host: localhost:8000" \
     --header "Origin: http://localhost:8000" \
     --header "Sec-WebSocket-Key: xxx" \
     --header "Sec-WebSocket-Version: 13" \
http://localhost:8000/v2/ws
