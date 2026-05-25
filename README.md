# Crestbridge
A custom "bridge" between MQTT and the crestron XML format, written in GO.

# About
This app is a sole developer project to learn about crestron, MQTT, and smarthome automation. I enjoy writing it in my free time, and hope to continue to expand its functionality

I took this project on to bring my houses's lighting system into the modern age. When it was built, it had what was considered pretty advanced AV control with a centralized crestron panel. Unfortunately, that is no longer bleeding edge. Since the touch panel connects to a central processor over ethernet, I figured I could get my home server to talk to the processor if I could figure out the protocol. I began to intercept traffic between the touch panel and processor until I had a working understanding of the protocol (which is pretty simple xml over tcp). This projects makes use of that knowledge and converts it into MQTT events which allowed me to add anything in crestron to home assistant and apple home.