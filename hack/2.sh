# run in webhookServer
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -subj "/CN=your-webhook" -nodes -addext "subjectAltName = IP:34.66.103.162"
export CA_BUNDLE=`cat cert.pem | base64 -w 0`