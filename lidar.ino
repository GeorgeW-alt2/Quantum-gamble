#include <Wire.h>
#include "SparkFun_VL53L1X.h" //Click here to get the library: http://librarymanager/All#SparkFun_VL53L1X

SFEVL53L1X distanceSensor;

int previousSignalRate = 0;
int threshold = 10000; // Threshold to detect a peak
bool peakDetected = false; // Tracks if a peak was detected
bool outputOne = false; // Tracks if output 1 should be sent
int counter = 0;
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

void loop(void)
{
  distanceSensor.startRanging();
  while (!distanceSensor.checkForDataReady())
  {
    delay(1);
  }

  int signalRate = distanceSensor.getSignalRate();
  //Serial.println(signalRate);

    if (signalRate < 0 && peakDetected && counter > 0)
    {
      Serial.write("0"); // Output 0 if there is a dip within 10 cycles
      peakDetected = false; // Reset peakDetected after detecting the dip
      counter = 0;
    }

    if (signalRate >= threshold && counter > 5)
    {
        Serial.write("1"); // Output 1 if no dip was detected within 10 cycles
        peakDetected = true; // Reset peakDetected after reporting 1
        counter = 0;

    }

  

  counter+=1;
}
