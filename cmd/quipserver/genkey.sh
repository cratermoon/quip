#!/bin/sh

openssl req -x509 -nodes -newkey rsa:4096 -keyout $1.crt -out $1.key -subj '/C=OR/ST=Oregon/L=Portland/O=Crater Moon Development/CN=cmdev.com/emailAddress=steven.e.newton@gmail.com'
