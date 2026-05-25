# Crestbridge
A custom "bridge" between MQTT and the crestron XML format, written in GO.

# About
This app is a sole developer project to learn about crestron, MQTT, and smarthome automation. I enjoy writing it in my free time, and hope to continue to expand its functionality

I took this project on to bring my houses's lighting system into the modern age. When it was built, it had what was considered pretty advanced AV control with a centralized crestron panel. Unfortunately, that is no longer bleeding edge. Since the touch panel connects to a central processor over ethernet, I figured I could get my home server to talk to the processor if I could figure out the protocol. I began to intercept traffic between the smartphone control app and processor until I had a working understanding of the protocol (which is pretty simple xml over tcp). This projects makes use of that knowledge and converts it into MQTT events which allowed me to add anything in crestron to home assistant and apple home. 

# How to use

Here's a sample crestbridge.toml (placed in config folder):

This is a bit cryptic, so if someone is trying to do something similar for some reason, feel free to reach out! rohan@rohanjamal.ca

CrestronIP="YOURIP"
Port="YOURPORT"
Password="YOUR-PASSWORD"

[[Devices]]     #1 device, add one of these sections for each device  
Name="Island"   #Name for MQTT channel  
Room="Kitchen"  #Room for MQTT channel  
ID=12           #ID. Used to update states  
Up=514          #ID of the digital join to increase value  
Down=515        #ID of the digital join to decrease value  
Type="analog"   #Type can be analog or digital (ie dimmer vs fan)  

[[Devices]]  
Name="Amp"  
Room="Kitchen"  
ID=20                   #ID for state tracking  
CTRLON=[15,30,31,13]    #Sequence of digital joins to turn on the element  
CTRLOFF=[30,31,26]      #Sequence of digital joins to turn off the element  
Type="digital"  
StateString="Apple TV"  #String which reflects device being in "on" state  

You can find this specific information by placing a device between your touch panel/control app and processor and seeing what happens/what IDs are activated when specific devices are turned on/off.