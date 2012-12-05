#! /usr/bin/env python

import irc.client
import sys
import re
from bitly import *

GRUBER_URLINTEXT_PAT = re.compile(ur'(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s`!()\[\]{};:\'".,<>?\xab\xbb\u201c\u201d\u2018\u2019]))')

a=Api(login="bitly_username",apikey="bitly_apikey")

channel = "#channel"

def on_connect(connection, event):
    connection.join(channel)

def on_join(connection, event):
	pass

def on_pubmsg(connection, event):
    r = GRUBER_URLINTEXT_PAT.findall(event.arguments()[0])
    for url in r:
	if len(url[0]) > 20 and not "jenkins.local.ironclad.mobi/" in url[0]:
		try:
			text = a.shorten(url[0])
			connection.privmsg(channel, text)
		except:
			pass

def on_disconnect(connection, event):
    raise SystemExit()

def main():

    client = irc.client.IRC()
    try:
        c = client.server().connect("serverurl", 6697, "LinkBot", ssl=True, password="password")
    except irc.client.ServerConnectionError, x:
        print x
        raise SystemExit(1)

    c.add_global_handler("welcome", on_connect)
    c.add_global_handler("join", on_join)
    c.add_global_handler("pubmsg", on_pubmsg)
    c.add_global_handler("disconnect", on_disconnect)

    client.process_forever()

if __name__ == '__main__':
    main()
