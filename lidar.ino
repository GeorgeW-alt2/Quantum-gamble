#include <Wire.h>
#include "SparkFun_VL53L1X.h" //Click here to get the library: http://librarymanager/All#SparkFun_VL53L1X

SFEVL53L1X distanceSensor;

int threshold = 10000; // Threshold to detect a peak
bool peakDetected = false; // Tracks if a peak was detected
int counter = 0;
int timing = 15;
bool tri = false;
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
    if (signalRate >= threshold)
    {
        peakDetected = true; // Reset peakDetected after reporting 1
        tri = true;

    }
    if (signalRate < 0 && peakDetected)
    {
      peakDetected = false; // Reset peakDetected after detecting the dip
      tri = true;
    }
  if (counter > timing && tri == true){
     tri = false;
     if (peakDetected)
    {
        Serial.write("1"); // Output 0 if there is a dip within 10 cycles
        counter = 0;
    }
    if (!peakDetected)
    {
        Serial.write("0"); // Output 0 if there is a dip within 10 cycles
        counter = 0;
    }
    peakDetected = false; 
    counter = 0;
  }
  counter+=1;
}
