For work diaginfra
ssh keys for connection to servers in directory ~/.ssh/
conf/config.yaml  in directory /app/conf/

docker run -d -p 3000:3000 -v conf:/app/conf -v ~/.ssh:~/.ssh diaginfra