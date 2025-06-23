// filepath: esp01-led-flash/esp01-led-flash/src/main.cpp
#include <Arduino.h>
#include <ESP8266WiFi.h>   // <-- Use this for ESP8266
#include <PubSubClient.h>
#include <time.h>
#include <DHT.h>
#include "secrets.h"
#include <vector>

// Room name
const char* roomName = "office";

// MQTT client ID
const char* mqtt_client_id = "ESP12E-Office-DHT11";

// NTP server for timestamp
const char* ntpServer = "pool.ntp.org";
const long  gmtOffset_sec = 0;
const int   daylightOffset_sec = 0;

// DHT11 settings
#define DHTPIN 2 // GPIO pin where the DHT11 is connected labeled as D4 on ESP8266 12E      
#define DHTTYPE DHT11
DHT dht(DHTPIN, DHTTYPE);

WiFiClient espClient;
PubSubClient client(espClient);

struct MqttMessage {
  String topic;
  String payload;
};

std::vector<MqttMessage> messageQueue;

// Helper to enqueue a message
void enqueueMessage(const String& topic, const String& payload) {
  Serial.println("Enqueuing message: " + payload + " to topic: " + topic);
  const size_t MAX_QUEUE = 60 * 12;
  if (messageQueue.size() == MAX_QUEUE) {
    messageQueue.erase(messageQueue.begin());
  }
  messageQueue.push_back({topic, payload});
}

// Helper to send all queued messages
void sendQueuedMessages() {
  while (!messageQueue.empty() && client.connected()) {
    MqttMessage& msg = messageQueue.front();
    Serial.println("Sending queued message to topic: " + msg.topic);
    client.publish(msg.topic.c_str(), msg.payload.c_str());
    messageQueue.erase(messageQueue.begin());
  }
}

void setup_wifi() {
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
  }
  Serial.print("Connected to WiFi. IP address: ");
  Serial.println(WiFi.localIP());
}

void reconnect() {
    client.connect(mqtt_client_id, mqtt_user, mqtt_pass);
}

void setup() {
  delay(5000);
  Serial.begin(9600); // Use 9600 baud rate for ESP8266
  setup_wifi();

  // Set time via NTP
  configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);

  // Wait for time to be set
  Serial.print("Waiting for NTP time sync: ");
  time_t now = time(nullptr);
  int retry = 0;
  const int retry_count = 20;
  while (now < 8 * 3600 * 2 && retry < retry_count) {
    delay(500);
    Serial.print(".");
    now = time(nullptr);
    retry++;
  }
  Serial.println();

  if (now < 8 * 3600 * 2) {
    Serial.println("Failed to get time from NTP server.");
  } else {
    Serial.println("Time synchronized.");
  }

  // Initialize MQTT client
  client.setServer(mqtt_server, mqtt_port);
  Serial.println("MQTT client initialized.");

  dht.begin(); // Initialize DHT sensor
  Serial.println("DHT sensor initialized.");
}

String getISO8601Time() {
  time_t now = time(nullptr);
  struct tm *timeinfo = gmtime(&now);
  if (timeinfo == nullptr) {
    return "";
  }
  char buf[25];
  strftime(buf, sizeof(buf), "%Y-%m-%dT%H:%M:%SZ", timeinfo);
  return String(buf);
}

unsigned long lastPublish = 0;
const unsigned long publishInterval = 60000; // 1 minute

void loop() {
  unsigned long nowMillis = millis();
  client.loop(); 
  
  if (nowMillis - lastPublish >= publishInterval) {
    lastPublish = nowMillis;

    float humidity = dht.readHumidity();
    float temperature = dht.readTemperature();

    if (isnan(temperature) || isnan(humidity)) {
      Serial.println("Failed to read from DHT sensor!");
      return;
    }

    String timestamp = getISO8601Time();

    String payload = "{";
    payload += "\"room\":\"" + String(roomName) + "\",";
    payload += "\"temperature\":" + String(temperature, 1) + ",";
    payload += "\"humidity\":" + String(humidity, 1) + ",";
    payload += "\"timestamp\":\"" + timestamp + "\"";
    payload += "}";

    String topic = "home/";
    String roomNoSpaces = String(roomName);
    roomNoSpaces.replace(" ", "");
    topic += roomNoSpaces;

    enqueueMessage(topic, payload);

    if (!client.connected()) {
      reconnect();
    }

    if (client.connected()) {
      sendQueuedMessages();
    }
  }
}