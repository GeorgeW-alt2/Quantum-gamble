#include <Wire.h>
#include "SparkFun_VL53L1X.h" //Click here to get the library: http://librarymanager/All#SparkFun_VL53L1X

SFEVL53L1X distanceSensor;
int threshold = 10000; // Threshold to detect a peak

void setup(void)
{
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);
  Wire.begin();
  Serial.begin(115200);
  if (distanceSensor.begin() != 0) // Begin returns 0 on a good init
  {
    while (1)
      ;
  }
}

#define ARRAY_SIZE 20 // Define the size of the array to store signal rates

int signalRates[ARRAY_SIZE]; // Array to store signal rates
int currentIndex = 0;        // Index to keep track of current position in the array
int minSignalRate, maxSignalRate; // Variables to store minimum and maximum signal rates
bool minDetected = false;    // Flag to indicate if a minimum is detected
bool maxDetected = false;    // Flag to indicate if a maximum is detected

void loop(void)
{
  distanceSensor.startRanging();
  while (!distanceSensor.checkForDataReady())
  {
    delay(1);
  }

  int signalRate = distanceSensor.getSignalRate();
  //Serial.println(abs(signalRate));
  if (abs(signalRate) > 22000){
    Serial.print("1");
    }
  if (abs(signalRate) > 15000 && abs(signalRate) < 22000 ){
    Serial.print("0");
    }
  delay(10);
}
