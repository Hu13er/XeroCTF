echo "generate Priv key"
openssl genrsa -out rsa.priv 2048

echo "extract Pub key"
openssl rsa -in rsa.priv -outform PEM -pubout -out rsa.pub
