#include "MQTT.h"

void callback(char* topic, byte* payload, unsigned int length);


 byte server[] = {127,0,0,1};
 MQTT client("m13.cloudmqtt.com", 12423, callback);

// recieve message
void callback(char* topic, byte* payload, unsigned int length) {
    char p[length + 1];
    memcpy(p, payload, length);
    p[length] = NULL;


    //LED color changed based on published message under topic "Colorchange"
    //used to make sure message queueing is working
    if (!strcmp(p, "RED"))
        RGB.color(255, 0, 0);
    else if (!strcmp(p, "GREEN"))
        RGB.color(0, 255, 0);
    else if (!strcmp(p, "BLUE"))
        RGB.color(0, 0, 255);
    else
        RGB.color(255, 255, 255);
    delay(1000);
}


void setup() {
    RGB.control(true);

    // connect to the server
    client.connect("sparkclient", "skupntaq", "yRS1mpvJp8su");

    // publish/subscribe
    if (client.isConnected()) {
        client.subscribe("colorChange");
    }
}

void loop() {
    if (client.isConnected())
        client.loop();
        
    if (client.isConnected()) {
        client.publish("testTopic","hello world");
    }
    
    delay(5000);
}