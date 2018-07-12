#!/bin/sh
openssl req -x509 -nodes -sha256 -newkey rsa:4096 -keyout $1.key -out $1.crt -subj '/C=OR/ST=Oregon/L=Portland/O=Crater Moon Development/CN=cmdev.com/emailAddress=steven.e.newton@gmail.com'
