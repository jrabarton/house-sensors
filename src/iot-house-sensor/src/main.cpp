// filepath: esp01-led-flash/esp01-led-flash/src/main.cpp
#include <Arduino.h>
#include <ESP8266WiFi.h>   // <-- Use this for ESP8266
#include <PubSubClient.h>
#include <time.h>
#include <DHT.h>
#include "secrets.h"

// Room name
const char* roomName = "office";

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
  Serial.print("Attempting MQTT connection...");
  while (!client.connected()) {
    client.connect("ESP01Client", mqtt_user, mqtt_pass);
    delay(500);
  }
  Serial.println("Connected to MQTT broker");
}

void setup() {
  delay(5000);
  Serial.begin(9600); // Use 9600 baud rate for ESP8266
  Serial.println("Starting setup...");
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
  Serial.println("Initializing MQTT client...");
  client.setServer(mqtt_server, mqtt_port);
  Serial.println("MQTT client initialized.");

  Serial.println("Connecting to DHT sensor");
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

void loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  // Read temperature and humidity from DHT11
  float humidity = dht.readHumidity();
  float temperature = dht.readTemperature();

  // Check if any reads failed
  if (isnan(temperature) || isnan(humidity)) {
    Serial.println("Failed to read from DHT sensor!");
    delay(2000);
    return;
  }

  // Get ISO8601 timestamp
  String timestamp = getISO8601Time();

  // Prepare JSON payload
  String payload = "{";
  payload += "\"room\":\"" + String(roomName) + "\",";
  payload += "\"temperature\":" + String(temperature, 1) + ",";
  payload += "\"humidity\":" + String(humidity, 1) + ",";
  payload += "\"timestamp\":\"" + timestamp + "\"";
  payload += "}";

  // Prepare topic with spaces removed from roomName
  String topic = "home/";
  String roomNoSpaces = String(roomName);
  roomNoSpaces.replace(" ", "");
  topic += roomNoSpaces;

  Serial.println("Publishing: " + payload);
  Serial.println("Topic: " + topic);

  client.publish(topic.c_str(), payload.c_str());

  delay(60000); // Wait 1 minute
}